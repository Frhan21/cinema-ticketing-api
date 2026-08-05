package dto

type StudioResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Capacity   int    `json:"capacity"`
	Facilities string `json:"facilities"`
}
