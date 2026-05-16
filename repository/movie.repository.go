package repository

import (
	"cinema-ticketing-api/models"

	"gorm.io/gorm"
)

type movieRepository struct {
	db *gorm.DB
}

// Create implements [MovieRepository].
func (m *movieRepository) Create(movie *models.Movie) error {
	return m.db.Table("movies").Create(movie).Error
}

// Delete implements [MovieRepository].
func (m *movieRepository) Delete(id string) error {
	return m.db.Table("movies").Where("id = ?", id).Delete(&models.Movie{}).Error
}

// FindAll implements [MovieRepository].
func (m *movieRepository) FindAll() ([]models.Movie, error) {
	var movies []models.Movie
	return movies, m.db.Table("movies").Find(&movies).Error
}

// FindByID implements [MovieRepository].
func (m *movieRepository) FindByID(id string) (*models.Movie, error) {
	var movie *models.Movie
	err := m.db.Table("movies").Where("id = ?", id).First(&movie).Error
	if err != nil {
		return nil, err
	}
	return movie, nil
}

// Update implements [MovieRepository].
func (m *movieRepository) Update(movie *models.Movie) error {
	return m.db.Table("movies").Where("id = ?", movie.ID).Updates(movie).Error
}

type MovieRepository interface {
	Create(movie *models.Movie) error
	Update(movie *models.Movie) error
	FindAll() ([]models.Movie, error)
	FindByID(id string) (*models.Movie, error)
	Delete(id string) error
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}
