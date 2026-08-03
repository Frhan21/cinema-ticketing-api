package promo

import "time"

type CreatePromoRequest struct {
	Code        string    `json:"code" binding:"required,min=3,max=50"`
	Description string    `json:"description"`
	Discount    float64   `json:"discount" binding:"required,min=1,max=100"`
	MaxUsage    int       `json:"max_usage"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

type UpdatePromoRequest struct {
	Description string    `json:"description"`
	Discount    float64   `json:"discount" binding:"omitempty,min=1,max=100"`
	MaxUsage    int       `json:"max_usage"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsActive    *bool     `json:"is_active"`
}
