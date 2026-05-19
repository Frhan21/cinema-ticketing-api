package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrStudioNotFound  = errors.New("studio not found")
	ErrInvalidStudioID = errors.New("invalid studio id")
)

type StudioService interface {
	FindAll(req request.PaginationRequest) ([]model.Studio, int64, error)
	FindByID(id string) (*model.Studio, error)
	Create(studio *model.Studio) error
	Update(studio *model.Studio) error
	Delete(id string) error
}

type studioService struct {
	studioRepo repository.StudioRepository
}

// Create implements [StudioService].
func (s *studioService) Create(studio *model.Studio) error {
	studio.ID = uuid.New()
	return s.studioRepo.Create(studio)
}

// Delete implements [StudioService].
func (s *studioService) Delete(id string) error {
	_, err := s.FindByID(id)
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return ErrInvalidStudioID
	}
	return s.studioRepo.Delete(parsedID)
}

// FindAll implements [StudioService].
func (s *studioService) FindAll(req request.PaginationRequest) ([]model.Studio, int64, error) {
	return s.studioRepo.FindAll(req.Page, req.PerPage)
}

// FindByID implements [StudioService].
func (s *studioService) FindByID(id string) (*model.Studio, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, ErrInvalidStudioID
	}

	studio, err := s.studioRepo.FindByID(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudioNotFound
		}
		return nil, err
	}

	return studio, nil
}

// Update implements [StudioService].
func (s *studioService) Update(studio *model.Studio) error {
	return s.studioRepo.Update(studio)
}

func NewStudioService(studioRepo repository.StudioRepository) StudioService {
	return &studioService{studioRepo: studioRepo}
}
