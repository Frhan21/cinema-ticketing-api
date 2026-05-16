package model

import (
	"time"

	"github.com/google/uuid"
)

type Studio struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	Capacity   int       `json:"capacity" gorm:"not null"`
	Facilities string    `json:"facilities" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
