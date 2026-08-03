package promo

import (
	"time"

	"github.com/google/uuid"
)

type PromoResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Discount    float64   `json:"discount"`
	MaxUsage    int       `json:"max_usage"`
	UsedCount   int       `json:"used_count"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	IsActive    bool      `json:"is_active"`
}

type ValidatePromoResponse struct {
	Code            string  `json:"code"`
	Discount        float64 `json:"discount"`
	DiscountedPrice float64 `json:"discounted_price"`
	OriginalPrice   float64 `json:"original_price"`
}
