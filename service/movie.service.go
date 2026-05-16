package service

import (
	"cinema-ticketing-api/models"
	"cinema-ticketing-api/repository"
	"errors"

	"github.com/google/uuid"
)

type movieService struct {
	movieRepository repository.MovieRepository
}

// Create implements [MovieService].
func (m *movieService) Create(movie *models.Movie) error {
	if movie.Title == "" || movie.Genre == "" || movie.Duration == 0 || movie.PosterUrl == "" {
		return errors.New("Fill the all required fields")
	}

	movie.ID = uuid.New()

	return m.movieRepository.Create(movie)
}

// Delete implements [MovieService].
func (m *movieService) Delete(id string) error {
	_, err := m.movieRepository.GetById(id)

	if err != nil {
		return err
	}

	return m.movieRepository.Delete(id)
}

// GetAll implements [MovieService].
func (m *movieService) GetAll() ([]models.Movie, error) {
	movies, err := m.movieRepository.GetAll()

	if err != nil {
		return nil, err
	}

	return movies, nil
}

// GetById implements [MovieService].
func (m *movieService) GetById(id string) (*models.Movie, error) {
	movie, err := m.movieRepository.GetById(id)

	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// Update implements [MovieService].
func (m *movieService) Update(id string, movie *models.Movie) (*models.Movie, error) {
	return m.movieRepository.Update(id, movie)
}

type MovieService interface {
	Create(movie *models.Movie) error
	Delete(id string) error
	Update(id string, movie *models.Movie) (*models.Movie, error)
	GetAll() ([]models.Movie, error)
	GetById(id string) (*models.Movie, error)
}

func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return &movieService{movieRepository: movieRepo}
}
