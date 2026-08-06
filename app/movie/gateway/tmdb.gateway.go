package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"cinema-ticketing-api/config"
)

var ErrMovieNotFound = errors.New("TMDB movie not found")

type MovieDetail struct {
	TMDBID      int64
	Title       string
	Genre       string
	Description string
	Duration    int
	PosterURL   string
}

type MovieGateway interface {
	FindByID(ctx context.Context, tmdbID int64) (*MovieDetail, error)
}

type tmdbGateway struct {
	baseURL     string
	accessToken string
	httpClient  *http.Client
}

type tmdbMovieResponse struct {
	ID         int64  `json:"id"`
	Title      string `json:"title"`
	Overview   string `json:"overview"`
	Runtime    int    `json:"runtime"`
	PosterPath string `json:"poster_path"`
	Genres     []struct {
		Name string `json:"name"`
	} `json:"genres"`
}

func NewTMDBGateway(cfg config.TMDBConfig) (MovieGateway, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if _, err := url.ParseRequestURI(baseURL); err != nil || baseURL == "" {
		return nil, fmt.Errorf("valid TMDB base URL is required")
	}
	accessToken := strings.TrimSpace(cfg.AccessToken)
	if accessToken == "" {
		return nil, fmt.Errorf("TMDB access token is required")
	}

	return &tmdbGateway{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (g *tmdbGateway) FindByID(ctx context.Context, tmdbID int64) (*MovieDetail, error) {
	requestURL := fmt.Sprintf("%s/movie/%d?language=id-ID", g.baseURL, tmdbID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create TMDB request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.accessToken)
	req.Header.Set("Accept", "application/json")

	res, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request TMDB movie: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrMovieNotFound
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("TMDB returned status %d", res.StatusCode)
	}

	var payload tmdbMovieResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode TMDB movie: %w", err)
	}

	genres := make([]string, 0, len(payload.Genres))
	for _, genre := range payload.Genres {
		if name := strings.TrimSpace(genre.Name); name != "" {
			genres = append(genres, name)
		}
	}
	posterURL := ""
	if payload.PosterPath != "" {
		posterURL = "https://image.tmdb.org/t/p/w500" + payload.PosterPath
	}

	return &MovieDetail{
		TMDBID:      payload.ID,
		Title:       payload.Title,
		Genre:       strings.Join(genres, ", "),
		Description: payload.Overview,
		Duration:    payload.Runtime,
		PosterURL:   posterURL,
	}, nil
}
