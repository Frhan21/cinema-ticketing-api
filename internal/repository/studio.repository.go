package repository

import (
	"cinema-ticketing-api/internal/model"

	"gorm.io/gorm"
)

type StudioRepository interface {
	FindAll() ([]model.Studio, error)
	FindByID(id string) (*model.Studio, error)
	Create(studio *model.Studio) error
	Update(studio *model.Studio) error
	Delete(id string) error
}

type studioRepository struct {
	db *gorm.DB
}

// FindAll implements [StudioRepository].
func (s *studioRepository) FindAll() ([]model.Studio, error) {
	studios := []model.Studio{}
	err := s.db.Table("studios").Find(&studios).Error
	return studios, err
}

// FindByID implements [StudioRepository].
func (s *studioRepository) FindByID(id string) (*model.Studio, error) {
	studio := &model.Studio{}
	err := s.db.Table("studios").Where("id = ?", id).First(&studio).Error
	return studio, err
}

// Create implements [StudioRepository].
func (s *studioRepository) Create(studio *model.Studio) error {
	return s.db.Table("studios").Create(studio).Error
}

// Delete implements [StudioRepository].
func (s *studioRepository) Delete(id string) error {
	return s.db.Table("studios").Where("id = ?", id).Delete(&model.Studio{}).Error
}

// Update implements [StudioRepository].
func (s *studioRepository) Update(studio *model.Studio) error {
	return s.db.Table("studios").Where("id = ?", studio.ID).Updates(&studio).Error
}

func NewStudioRepository(db *gorm.DB) StudioRepository {
	return &studioRepository{db: db}
}
