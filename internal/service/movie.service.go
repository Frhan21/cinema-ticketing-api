package service

import (
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/pkg/apperror"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type movieService struct {
	movieRepository repository.MovieRepository
}

// Create implements [MovieService].
func (m *movieService) Create(movie *model.Movie) error {
	if movie.Title == "" || movie.Genre == "" || movie.Duration == 0 || movie.PosterUrl == "" {
		return apperror.NewBadRequestError("Fill the all required fields")
	}

	movie.ID = uuid.New()

	return m.movieRepository.Create(movie)
}

// Delete implements [MovieService].
func (m *movieService) Delete(id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return apperror.NewBadRequestError("invalid movie id")
	}

	_, err = m.movieRepository.FindByID(parsedID)

	if err != nil {
		return err
	}

	return m.movieRepository.Delete(parsedID)
}

// FindAll implements [MovieService].
func (m *movieService) FindAll(req request.PaginationRequest) ([]model.Movie, int64, error) {
	movies, count, err := m.movieRepository.FindAll(req.Page, req.PerPage)

	if err != nil {
		return nil, 0, err
	}

	return movies, count, nil
}

// FindByID implements [MovieService].
func (m *movieService) FindByID(id string) (*model.Movie, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, apperror.NewBadRequestError("invalid movie id")
	}

	movie, err := m.movieRepository.FindByID(parsedID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewNotFoundError("movie not found")
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
	FindAll(req request.PaginationRequest) ([]model.Movie, int64, error)
	FindByID(id string) (*model.Movie, error)
}

func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return &movieService{movieRepository: movieRepo}
}
