package service

import (
	"cinema-ticketing-api/models"
	"cinema-ticketing-api/repository"

	"github.com/google/uuid"
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
	return s.studioRepo.Delete(id)
}

// FindAll implements [StudioService].
func (s *studioService) FindAll() ([]models.Studio, error) {
	return s.studioRepo.FindAll()
}

// FindByID implements [StudioService].
func (s *studioService) FindByID(id string) (*models.Studio, error) {
	studio, err := s.studioRepo.FindByID(id)
	if err != nil {
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
