package service

import (
	moviegateway "cinema-ticketing-api/app/movie/gateway"
	movierepository "cinema-ticketing-api/app/movie/repository"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/request"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type movieService struct {
	movieRepository movierepository.MovieRepository
	movieGateway    moviegateway.MovieGateway
}

// Create implements [MovieService].
func (m *movieService) Create(movie *entities.Movie) error {
	if movie.Title == "" || movie.Genre == "" || movie.Duration == 0 || movie.PosterUrl == "" {
		return apperror.NewBadRequestError("Fill the all required fields")
	}

	movie.ID = uuid.New()

	return m.movieRepository.Create(movie)
}

func (m *movieService) Import(ctx context.Context, tmdbID int64) (*entities.Movie, error) {
	if tmdbID <= 0 {
		return nil, apperror.NewBadRequestError("tmdb_id must be greater than zero")
	}

	movie, err := m.movieRepository.FindByTMDBID(tmdbID)
	if err == nil {
		return movie, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewInternalServerError("failed to find imported movie")
	}

	detail, err := m.movieGateway.FindByID(ctx, tmdbID)
	if err != nil {
		if errors.Is(err, moviegateway.ErrMovieNotFound) {
			return nil, apperror.NewNotFoundError("movie not found on TMDB")
		}
		return nil, apperror.NewInternalServerError("failed to fetch movie from TMDB")
	}
	if detail.TMDBID <= 0 || detail.Title == "" || detail.Genre == "" || detail.Duration <= 0 || detail.PosterURL == "" {
		return nil, apperror.NewBadRequestError("TMDB movie has incomplete booking metadata")
	}

	movie = &entities.Movie{
		ID:          uuid.New(),
		TMDBID:      &detail.TMDBID,
		Title:       detail.Title,
		Genre:       detail.Genre,
		Description: detail.Description,
		Duration:    detail.Duration,
		PosterUrl:   detail.PosterURL,
	}
	if err := m.movieRepository.Create(movie); err != nil {
		// A concurrent import may have inserted the same TMDB movie first.
		existing, findErr := m.movieRepository.FindByTMDBID(tmdbID)
		if findErr == nil {
			return existing, nil
		}
		return nil, apperror.NewInternalServerError("failed to import movie")
	}

	return movie, nil
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
func (m *movieService) FindAll(req request.PaginationRequest) ([]entities.Movie, int64, error) {
	movies, count, err := m.movieRepository.FindAll(req.Page, req.PerPage)

	if err != nil {
		return nil, 0, err
	}

	return movies, count, nil
}

// FindByID implements [MovieService].
func (m *movieService) FindByID(id string) (*entities.Movie, error) {
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
func (m *movieService) Update(movie *entities.Movie) error {
	return m.movieRepository.Update(movie)
}

type MovieService interface {
	Import(ctx context.Context, tmdbID int64) (*entities.Movie, error)
	Create(movie *entities.Movie) error
	Delete(id string) error
	Update(movie *entities.Movie) error
	FindAll(req request.PaginationRequest) ([]entities.Movie, int64, error)
	FindByID(id string) (*entities.Movie, error)
}

func NewMovieService(movieRepo movierepository.MovieRepository, movieGateway moviegateway.MovieGateway) MovieService {
	return &movieService{movieRepository: movieRepo, movieGateway: movieGateway}
}
