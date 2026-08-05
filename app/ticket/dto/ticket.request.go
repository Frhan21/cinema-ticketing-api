package dto

import "github.com/google/uuid"

type BookTicketRequest struct {
	ScheduleId uuid.UUID   `json:"schedule_id" binding:"required"`
	SeatIds    []uuid.UUID `json:"seat_id" binding:"required"`
	PromoCode  string      `json:"promo_code"` // opsional
}
