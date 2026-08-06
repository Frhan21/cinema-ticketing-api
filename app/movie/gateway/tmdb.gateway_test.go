package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"cinema-ticketing-api/config"
)

func TestTMDBGatewayFindByID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/movie/299534" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("language") != "id-ID" {
			t.Fatalf("expected Indonesian language query")
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Fatalf("unexpected authorization header")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": 299534,
			"title": "Avengers: Endgame",
			"overview": "Overview",
			"runtime": 181,
			"poster_path": "/poster.jpg",
			"genres": [{"name": "Action"}, {"name": "Adventure"}]
		}`))
	}))
	defer server.Close()

	gateway, err := NewTMDBGateway(config.TMDBConfig{BaseURL: server.URL, AccessToken: "test-token"})
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	movie, err := gateway.FindByID(context.Background(), 299534)
	if err != nil {
		t.Fatalf("unexpected find error: %v", err)
	}
	if movie.Title != "Avengers: Endgame" || movie.Duration != 181 {
		t.Fatalf("unexpected movie: %+v", movie)
	}
	if movie.Genre != "Action, Adventure" {
		t.Fatalf("unexpected genres: %s", movie.Genre)
	}
	if movie.PosterURL != "https://image.tmdb.org/t/p/w500/poster.jpg" {
		t.Fatalf("unexpected poster URL: %s", movie.PosterURL)
	}
}
