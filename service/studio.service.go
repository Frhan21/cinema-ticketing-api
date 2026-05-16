package service

import (
	"cinema-ticketing-api/models"
	"cinema-ticketing-api/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrStudioNotFound  = errors.New("studio not found")
	ErrInvalidStudioID = errors.New("invalid studio id")
)

type StudioService interface {
	FindAll() ([]models.Studio, error)
	FindByID(id string) (*models.Studio, error)
	Create(studio *models.Studio) error
	Update(studio *models.Studio) error
	Delete(id string) error
}

type studioService struct {
	studioRepo repository.StudioRepository
}

// Create implements [StudioService].
func (s *studioService) Create(studio *models.Studio) error {
	studio.ID = uuid.New()
	return s.studioRepo.Create(studio)
}

// Delete implements [StudioService].
func (s *studioService) Delete(id string) error {
	_, err := s.FindByID(id)
	if err != nil {
		return err
	}
	return s.studioRepo.Delete(id)
}

// FindAll implements [StudioService].
func (s *studioService) FindAll() ([]models.Studio, error) {
	return s.studioRepo.FindAll()
}

// FindByID implements [StudioService].
func (s *studioService) FindByID(id string) (*models.Studio, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidStudioID
	}

	studio, err := s.studioRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStudioNotFound
		}
		return nil, err
	}

	return studio, nil
}

// Update implements [StudioService].
func (s *studioService) Update(studio *models.Studio) error {
	return s.studioRepo.Update(studio)
}

func NewStudioService(studioRepo repository.StudioRepository) StudioService {
	return &studioService{studioRepo: studioRepo}
}
