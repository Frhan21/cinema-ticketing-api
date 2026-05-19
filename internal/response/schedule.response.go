package response

import "github.com/google/uuid"

type ScheduleResponse struct {
	ID        uuid.UUID `json:"id"`
	MovieID   uuid.UUID `json:"movie_id"`
	StudioID  uuid.UUID `json:"studio_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Price     float64   `json:"price"`
}

type CreateScheduleResponse struct {
	ID        uuid.UUID `json:"id"`
	MovieID   uuid.UUID `json:"movie_id"`
	StudioID  uuid.UUID `json:"studio_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Price     float64   `json:"price"`
}

type UpdateScheduleResponse struct {
	ID        uuid.UUID `json:"id"`
	MovieID   uuid.UUID `json:"movie_id"`
	StudioID  uuid.UUID `json:"studio_id"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	Price     float64   `json:"price"`
}
