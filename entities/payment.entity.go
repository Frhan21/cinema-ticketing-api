package entities

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID                    uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TransactionID         uuid.UUID     `gorm:"type:uuid;not null;uniqueIndex" json:"transaction_id"`
	OrderID               string        `gorm:"type:varchar(50);not null;uniqueIndex" json:"order_id"`
	MidtransTransactionID *string       `gorm:"type:varchar(100);uniqueIndex" json:"midtrans_transaction_id,omitempty"`
	SnapToken             string        `gorm:"type:text;not null" json:"-"`
	RedirectURL           string        `gorm:"type:text;not null" json:"redirect_url"`
	PaymentType           *string       `gorm:"type:varchar(50)" json:"payment_type,omitempty"`
	PaymentStatus         PaymentStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"payment_status"`
	GatewayStatus         string        `gorm:"type:varchar(30);not null;default:'pending'" json:"gateway_status"`
	GrossAmount           int64         `gorm:"type:bigint;not null" json:"gross_amount"`
	PaidAt                *time.Time    `gorm:"type:timestamp" json:"paid_at,omitempty"`
	CreatedAt             time.Time     `gorm:"type:timestamp;autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time     `gorm:"type:timestamp;autoUpdateTime" json:"updated_at"`

	Transaction Transaction `gorm:"foreignKey:TransactionID;references:ID" json:"-"`
}

type PaymentHistory struct {
	ID            uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PaymentID     uuid.UUID     `gorm:"type:uuid;not null;index" json:"payment_id"`
	PaymentStatus PaymentStatus `gorm:"type:varchar(20);not null" json:"payment_status"`
	GatewayStatus string        `gorm:"type:varchar(30);not null" json:"gateway_status"`
	CreatedAt     time.Time     `gorm:"type:timestamp;autoCreateTime" json:"created_at"`

	Payment Payment `gorm:"foreignKey:PaymentID;references:ID" json:"-"`
}
