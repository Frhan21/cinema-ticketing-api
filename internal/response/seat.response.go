package response

type SeatResponse struct {
	ID         string `json:"id"`
	StudioID   string `json:"studio_id"`
	SeatNumber  string `json:"seat_number"`
	IsAvailable bool   `json:"is_available"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type SeatAvailabilityResponse struct {
	ID         string `json:"id"`
	SeatNumber string `json:"seat_number"`
}
