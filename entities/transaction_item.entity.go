package entities

import (
	"time"

	"github.com/google/uuid"
)

type TransactionItem struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TransactionID uuid.UUID `gorm:"column:transaction_id;not null" json:"transaction_id"`
	TicketID      uuid.UUID `gorm:"column:ticket_id;not null;uniqueIndex" json:"ticket_id"`
	Price         float64   `gorm:"column:price;not null" json:"price"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	Transaction Transaction `gorm:"foreignKey:TransactionID;references:ID" json:"transaction"`
	Ticket      Ticket      `gorm:"foreignKey:TicketID;references:ID" json:"ticket"`
}
