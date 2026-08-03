package seat

import (
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/request"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type seatService struct {
	seatRepo SeatRepository
}

// Create implements [SeatService].
func (s *seatService) Create(seat *entities.Seat) error {
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
func (s *seatService) FindAll(req request.PaginationRequest) ([]entities.Seat, int64, error) {
	return s.seatRepo.FindAll(req.Page, req.PerPage)
}

// FindByID implements [SeatService].
func (s *seatService) FindByID(id string) (*entities.Seat, error) {
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
func (s *seatService) FindByStudioID(studioID string) ([]entities.Seat, error) {
	parsedStudioID, err := uuid.Parse(studioID)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid studio id")
	}
	return s.seatRepo.FindByStudioID(parsedStudioID)
}

// Update implements [SeatService].
func (s *seatService) Update(seat *entities.Seat) error {
	_, err := s.FindByID(seat.ID.String())
	if err != nil {
		return err
	}
	return s.seatRepo.Update(seat)
}

type SeatService interface {
	FindAll(req request.PaginationRequest) ([]entities.Seat, int64, error)
	FindByID(id string) (*entities.Seat, error)
	FindByStudioID(studioID string) ([]entities.Seat, error)
	Create(seat *entities.Seat) error
	Update(seat *entities.Seat) error
	Delete(id string) error
}

func NewSeatService(seatRepo SeatRepository) SeatService {
	return &seatService{seatRepo: seatRepo}
}
