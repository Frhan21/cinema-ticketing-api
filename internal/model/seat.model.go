package model

import (
	"time"

	"github.com/google/uuid"
)

type Seat struct {
	ID          uuid.UUID `json:"id" gorm:"type:varchar(36);primaryKey"`
	StudioID    uuid.UUID `json:"studio_id" gorm:"type:varchar(36);not null"`
	SeatNumber  string    `json:"seat_number" gorm:"type:varchar(100);not null"`
	IsAvailable bool      `json:"is_available" gorm:"type:boolean;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Relation
	Studio Studio `json:"studio" gorm:"foreignKey:StudioID;references:ID;constraint:OnDelete:CASCADE"`
}
