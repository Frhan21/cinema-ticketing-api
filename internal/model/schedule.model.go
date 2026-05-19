package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MovieID   uuid.UUID `gorm:"type:uuid;not null" json:"movie_id"`
	StudioID  uuid.UUID `gorm:"type:uuid;not null" json:"studio_id"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relasi database
	Movie  Movie  `gorm:"foreignKey:MovieID;references:ID" json:"movie"`
	Studio Studio `gorm:"foreignKey:StudioID;references:ID" json:"studio"`
}
