package dto

import "github.com/google/uuid"

type ScheduleResponse struct {
	ID        uuid.UUID      `json:"id"`
	MovieID   uuid.UUID      `json:"movie_id"`
	TMDBID    *int64         `json:"tmdb_id,omitempty"`
	StudioID  uuid.UUID      `json:"studio_id"`
	Studio    StudioResponse `json:"studio"`
	StartTime string         `json:"start_time"`
	EndTime   string         `json:"end_time"`
	Price     float64        `json:"price"`
}

type StudioResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type CreateScheduleResponse struct {
	ID        uuid.UUID `json:"id"`
	MovieID   uuid.UUID `json:"movie_id"`
	TMDBID    *int64    `json:"tmdb_id,omitempty"`
	StudioID  uuid.UUID `json:"studio_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Price     float64   `json:"price"`
}

type UpdateScheduleResponse struct {
	ID        uuid.UUID `json:"id"`
	MovieID   uuid.UUID `json:"movie_id"`
	TMDBID    *int64    `json:"tmdb_id,omitempty"`
	StudioID  uuid.UUID `json:"studio_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Price     float64   `json:"price"`
}
