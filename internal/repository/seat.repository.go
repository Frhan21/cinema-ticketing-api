package repository

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type seatRepository struct {
	db *gorm.DB
}

// Create implements [SeatRepository].
func (s *seatRepository) Create(seat *model.Seat) error {
	return s.db.Table("seats").Create(seat).Error
}

// Delete implements [SeatRepository].
func (s *seatRepository) Delete(id uuid.UUID) error {
	return s.db.Table("seats").Where("id = ?", id).Delete(&model.Seat{}).Error
}

// FindAll implements [SeatRepository].
func (s *seatRepository) FindAll(page, perPage int) ([]model.Seat, int64, error) {
	var seats []model.Seat
	var count int64

	err := s.db.Model(&model.Seat{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Model(&model.Seat{}).Preload("Studio").Scopes(pagination.Paginate(page, perPage)).Find(&seats).Error
	return seats, count, err
}

// FindByID implements [SeatRepository].
func (s *seatRepository) FindByID(id uuid.UUID) (*model.Seat, error) {
	var seat *model.Seat
	err := s.db.Preload("Studio").Where("id = ?", id).First(&seat).Error
	if err != nil {
		return nil, err
	}
	return seat, nil
}

// FindByStudioID implements [SeatRepository].
func (s *seatRepository) FindByStudioID(studioID uuid.UUID) ([]model.Seat, error) {
	var seats []model.Seat
	return seats, s.db.Preload("Studio").Where("studio_id = ?", studioID).Find(&seats).Error
}

// Update implements [SeatRepository].
func (s *seatRepository) Update(seat *model.Seat) error {
	return s.db.Table("seats").Save(seat).Error
}

type SeatRepository interface {
	FindAll(page, perPage int) ([]model.Seat, int64, error)
	FindByID(id uuid.UUID) (*model.Seat, error)
	FindByStudioID(studioID uuid.UUID) ([]model.Seat, error)
	Create(seat *model.Seat) error
	Update(seat *model.Seat) error
	Delete(id uuid.UUID) error
}

func NewSeatRepository(db *gorm.DB) SeatRepository {
	return &seatRepository{db: db}
}
