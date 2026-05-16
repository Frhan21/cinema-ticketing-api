package repository

import (
	"cinema-ticketing-api/internal/model"

	"gorm.io/gorm"
)

type movieRepository struct {
	db *gorm.DB
}

// Create implements [MovieRepository].
func (m *movieRepository) Create(movie *model.Movie) error {
	return m.db.Table("movies").Create(movie).Error
}

// Delete implements [MovieRepository].
func (m *movieRepository) Delete(id string) error {
	return m.db.Table("movies").Where("id = ?", id).Delete(&model.Movie{}).Error
}

// FindAll implements [MovieRepository].
func (m *movieRepository) FindAll() ([]model.Movie, error) {
	var movies []model.Movie
	return movies, m.db.Table("movies").Find(&movies).Error
}

// FindByID implements [MovieRepository].
func (m *movieRepository) FindByID(id string) (*model.Movie, error) {
	var movie *model.Movie
	err := m.db.Table("movies").Where("id = ?", id).First(&movie).Error
	if err != nil {
		return nil, err
	}
	return movie, nil
}

// Update implements [MovieRepository].
func (m *movieRepository) Update(movie *model.Movie) error {
	return m.db.Table("movies").Where("id = ?", movie.ID).Updates(movie).Error
}

type MovieRepository interface {
	Create(movie *model.Movie) error
	Update(movie *model.Movie) error
	FindAll() ([]model.Movie, error)
	FindByID(id string) (*model.Movie, error)
	Delete(id string) error
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}
