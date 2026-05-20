package repository

import (
	"cinema-ticketing-api/internal/enums"
	"cinema-ticketing-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateMany(tickets []model.Ticket) error
	FindByID(id uuid.UUID) (*model.Ticket, error)
	FindByUserID(userID uuid.UUID) ([]model.Ticket, error)
	FindByBookedSeats(scheduleID uuid.UUID, seatID uuid.UUID) ([]model.Ticket, error)
	FindBookedSeatIDsBySchedule(scheduleID uuid.UUID) ([]uuid.UUID, error)
	UpdateStatusByIDs(ticketIDs []uuid.UUID, status enums.TicketStatus) error
	Delete(id uuid.UUID) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

// CreateMany implements [TicketRepository].
func (r *ticketRepository) CreateMany(tickets []model.Ticket) error {
	return r.db.Create(&tickets).Error
}

// FindByID implements [TicketRepository].
func (r *ticketRepository) FindByID(id uuid.UUID) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.Where("id = ?", id).First(&ticket).Error
	return &ticket, err
}

// FindByUserID implements [TicketRepository].
func (r *ticketRepository) FindByUserID(userID uuid.UUID) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Where("user_id = ?", userID).Find(&tickets).Error
	return tickets, err
}

// FindByBookedSeats implements [TicketRepository].
// Digunakan untuk cek apakah SATU seat tertentu sudah di-booking pada schedule ini.
// Hanya tiket yang statusnya bukan 'cancelled' yang dianggap aktif.
func (r *ticketRepository) FindByBookedSeats(scheduleID uuid.UUID, seatID uuid.UUID) ([]model.Ticket, error) {
	var tickets []model.Ticket
	err := r.db.Where("schedule_id = ? AND seat_id = ? AND status != ?", scheduleID, seatID, "cancelled").Find(&tickets).Error
	return tickets, err
}

// FindBookedSeatIDsBySchedule implements [TicketRepository].
// Digunakan untuk mengambil SEMUA seat_id yang sudah di-booking pada suatu schedule.
// Hasilnya dipakai untuk menyaring kursi mana saja yang masih tersedia.
func (r *ticketRepository) FindBookedSeatIDsBySchedule(scheduleID uuid.UUID) ([]uuid.UUID, error) {
	var seatIDs []uuid.UUID
	err := r.db.Model(&model.Ticket{}).
		Where("schedule_id = ? AND status != ?", scheduleID, "cancelled").
		Pluck("seat_id", &seatIDs).Error
	return seatIDs, err
}

// UpdateStatusByIDs implements [TicketRepository].
func (r *ticketRepository) UpdateStatusByIDs(ticketIDs []uuid.UUID, status enums.TicketStatus) error {
	return r.db.Model(&model.Ticket{}).Where("id IN ?", ticketIDs).Update("status", status).Error
}

// Delete implements [TicketRepository].
func (r *ticketRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Ticket{}).Error
}
