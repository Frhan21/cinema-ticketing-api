package repository

import (
	"cinema-ticketing-api/models"

	"gorm.io/gorm"
)

type StudioRepository interface {
	FindAll() ([]models.Studio, error)
	FindByID(id string) (*models.Studio, error)
	Create(studio *models.Studio) error
	Update(studio *models.Studio) error
	Delete(id string) error
}

type studioRepository struct {
	db *gorm.DB
}

// FindAll implements [StudioRepository].
func (s *studioRepository) FindAll() ([]models.Studio, error) {
	studios := []models.Studio{}
	err := s.db.Table("studios").Find(&studios).Error
	return studios, err
}

// FindByID implements [StudioRepository].
func (s *studioRepository) FindByID(id string) (*models.Studio, error) {
	studio := &models.Studio{}
	err := s.db.Table("studios").Where("id = ?", id).First(&studio).Error
	return studio, err
}

// Create implements [StudioRepository].
func (s *studioRepository) Create(studio *models.Studio) error {
	return s.db.Table("studios").Create(studio).Error
}

// Delete implements [StudioRepository].
func (s *studioRepository) Delete(id string) error {
	return s.db.Table("studios").Where("id = ?", id).Delete(&models.Studio{}).Error
}

// Update implements [StudioRepository].
func (s *studioRepository) Update(studio *models.Studio) error {
	return s.db.Table("studios").Where("id = ?", studio.ID).Updates(&studio).Error
}

func NewStudioRepository(db *gorm.DB) StudioRepository {
	return &studioRepository{db: db}
}
