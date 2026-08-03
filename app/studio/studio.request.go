package studio

type StudioRequest struct {
	Name       string `json:"name" binding:"required"`
	Capacity   int    `json:"capacity" binding:"required,gt=0"`
	Facilities string `json:"facilities" binding:"required"`
}

type UpdateStudioRequest struct {
	Name       *string `json:"name"`
	Capacity   *int    `json:"capacity" binding:"omitempty,gt=0"`
	Facilities *string `json:"facilities"`
}
