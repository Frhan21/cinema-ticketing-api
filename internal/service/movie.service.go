package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrMovieNotFound  = errors.New("movie not found")
	ErrInvalidMovieID = errors.New("invalid movie id")
)

type movieService struct {
	movieRepository repository.MovieRepository
}

// Create implements [MovieService].
func (m *movieService) Create(movie *model.Movie) error {
	if movie.Title == "" || movie.Genre == "" || movie.Duration == 0 || movie.PosterUrl == "" {
		return errors.New("Fill the all required fields")
	}

	movie.ID = uuid.New()

	return m.movieRepository.Create(movie)
}

// Delete implements [MovieService].
func (m *movieService) Delete(id string) error {
	_, err := m.movieRepository.FindByID(id)

	if err != nil {
		return err
	}

	return m.movieRepository.Delete(id)
}

// FindAll implements [MovieService].
func (m *movieService) FindAll() ([]model.Movie, error) {
	movies, err := m.movieRepository.FindAll()

	if err != nil {
		return nil, err
	}

	return movies, nil
}

// FindByID implements [MovieService].
func (m *movieService) FindByID(id string) (*model.Movie, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidMovieID
	}

	movie, err := m.movieRepository.FindByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMovieNotFound
		}
		return nil, err
	}

	return movie, nil
}

// Update implements [MovieService].
func (m *movieService) Update(movie *model.Movie) error {
	return m.movieRepository.Update(movie)
}

type MovieService interface {
	Create(movie *model.Movie) error
	Delete(id string) error
	Update(movie *model.Movie) error
	FindAll() ([]model.Movie, error)
	FindByID(id string) (*model.Movie, error)
}

func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return &movieService{movieRepository: movieRepo}
}
