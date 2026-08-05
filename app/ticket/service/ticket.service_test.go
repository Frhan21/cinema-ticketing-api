package service

import (
	ticketdto "cinema-ticketing-api/app/ticket/dto"
	ticketrepository "cinema-ticketing-api/app/ticket/repository"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

type ticketRepositoryStub struct {
	bookingParams *ticketrepository.CreateBookingParams
}

func (r *ticketRepositoryStub) CreateBooking(params ticketrepository.CreateBookingParams) error {
	r.bookingParams = &params
	return nil
}

func (r *ticketRepositoryStub) FindByBookedSeats(uuid.UUID, uuid.UUID) ([]entities.Ticket, error) {
	return nil, nil
}

func (r *ticketRepositoryStub) FindBookedSeatIDsBySchedule(uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

type seatRepositoryStub struct {
	seats map[uuid.UUID]*entities.Seat
}

func (r *seatRepositoryStub) FindByID(id uuid.UUID) (*entities.Seat, error) {
	seat, ok := r.seats[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return seat, nil
}

func (r *seatRepositoryStub) FindByStudioID(studioID uuid.UUID) ([]entities.Seat, error) {
	return nil, nil
}

type bookingTransactionRepositoryStub struct{}

func (r *bookingTransactionRepositoryStub) FindByUserID(uuid.UUID) ([]entities.Transaction, error) {
	return nil, nil
}

type scheduleRepositoryStub struct {
	schedule entities.Schedule
}

func (r *scheduleRepositoryStub) FindByID(uuid.UUID) (entities.Schedule, error) {
	return r.schedule, nil
}

type promoRepositoryStub struct {
	promo *entities.Promo
}

func (r *promoRepositoryStub) FindByCode(string) (*entities.Promo, error) {
	if r.promo == nil {
		return nil, errors.New("not found")
	}
	return r.promo, nil
}

func TestBookTicketCreatesAggregateWithPromo(t *testing.T) {
	studioID := uuid.New()
	scheduleID := uuid.New()
	seatOne := uuid.New()
	seatTwo := uuid.New()
	promoID := uuid.New()
	ticketRepo := &ticketRepositoryStub{}
	service := NewTicketService(
		ticketRepo,
		&seatRepositoryStub{seats: map[uuid.UUID]*entities.Seat{
			seatOne: {ID: seatOne, StudioID: studioID, SeatNumber: "A1", IsAvailable: true},
			seatTwo: {ID: seatTwo, StudioID: studioID, SeatNumber: "A2", IsAvailable: true},
		}},
		&bookingTransactionRepositoryStub{},
		&scheduleRepositoryStub{schedule: entities.Schedule{ID: scheduleID, StudioID: studioID, Price: 50000}},
		&promoRepositoryStub{promo: &entities.Promo{
			ID:        promoID,
			Code:      "SAVE20",
			Discount:  20,
			StartDate: time.Now().Add(-time.Hour),
			EndDate:   time.Now().Add(time.Hour),
			IsActive:  true,
		}},
	)

	transaction, ticketIDs, err := service.BookTicket(uuid.New(), &ticketdto.BookTicketRequest{
		ScheduleId: scheduleID,
		SeatIds:    []uuid.UUID{seatOne, seatTwo},
		PromoCode:  "SAVE20",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if transaction.TotalPrice != 80000 {
		t.Fatalf("expected total price 80000, got %.0f", transaction.TotalPrice)
	}
	if len(ticketIDs) != 2 || ticketRepo.bookingParams == nil {
		t.Fatal("expected complete booking aggregate")
	}
	if ticketRepo.bookingParams.PromoID == nil || *ticketRepo.bookingParams.PromoID != promoID {
		t.Fatal("expected promo usage to be included in atomic booking")
	}
	if len(ticketRepo.bookingParams.Tickets) != 2 || len(ticketRepo.bookingParams.Items) != 2 {
		t.Fatal("expected two tickets and two transaction items")
	}
}

func TestBookTicketRejectsSeatFromAnotherStudio(t *testing.T) {
	studioID := uuid.New()
	seatID := uuid.New()
	ticketRepo := &ticketRepositoryStub{}
	service := NewTicketService(
		ticketRepo,
		&seatRepositoryStub{seats: map[uuid.UUID]*entities.Seat{
			seatID: {ID: seatID, StudioID: uuid.New(), SeatNumber: "A1", IsAvailable: true},
		}},
		&bookingTransactionRepositoryStub{},
		&scheduleRepositoryStub{schedule: entities.Schedule{ID: uuid.New(), StudioID: studioID, Price: 50000}},
		&promoRepositoryStub{},
	)

	_, _, err := service.BookTicket(uuid.New(), &ticketdto.BookTicketRequest{SeatIds: []uuid.UUID{seatID}})
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %v", err)
	}
	if ticketRepo.bookingParams != nil {
		t.Fatal("booking must not be persisted")
	}
}
