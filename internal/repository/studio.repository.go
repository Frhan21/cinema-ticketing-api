package repository

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudioRepository interface {
	FindAll(page, perPage int) ([]model.Studio, int64, error)
	FindByID(id uuid.UUID) (*model.Studio, error)
	Create(studio *model.Studio) error
	Update(studio *model.Studio) error
	Delete(id uuid.UUID) error
}

type studioRepository struct {
	db *gorm.DB
}

// FindAll implements [StudioRepository].
func (s *studioRepository) FindAll(page, perPage int) ([]model.Studio, int64, error) {
	studios := []model.Studio{}
	var count int64

	err := s.db.Model(&model.Studio{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = s.db.Model(&model.Studio{}).Scopes(pagination.Paginate(page, perPage)).Find(&studios).Error
	return studios, count, err
}

// FindByID implements [StudioRepository].
func (s *studioRepository) FindByID(id uuid.UUID) (*model.Studio, error) {
	studio := &model.Studio{}
	err := s.db.Table("studios").Where("id = ?", id).First(&studio).Error
	return studio, err
}

// Create implements [StudioRepository].
func (s *studioRepository) Create(studio *model.Studio) error {
	return s.db.Table("studios").Create(studio).Error
}

// Delete implements [StudioRepository].
func (s *studioRepository) Delete(id uuid.UUID) error {
	return s.db.Table("studios").Where("id = ?", id).Delete(&model.Studio{}).Error
}

// Update implements [StudioRepository].
func (s *studioRepository) Update(studio *model.Studio) error {
	return s.db.Table("studios").Save(studio).Error
}

func NewStudioRepository(db *gorm.DB) StudioRepository {
	return &studioRepository{db: db}
}
