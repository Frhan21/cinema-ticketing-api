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

// GetAll implements [MovieRepository].
func (m *movieRepository) GetAll() ([]models.Movie, error) {
	var movies []models.Movie
	return movies, m.db.Table("movies").Find(&movies).Error
}

// GetById implements [MovieRepository].
func (m *movieRepository) GetById(id string) (models.Movie, error) {
	var movie *models.Movie
	err := m.db.Table("movies").Where("id = ?", id).First(&movie).Error
	if err != nil {
		return models.Movie{}, err
	}
	return *movie, nil
}

// Update implements [MovieRepository].
func (m *movieRepository) Update(id string, movie *models.Movie) (*models.Movie, error) {
	err := m.db.Table("movies").Where("id = ?", id).Updates(movie).Error
	if err != nil {
		return nil, err
	}
	return movie, nil
}

type MovieRepository interface {
	Create(movie *models.Movie) error
	Update(id string, movie *models.Movie) (*models.Movie, error)
	GetAll() ([]models.Movie, error)
	GetById(id string) (models.Movie, error)
	Delete(id string) error
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}
