package request

type SeatRequest struct {
	StudioID    string `json:"studio_id" binding:"required,uuid"`
	SeatNumber  string `json:"seat_number" binding:"required"`
	IsAvailable *bool  `json:"is_available" binding:"omitempty"`
}

type UpdateSeatRequest struct {
	ID          string `json:"id" binding:"required,uuid"`
	StudioID    string `json:"studio_id" binding:"required,uuid"`
	SeatNumber  string `json:"seat_number" binding:"required"`
	IsAvailable *bool  `json:"is_available" binding:"omitempty"`
}
