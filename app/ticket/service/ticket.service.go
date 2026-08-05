package service

import (
	ticketdto "cinema-ticketing-api/app/ticket/dto"
	ticketrepository "cinema-ticketing-api/app/ticket/repository"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
)

type TicketService interface {
	BookTicket(userID uuid.UUID, req *ticketdto.BookTicketRequest) (*entities.Transaction, []uuid.UUID, error)
	GetUserHistory(userID uuid.UUID) ([]entities.Transaction, error)
	GetAvailableSeats(scheduleID uuid.UUID) ([]entities.Seat, error)
}

type ticketRepository interface {
	CreateBooking(params ticketrepository.CreateBookingParams) error
	FindByBookedSeats(scheduleID uuid.UUID, seatID uuid.UUID) ([]entities.Ticket, error)
	FindBookedSeatIDsBySchedule(scheduleID uuid.UUID) ([]uuid.UUID, error)
}

type seatRepository interface {
	FindByID(id uuid.UUID) (*entities.Seat, error)
	FindByStudioID(studioID uuid.UUID) ([]entities.Seat, error)
}

type transactionRepository interface {
	FindByUserID(userID uuid.UUID) ([]entities.Transaction, error)
}

type scheduleRepository interface {
	FindByID(id uuid.UUID) (entities.Schedule, error)
}

type promoRepository interface {
	FindByCode(code string) (*entities.Promo, error)
}

type ticketService struct {
	ticketRepo      ticketRepository
	seatRepo        seatRepository
	transactionRepo transactionRepository
	scheduleRepo    scheduleRepository
	promoRepo       promoRepository
}

func NewTicketService(
	ticketRepo ticketRepository,
	seatRepo seatRepository,
	transactionRepo transactionRepository,
	scheduleRepo scheduleRepository,
	promoRepo promoRepository,
) TicketService {
	return &ticketService{
		ticketRepo:      ticketRepo,
		seatRepo:        seatRepo,
		transactionRepo: transactionRepo,
		scheduleRepo:    scheduleRepo,
		promoRepo:       promoRepo,
	}
}

// GetAvailableSeats mengembalikan daftar kursi yang BELUM di-booking untuk suatu schedule.
//
// Algoritma:
//  1. Ambil schedule → dapat studioID
//  2. Ambil SEMUA kursi di studio tersebut
//  3. Ambil semua seat_id yang sudah punya tiket aktif di schedule ini
//  4. Saring: kursi yang seat_id-nya TIDAK ADA di daftar booked = kursi tersedia
func (s *ticketService) GetAvailableSeats(scheduleID uuid.UUID) ([]entities.Seat, error) {
	// 1. Validasi schedule dan ambil studioID
	schedule, err := s.scheduleRepo.FindByID(scheduleID)
	if err != nil {
		return nil, apperror.NewNotFoundError("schedule not found")
	}

	// 2. Ambil semua kursi di studio tersebut
	allSeats, err := s.seatRepo.FindByStudioID(schedule.StudioID)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to get seats")
	}

	// 3. Ambil semua seat_id yang sudah di-booking pada schedule ini (status != cancelled)
	bookedSeatIDs, err := s.ticketRepo.FindBookedSeatIDsBySchedule(scheduleID)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to check booked seats")
	}

	// Buat map agar pencarian O(1) — lebih efisien daripada nested loop
	bookedMap := make(map[uuid.UUID]bool, len(bookedSeatIDs))
	for _, seatID := range bookedSeatIDs {
		bookedMap[seatID] = true
	}

	// 4. Filter: hanya kembalikan kursi yang belum ada di bookedMap
	var availableSeats []entities.Seat
	for _, seat := range allSeats {
		if !bookedMap[seat.ID] {
			availableSeats = append(availableSeats, seat)
		}
	}

	return availableSeats, nil
}

