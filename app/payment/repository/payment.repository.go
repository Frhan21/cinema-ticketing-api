package repository

import (
	"cinema-ticketing-api/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type paymentRepository struct {
	db *gorm.DB
}

type PaymentRepository interface {
	Create(payload *entities.Payment) error
	FindByID(id uuid.UUID) (*entities.Payment, error)
	FindByTransactionID(transactionID uuid.UUID) (*entities.Payment, error)
	FindByOrderID(orderID string) (*entities.Payment, error)
	Update(payload *entities.Payment) error
	CreateHistory(history *entities.PaymentHistory) error
	ApplyStatus(payment *entities.Payment, history *entities.PaymentHistory, ticketStatus entities.TicketStatus) (bool, error)
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// Create implements [PaymentRepository].
func (p *paymentRepository) Create(payload *entities.Payment) error {
	return p.db.Create(payload).Error
}

// CreateHistory implements [PaymentRepository].
func (p *paymentRepository) CreateHistory(history *entities.PaymentHistory) error {
	return p.db.Create(history).Error
}

// FindByID implements [PaymentRepository].
func (p *paymentRepository) FindByID(id uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := p.db.First(&payment, "id = ?", id).Error
	return &payment, err
}

// FindByOrderID implements [PaymentRepository].
func (p *paymentRepository) FindByOrderID(orderID string) (*entities.Payment, error) {
	var payment entities.Payment
	err := p.db.Preload("Transaction.User").
		Preload("Transaction.Items.Ticket.Schedule.Movie").
		First(&payment, "order_id = ?", orderID).Error
	return &payment, err
}

// FindByTransactionID implements [PaymentRepository].
func (p *paymentRepository) FindByTransactionID(transactionID uuid.UUID) (*entities.Payment, error) {
	var payment entities.Payment
	err := p.db.First(&payment, "transaction_id = ?", transactionID).Error
	return &payment, err
}

// Update implements [PaymentRepository].
func (p *paymentRepository) Update(payload *entities.Payment) error {
	return p.db.Model(&entities.Payment{}).
		Where("id = ?", payload.ID).
		Updates(map[string]interface{}{
			"midtrans_transaction_id": payload.MidtransTransactionID,
			"payment_type":            payload.PaymentType,
			"payment_status":          payload.PaymentStatus,
			"gateway_status":          payload.GatewayStatus,
			"paid_at":                 payload.PaidAt,
		}).Error
}

// ApplyStatus updates the complete payment aggregate in one database transaction.
func (p *paymentRepository) ApplyStatus(payment *entities.Payment, history *entities.PaymentHistory, ticketStatus entities.TicketStatus) (bool, error) {
	statusTransitioned := false
	err := p.db.Transaction(func(tx *gorm.DB) error {
		var current entities.Payment
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&current, "id = ?", payment.ID).Error; err != nil {
			return err
		}
		statusChanged := current.PaymentStatus != payment.PaymentStatus || current.GatewayStatus != payment.GatewayStatus
		metadataChanged := !equalOptionalString(current.MidtransTransactionID, payment.MidtransTransactionID) ||
			!equalOptionalString(current.PaymentType, payment.PaymentType) ||
			!equalOptionalTime(current.PaidAt, payment.PaidAt)
		if !statusChanged && !metadataChanged {
			return nil
		}

		if err := tx.Model(&entities.Payment{}).
			Where("id = ?", payment.ID).
			Updates(map[string]interface{}{
				"midtrans_transaction_id": payment.MidtransTransactionID,
				"payment_type":            payment.PaymentType,
				"payment_status":          payment.PaymentStatus,
				"gateway_status":          payment.GatewayStatus,
				"paid_at":                 payment.PaidAt,
			}).Error; err != nil {
			return err
		}

		if statusChanged {
			history.PaymentID = payment.ID
			if err := tx.Create(history).Error; err != nil {
				return err
			}

			transactionUpdates := map[string]interface{}{
				"payment_status": payment.PaymentStatus,
			}
			if payment.PaymentType != nil {
				transactionUpdates["payment_method"] = *payment.PaymentType
			}
			if err := tx.Model(&entities.Transaction{}).
				Where("id = ?", payment.TransactionID).
				Updates(transactionUpdates).Error; err != nil {
				return err
			}

			ticketIDs := tx.Model(&entities.TransactionItem{}).
				Select("ticket_id").
				Where("transaction_id = ?", payment.TransactionID)
			if err := tx.Model(&entities.Ticket{}).
				Where("id IN (?)", ticketIDs).
				Update("status", ticketStatus).Error; err != nil {
				return err
			}
			statusTransitioned = true
		}
		return nil
	})
	return statusTransitioned, err
}

func equalOptionalString(left, right *string) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func equalOptionalTime(left, right *time.Time) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return left.Equal(*right)
}
