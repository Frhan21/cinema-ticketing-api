package dto

import "github.com/google/uuid"

type CreateSchedule struct {
	MovieID   uuid.UUID `json:"movie_id" binding:"required"`
	StudioID  uuid.UUID `json:"studio_id" binding:"required"`
	StartTime string    `json:"start_time" binding:"required"`
	Price     float64   `json:"price" binding:"required"`
}

type UpdateSchedule struct {
	MovieID   *uuid.UUID `json:"movie_id" binding:"omitempty"`
	StudioID  *uuid.UUID `json:"studio_id" binding:"omitempty"`
	StartTime *string    `json:"start_time" binding:"omitempty"`
	Price     *float64   `json:"price" binding:"omitempty,gt=0"`
}
