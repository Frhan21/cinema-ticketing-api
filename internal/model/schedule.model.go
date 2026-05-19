package model

import (
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	MovieID   uuid.UUID `gorm:"type:uuid;not null" json:"movie_id"`
	StudioID  uuid.UUID `gorm:"type:uuid;not null" json:"studio_id"`
	StartTime time.Time `gorm:"type:timestamptz;not null" json:"start_time"`
	EndTime   time.Time `gorm:"type:timestamptz;not null" json:"end_time"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `gorm:"type:timestamptz;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamptz;autoUpdateTime" json:"updated_at"`

	// Relasi database
	Movie  Movie  `gorm:"foreignKey:MovieID;references:ID" json:"movie"`
	Studio Studio `gorm:"foreignKey:StudioID;references:ID" json:"studio"`
}
