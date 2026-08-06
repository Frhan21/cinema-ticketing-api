package dto

type MovieRequest struct {
	TMDBID      *int64 `json:"tmdb_id" binding:"omitempty,gt=0"`
	Title       string `json:"title" binding:"required"`
	Genre       string `json:"genre" binding:"required"`
	Description string `json:"description" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
	PosterUrl   string `json:"poster_url" binding:"required"`
}

type ImportMovieRequest struct {
	TMDBID int64 `json:"tmdb_id" binding:"required,gt=0"`
}

type UpdateMovieRequest struct {
	Title       *string `json:"title" binding:"omitempty"`
	Genre       *string `json:"genre" binding:"omitempty"`
	Description *string `json:"description" binding:"omitempty"`
	Duration    *int    `json:"duration" binding:"omitempty"`
	PosterUrl   *string `json:"poster_url" binding:"omitempty"`
}
