package entities

import (
	"time"

	"github.com/google/uuid"
)

type Movie struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TMDBID      *int64    `gorm:"column:tmdb_id;uniqueIndex" json:"tmdb_id,omitempty"`
	Title       string    `json:"name" gorm:"not null"`
	Genre       string    `json:"genre" gorm:"not null"`
	Description string    `json:"description"`
	Duration    int       `json:"duration" gorm:"not null"`
	PosterUrl   string    `json:"poster_url" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
