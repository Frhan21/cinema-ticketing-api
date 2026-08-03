package ticket

import (
	"cinema-ticketing-api/app/promo"
	"cinema-ticketing-api/app/schedule"
	"cinema-ticketing-api/app/seat"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"

	"github.com/google/uuid"
)

type TicketService interface {
	BookTicket(userID uuid.UUID, req *BookTicketRequest) (*entities.Transaction, []uuid.UUID, error)
	GetUserHistory(userID uuid.UUID) ([]entities.Transaction, error)
	GetAvailableSeats(scheduleID uuid.UUID) ([]entities.Seat, error)
}

type ticketService struct {
	ticketRepo          TicketRepository
	seatRepo            seat.SeatRepository
	transactionRepo     TransactionRepository
	transactionItemRepo TransactionItemRepository
	scheduleRepo        schedule.ScheduleRepository
	promoRepo           promo.PromoRepository
}

func NewTicketService(
	ticketRepo TicketRepository,
	seatRepo seat.SeatRepository,
	transactionRepo TransactionRepository,
	transactionItemRepo TransactionItemRepository,
	scheduleRepo schedule.ScheduleRepository,
	promoRepo promo.PromoRepository,
) TicketService {
	return &ticketService{
		ticketRepo:          ticketRepo,
		seatRepo:            seatRepo,
		transactionRepo:     transactionRepo,
		transactionItemRepo: transactionItemRepo,
		scheduleRepo:        scheduleRepo,
		promoRepo:           promoRepo,
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
func (s *ticketService) BookTicket(userID uuid.UUID, req *BookTicketRequest) (*entities.Transaction, []uuid.UUID, error) {
	// 1. Validasi schedule
	schedule, err := s.scheduleRepo.FindByID(req.ScheduleId)
	if err != nil {
		return nil, nil, apperror.NewNotFoundError("schedule not found")
	}

	// 2. Validasi setiap seat
	for _, seatID := range req.SeatIds {
		seat, err := s.seatRepo.FindByID(seatID)
		if err != nil {
			return nil, nil, apperror.NewNotFoundError("seat not found: " + seatID.String())
		}
		if !seat.IsAvailable {
			return nil, nil, apperror.NewBadRequestError("seat is not available: " + seat.SeatNumber)
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

	if err := s.ticketRepo.CreateMany(tickets); err != nil {
		return nil, nil, apperror.NewInternalServerError("failed to create tickets: " + err.Error())
	}

	// 3a. Terapkan promo jika ada kode yang diberikan
	totalPrice := schedule.Price * float64(len(req.SeatIds))
	var promoID *uuid.UUID
	if req.PromoCode != "" {
		promo, err := s.promoRepo.FindByCode(req.PromoCode)
		if err != nil {
			return nil, nil, apperror.NewNotFoundError("promo code not found")
		}
		if !promo.IsActive {
			return nil, nil, apperror.NewBadRequestError("promo is not active")
		}
		if promo.MaxUsage > 0 && promo.UsedCount >= promo.MaxUsage {
			return nil, nil, apperror.NewBadRequestError("promo usage limit reached")
		}
		discount := totalPrice * (promo.Discount / 100)
		totalPrice -= discount
		promoID = &promo.ID
		// Increment usage count
		_ = s.promoRepo.IncrementUsage(promo.ID)
	}
	_ = promoID // reserved untuk future use (relasi transaksi ke promo)

	// 4. Buat Transaction
	transaction := entities.Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		TotalPrice:    totalPrice,
		PaymentStatus: "pending",
	}

	if err := s.transactionRepo.Create(&transaction); err != nil {
		return nil, nil, apperror.NewInternalServerError("failed to create transaction: " + err.Error())
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

	if err := s.transactionItemRepo.CreateMany(transactionItems); err != nil {
		return nil, nil, apperror.NewInternalServerError("failed to create transaction items: " + err.Error())
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