// BookTicket membuat tiket dan transaksi untuk booking kursi bioskop.
//
// Alur:
//  1. Validasi schedule ada
//  2. Validasi setiap seat: tersedia dan belum di-booking pada schedule ini
//  3. Buat Ticket (satu per seat)
//  4. Buat Transaction (header pembayaran)
//  5. Buat TransactionItem (penghubung Transaction → Ticket)
func (s *ticketService) BookTicket(userID uuid.UUID, req *ticketdto.BookTicketRequest) (*entities.Transaction, []uuid.UUID, error) {
	// 1. Validasi schedule
	schedule, err := s.scheduleRepo.FindByID(req.ScheduleId)
	if err != nil {
		return nil, nil, apperror.NewNotFoundError("schedule not found")
	}

	if len(req.SeatIds) == 0 {
		return nil, nil, apperror.NewBadRequestError("at least one seat is required")
	}
	if schedule.Price <= 0 || math.Trunc(schedule.Price) != schedule.Price {
		return nil, nil, apperror.NewInternalServerError("schedule price must be a positive whole IDR value")
	}

	// 2. Validasi setiap seat
	seenSeats := make(map[uuid.UUID]struct{}, len(req.SeatIds))
	for _, seatID := range req.SeatIds {
		if _, exists := seenSeats[seatID]; exists {
			return nil, nil, apperror.NewBadRequestError("duplicate seat selected: " + seatID.String())
		}
		seenSeats[seatID] = struct{}{}

		seat, err := s.seatRepo.FindByID(seatID)
		if err != nil {
			return nil, nil, apperror.NewNotFoundError("seat not found: " + seatID.String())
		}
		if !seat.IsAvailable {
			return nil, nil, apperror.NewBadRequestError("seat is not available: " + seat.SeatNumber)
		}
		if seat.StudioID != schedule.StudioID {
			return nil, nil, apperror.NewBadRequestError("seat does not belong to the schedule studio: " + seat.SeatNumber)
		}

		bookedTickets, err := s.ticketRepo.FindByBookedSeats(req.ScheduleId, seatID)
		if err != nil {
			return nil, nil, err
		}
		if len(bookedTickets) > 0 {
			return nil, nil, apperror.NewBadRequestError("seat already booked for this schedule: " + seat.SeatNumber)
		}
	}

	// 3. Buat Ticket (satu per seat)
	var tickets []entities.Ticket
	for _, seatID := range req.SeatIds {
		tickets = append(tickets, entities.Ticket{
			ID:         uuid.New(),
			UserId:     userID,
			ScheduleId: req.ScheduleId,
			SeatId:     seatID,
			Price:      schedule.Price,
			Status:     entities.TicketStatusPending,
		})
	}

	// 3a. Terapkan promo jika ada kode yang diberikan
	totalPrice := schedule.Price * float64(len(req.SeatIds))
	var promoID *uuid.UUID
	now := time.Now()
	if req.PromoCode != "" {
		promo, err := s.promoRepo.FindByCode(req.PromoCode)
		if err != nil {
			return nil, nil, apperror.NewNotFoundError("promo code not found")
		}
		if !promo.IsActive {
			return nil, nil, apperror.NewBadRequestError("promo is not active")
		}
		if now.Before(promo.StartDate) {
			return nil, nil, apperror.NewBadRequestError("promo has not started yet")
		}
		if now.After(promo.EndDate) {
			return nil, nil, apperror.NewBadRequestError("promo has expired")
		}
		if promo.MaxUsage > 0 && promo.UsedCount >= promo.MaxUsage {
			return nil, nil, apperror.NewBadRequestError("promo usage limit reached")
		}
		discount := math.Round(totalPrice * (promo.Discount / 100))
		totalPrice -= discount
		if totalPrice <= 0 {
			return nil, nil, apperror.NewBadRequestError("promo results in a non-payable transaction total")
		}
		promoID = &promo.ID
	}

	// 4. Buat Transaction
	transaction := entities.Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		TotalPrice:    totalPrice,
		PaymentStatus: "pending",
	}

	// 5. Buat TransactionItem (penghubung Transaction → Ticket)
	var transactionItems []entities.TransactionItem
	var ticketIDs []uuid.UUID
	for _, ticket := range tickets {
		transactionItems = append(transactionItems, entities.TransactionItem{
			ID:            uuid.New(),
			TransactionID: transaction.ID,
			TicketID:      ticket.ID,
			Price:         ticket.Price,
		})
		ticketIDs = append(ticketIDs, ticket.ID)
	}

	if err := s.ticketRepo.CreateBooking(ticketrepository.CreateBookingParams{
		Transaction: &transaction,
		Tickets:     tickets,
		Items:       transactionItems,
		PromoID:     promoID,
		Now:         now,
	}); err != nil {
		if errors.Is(err, ticketrepository.ErrPromoUnavailable) {
			return nil, nil, apperror.NewBadRequestError("promo is no longer available")
		}
		if errors.Is(err, ticketrepository.ErrSeatAlreadyBooked) {
			return nil, nil, apperror.NewBadRequestError("one or more selected seats are already booked")
		}
		return nil, nil, apperror.NewInternalServerError("failed to create booking")
	}

	return &transaction, ticketIDs, nil
}

// GetUserHistory mengambil semua transaksi milik user berdasarkan userID.
func (s *ticketService) GetUserHistory(userID uuid.UUID) ([]entities.Transaction, error) {
	transactions, err := s.transactionRepo.FindByUserID(userID)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to get user history")
	}
	return transactions, nil
}
