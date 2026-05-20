package model

import (
	"cinema-ticketing-api/internal/enums"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Email     string     `gorm:"unique;not null" json:"email"`
	Password  string     `gorm:"not null" json:"-"`
	Role      enums.Role `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
