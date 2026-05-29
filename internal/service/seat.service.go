package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/pkg/apperror"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
func (s *seatService) FindAll(req request.PaginationRequest) ([]model.Seat, int64, error) {
	return s.seatRepo.FindAll(req.Page, req.PerPage)
}

// FindByID implements [SeatService].
func (s *seatService) FindByID(id string) (*model.Seat, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid seat id")
	}

	seat, err := s.seatRepo.FindByID(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("seat not found")
		}
		return nil, err
	}

	return seat, nil
}

// FindByStudioID implements [SeatService].
func (s *seatService) FindByStudioID(studioID string) ([]model.Seat, error) {
	parsedStudioID, err := uuid.Parse(studioID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid studio id")
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
	FindAll(req request.PaginationRequest) ([]model.Seat, int64, error)
	FindByID(id string) (*model.Seat, error)
	FindByStudioID(studioID string) ([]model.Seat, error)
	Create(seat *model.Seat) error
	Update(seat *model.Seat) error
	Delete(id string) error
}

func NewSeatService(seatRepo repository.SeatRepository) SeatService {
	return &seatService{seatRepo: seatRepo}
}
