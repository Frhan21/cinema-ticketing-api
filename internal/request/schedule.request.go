package request

import "github.com/google/uuid"

type ScheduleRequest struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	MovieID   uuid.UUID `json:"movie_id" binding:"required"`
	StudioID  uuid.UUID `json:"studio_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
}

type CreateSchedule struct {
	MovieID   uuid.UUID `json:"movie_id" binding:"required"`
	StudioID  uuid.UUID `json:"studio_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
}

type UpdateSchedule struct {
	ID        uuid.UUID `json:"id" binding:"required"`
	MovieID   uuid.UUID `json:"movie_id" binding:"required"`
	StudioID  uuid.UUID `json:"studio_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	EndTime   string    `json:"end_time" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
}
