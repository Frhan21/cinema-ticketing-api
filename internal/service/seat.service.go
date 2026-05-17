package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrSeatNotFound        = errors.New("seat not found")
	ErrInvalidSeatID       = errors.New("invalid seat id")
	ErrInvalidSeatStudioID = errors.New("invalid studio id")
)

type seatService struct {
	seatRepo repository.SeatRepository
}

// Create implements [SeatService].
func (s *seatService) Create(seat *model.Seat) error {
	seat.ID = uuid.New()
	return s.seatRepo.Create(seat)
}

// Delete implements [SeatService].
func (s *seatService) Delete(id string) error {
	_, err := s.FindByID(id)
	if err != nil {
		return err
	}
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.seatRepo.Delete(parsedID)
}

// FindAll implements [SeatService].
func (s *seatService) FindAll() ([]model.Seat, error) {
	return s.seatRepo.FindAll()
}

// FindByID implements [SeatService].
func (s *seatService) FindByID(id string) (*model.Seat, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidSeatID
	}

	seat, err := s.seatRepo.FindByID(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSeatNotFound
		}
		return nil, err
	}

	return seat, nil
}

// FindByStudioID implements [SeatService].
func (s *seatService) FindByStudioID(studioID string) ([]model.Seat, error) {
	parsedStudioID, err := uuid.Parse(studioID)
	if err != nil {
		return nil, ErrInvalidSeatStudioID
	}
	return s.seatRepo.FindByStudioID(parsedStudioID)
}

// Update implements [SeatService].
func (s *seatService) Update(seat *model.Seat) error {
	_, err := s.FindByID(seat.ID.String())
	if err != nil {
		return err
	}
	return s.seatRepo.Update(seat)
}

type SeatService interface {
	FindAll() ([]model.Seat, error)
	FindByID(id string) (*model.Seat, error)
	FindByStudioID(studioID string) ([]model.Seat, error)
	Create(seat *model.Seat) error
	Update(seat *model.Seat) error
	Delete(id string) error
}

func NewSeatService(seatRepo repository.SeatRepository) SeatService {
	return &seatService{seatRepo: seatRepo}
}
