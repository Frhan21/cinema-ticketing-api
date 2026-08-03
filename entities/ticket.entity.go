package entities

import (
	"time"

	"github.com/google/uuid"
)

type Ticket struct {
	ID         uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	UserId     uuid.UUID    `gorm:"column:user_id;not null" json:"user_id"`
	ScheduleId uuid.UUID    `gorm:"column:schedule_id;not null" json:"schedule_id"`
	SeatId     uuid.UUID    `gorm:"column:seat_id;not null" json:"seat_id"`
	Price      float64      `gorm:"column:price;not null" json:"price"`
	Status     TicketStatus `gorm:"type:varchar(20);column:status;not null;default:'pending'" json:"status"`
	CreatedAt  time.Time    `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time    `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at"`

	// Relation
	User     User     `gorm:"foreignKey:UserId;references:ID" json:"user"`
	Schedule Schedule `gorm:"foreignKey:ScheduleId;references:ID" json:"schedule"`
	Seat     Seat     `gorm:"foreignKey:SeatId;references:ID" json:"seat"`
}
