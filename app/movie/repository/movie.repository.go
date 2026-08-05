package repository

import (
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type movieRepository struct {
	db *gorm.DB
}

// Create implements [MovieRepository].
func (m *movieRepository) Create(movie *entities.Movie) error {
	return m.db.Table("movies").Create(movie).Error
}

// Delete implements [MovieRepository].
func (m *movieRepository) Delete(id uuid.UUID) error {
	return m.db.Table("movies").Where("id = ?", id).Delete(&entities.Movie{}).Error
}

func (m *movieRepository) FindAll(page, perPage int) ([]entities.Movie, int64, error) {
	var movies []entities.Movie
	var count int64

	err := m.db.Model(&entities.Movie{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = m.db.Model(&entities.Movie{}).Scopes(pagination.Paginate(page, perPage)).Find(&movies).Error
	return movies, count, err
}

// FindByID implements [MovieRepository].
func (m *movieRepository) FindByID(id uuid.UUID) (*entities.Movie, error) {
	var movie *entities.Movie
	err := m.db.Table("movies").Where("id = ?", id).First(&movie).Error
	if err != nil {
		return nil, err
	}
	return movie, nil
}

// Update implements [MovieRepository].
func (m *movieRepository) Update(movie *entities.Movie) error {
	return m.db.Table("movies").Save(movie).Error
}

type MovieRepository interface {
	Create(movie *entities.Movie) error
	Update(movie *entities.Movie) error
	FindAll(page, perPage int) ([]entities.Movie, int64, error)
	FindByID(id uuid.UUID) (*entities.Movie, error)
	Delete(id uuid.UUID) error
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}
