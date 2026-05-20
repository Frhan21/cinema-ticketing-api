package model

import (
	"cinema-ticketing-api/internal/enums"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID            `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID            `gorm:"column:user_id;not null" json:"user_id"`
	TotalPrice    float64              `gorm:"column:total_price;not null" json:"total_price"`
	PaymentMethod enums.PaymentMethod  `gorm:"type:varchar(30)" json:"payment_method"`
	PaymentStatus enums.PaymentStatus  `gorm:"type:varchar(20);not null;default:'pending'" json:"payment_status"`
	CreatedAt     time.Time            `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time            `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	User  User              `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Items []TransactionItem `gorm:"foreignKey:TransactionID;references:ID" json:"items"`
}
