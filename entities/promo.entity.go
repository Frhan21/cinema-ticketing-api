package entities

import (
	"time"

	"github.com/google/uuid"
)

// Promo adalah model untuk kode diskon yang dapat digunakan saat booking tiket.
type Promo struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Discount    float64   `gorm:"type:decimal(5,2);not null" json:"discount"` // persentase: 20 = 20%
	MaxUsage    int       `gorm:"default:0" json:"max_usage"`                 // 0 = unlimited
	UsedCount   int       `gorm:"default:0" json:"used_count"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
