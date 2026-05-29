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

func (s *studioService) Delete(id string) error {
	_, err := s.FindByID(id)
	if err != nil {
		return err
	}
	parsedID, _ := uuid.Parse(id)
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
		return nil, apperror.NewBadRequestError("invalid studio id")
	}

	studio, err := s.studioRepo.FindByID(parsedID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("studio not found")
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
