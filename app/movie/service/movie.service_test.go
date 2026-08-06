package service

import (
	"context"
	"testing"

	moviegateway "cinema-ticketing-api/app/movie/gateway"
	"cinema-ticketing-api/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type movieRepositoryStub struct {
	movies map[int64]*entities.Movie
}

func (r *movieRepositoryStub) Create(movie *entities.Movie) error {
	r.movies[*movie.TMDBID] = movie
	return nil
}

func (r *movieRepositoryStub) Update(*entities.Movie) error { return nil }

func (r *movieRepositoryStub) FindAll(int, int) ([]entities.Movie, int64, error) {
	return nil, 0, nil
}

func (r *movieRepositoryStub) FindByID(uuid.UUID) (*entities.Movie, error) {
	return nil, gorm.ErrRecordNotFound
}

func (r *movieRepositoryStub) FindByTMDBID(tmdbID int64) (*entities.Movie, error) {
	movie, ok := r.movies[tmdbID]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return movie, nil
}

func (r *movieRepositoryStub) Delete(uuid.UUID) error { return nil }

type movieGatewayStub struct {
	calls int
}

func (g *movieGatewayStub) FindByID(context.Context, int64) (*moviegateway.MovieDetail, error) {
	g.calls++
	return &moviegateway.MovieDetail{
		TMDBID:      299534,
		Title:       "Avengers: Endgame",
		Genre:       "Action, Adventure",
		Description: "Overview",
		Duration:    181,
		PosterURL:   "https://image.tmdb.org/t/p/w500/poster.jpg",
	}, nil
}

func TestImportMovieIsIdempotent(t *testing.T) {
	repository := &movieRepositoryStub{movies: make(map[int64]*entities.Movie)}
	gateway := &movieGatewayStub{}
	service := NewMovieService(repository, gateway)

	first, err := service.Import(context.Background(), 299534)
	if err != nil {
		t.Fatalf("unexpected first import error: %v", err)
	}
	second, err := service.Import(context.Background(), 299534)
	if err != nil {
		t.Fatalf("unexpected second import error: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("expected the same local movie ID")
	}
	if gateway.calls != 1 {
		t.Fatalf("expected one TMDB request, got %d", gateway.calls)
	}
}
