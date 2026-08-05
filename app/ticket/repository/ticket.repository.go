package repository

import (
	"cinema-ticketing-api/entities"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrPromoUnavailable  = errors.New("promo is no longer available")
	ErrSeatAlreadyBooked = errors.New("seat is already booked")
)

type CreateBookingParams struct {
	Transaction *entities.Transaction
	Tickets     []entities.Ticket
	Items       []entities.TransactionItem
	PromoID     *uuid.UUID
	Now         time.Time
}

type TicketRepository interface {
	CreateBooking(params CreateBookingParams) error
	FindByID(id uuid.UUID) (*entities.Ticket, error)
	FindByUserID(userID uuid.UUID) ([]entities.Ticket, error)
	FindPaidTicketsBySchedule(scheduleID uuid.UUID) ([]entities.Ticket, error)
	FindByBookedSeats(scheduleID uuid.UUID, seatID uuid.UUID) ([]entities.Ticket, error)
	FindBookedSeatIDsBySchedule(scheduleID uuid.UUID) ([]uuid.UUID, error)
	UpdateStatusByIDs(ticketIDs []uuid.UUID, status entities.TicketStatus) error
	Delete(id uuid.UUID) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

// CreateBooking persists the booking aggregate and promo usage atomically.
func (r *ticketRepository) CreateBooking(params CreateBookingParams) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if params.PromoID != nil {
			result := tx.Model(&entities.Promo{}).
				Where(
					"id = ? AND is_active = ? AND start_date <= ? AND end_date >= ? AND (max_usage = 0 OR used_count < max_usage)",
					*params.PromoID,
					true,
					params.Now,
					params.Now,
				).
				UpdateColumn("used_count", gorm.Expr("used_count + 1"))
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return ErrPromoUnavailable
			}
		}

		if err := tx.Create(params.Transaction).Error; err != nil {
			return err
		}
		if err := tx.Create(&params.Tickets).Error; err != nil {
			if strings.Contains(err.Error(), "uq_active_ticket_schedule_seat") {
				return ErrSeatAlreadyBooked
			}
			return err
		}
		return tx.Create(&params.Items).Error
	})
}

// FindByID implements [TicketRepository].
func (r *ticketRepository) FindByID(id uuid.UUID) (*entities.Ticket, error) {
	var ticket entities.Ticket
	err := r.db.Where("id = ?", id).First(&ticket).Error
	return &ticket, err
}

// FindByUserID implements [TicketRepository].
func (r *ticketRepository) FindByUserID(userID uuid.UUID) ([]entities.Ticket, error) {
	var tickets []entities.Ticket
	err := r.db.Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

// FindPaidTicketsBySchedule implements [TicketRepository].
func (r *ticketRepository) FindPaidTicketsBySchedule(scheduleID uuid.UUID) ([]entities.Ticket, error) {
	var tickets []entities.Ticket
	err := r.db.Preload("User").
		Where("schedule_id = ? AND status = ?", scheduleID, entities.TicketStatusPaid).
		Find(&tickets).Error
	return tickets, err
}

// FindByBookedSeats implements [TicketRepository].
// Digunakan untuk cek apakah SATU seat tertentu sudah di-booking pada schedule ini.
// Hanya tiket yang statusnya bukan 'cancelled' yang dianggap aktif.
func (r *ticketRepository) FindByBookedSeats(scheduleID uuid.UUID, seatID uuid.UUID) ([]entities.Ticket, error) {
	var tickets []entities.Ticket
	err := r.db.Where("schedule_id = ? AND seat_id = ? AND status != ?", scheduleID, seatID, "cancelled").Find(&tickets).Error
	return tickets, err
}

// FindBookedSeatIDsBySchedule implements [TicketRepository].
// Digunakan untuk mengambil SEMUA seat_id yang sudah di-booking pada suatu schedule.
// Hasilnya dipakai untuk menyaring kursi mana saja yang masih tersedia.
func (r *ticketRepository) FindBookedSeatIDsBySchedule(scheduleID uuid.UUID) ([]uuid.UUID, error) {
	var seatIDs []uuid.UUID
	err := r.db.Model(&entities.Ticket{}).
		Where("schedule_id = ? AND status != ?", scheduleID, "cancelled").
		Pluck("seat_id", &seatIDs).Error
	return seatIDs, err
}

// UpdateStatusByIDs implements [TicketRepository].
func (r *ticketRepository) UpdateStatusByIDs(ticketIDs []uuid.UUID, status entities.TicketStatus) error {
	return r.db.Model(&entities.Ticket{}).Where("id IN ?", ticketIDs).Update("status", status).Error
}

// Delete implements [TicketRepository].
func (r *ticketRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&entities.Ticket{}).Error
}
