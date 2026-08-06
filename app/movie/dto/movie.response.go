package dto

type MovieResponse struct {
	ID          string `json:"id"`
	TMDBID      *int64 `json:"tmdb_id,omitempty"`
	Title       string `json:"title"`
	Genre       string `json:"genre"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	PosterUrl   string `json:"poster_url"`
}
