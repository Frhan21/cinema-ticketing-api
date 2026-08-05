package service

import (
	studiorepository "cinema-ticketing-api/app/studio/repository"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/request"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudioService interface {
	FindAll(req request.PaginationRequest) ([]entities.Studio, int64, error)
	FindByID(id string) (*entities.Studio, error)
	Create(studio *entities.Studio) error
	Update(studio *entities.Studio) error
	Delete(id string) error
}

type studioService struct {
	studioRepo studiorepository.StudioRepository
}

// Create implements [StudioService].
func (s *studioService) Create(studio *entities.Studio) error {
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
func (s *studioService) FindAll(req request.PaginationRequest) ([]entities.Studio, int64, error) {
	return s.studioRepo.FindAll(req.Page, req.PerPage)
}

// FindByID implements [StudioService].
func (s *studioService) FindByID(id string) (*entities.Studio, error) {
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
func (s *studioService) Update(studio *entities.Studio) error {
	return s.studioRepo.Update(studio)
}

func NewStudioService(studioRepo studiorepository.StudioRepository) StudioService {
	return &studioService{studioRepo: studioRepo}
}
