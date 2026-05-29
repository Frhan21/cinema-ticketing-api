package repository

import (
	"cinema-ticketing-api/internal/enums"
	"cinema-ticketing-api/internal/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	FindByID(id uuid.UUID) (*model.Transaction, error)
	FindByUserID(userID uuid.UUID) ([]model.Transaction, error)
	FindAll() ([]model.Transaction, error)
	FindExpiredPendingTransactions(expiredBefore time.Time) ([]model.Transaction, error)
	UpdateStatus(id uuid.UUID, status enums.PaymentStatus) error
	UpdateStatusAndMethod(id uuid.UUID, status enums.PaymentStatus, method enums.PaymentMethod) error
	Delete(id uuid.UUID) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create implements [TransactionRepository].
func (r *transactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

// FindByID implements [TransactionRepository].
func (r *transactionRepository) FindByID(id uuid.UUID) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Where("id = ?", id).First(&transaction).Error
	return &transaction, err
}

// FindByUserID implements [TransactionRepository].
func (r *transactionRepository) FindByUserID(userID uuid.UUID) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("Items.Ticket").Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

// FindAll implements [TransactionRepository].
func (r *transactionRepository) FindAll() ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("User").Preload("Items.Ticket").Find(&transactions).Error
	return transactions, err
}

// FindExpiredPendingTransactions implements [TransactionRepository].
func (r *transactionRepository) FindExpiredPendingTransactions(expiredBefore time.Time) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Preload("Items").Preload("User").
		Where("payment_status = ? AND created_at < ?", enums.PaymentStatusPending, expiredBefore).
		Find(&transactions).Error
	return transactions, err
}

// UpdateStatus implements [TransactionRepository].
func (r *transactionRepository) UpdateStatus(id uuid.UUID, status enums.PaymentStatus) error {
	return r.db.Model(&model.Transaction{}).Where("id = ?", id).Update("payment_status", status).Error
}

// UpdateStatusAndMethod implements [TransactionRepository].
func (r *transactionRepository) UpdateStatusAndMethod(id uuid.UUID, status enums.PaymentStatus, method enums.PaymentMethod) error {
	return r.db.Model(&model.Transaction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"payment_status": status,
			"payment_method": method,
		}).Error
}

// Delete implements [TransactionRepository].
func (r *transactionRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Transaction{}).Error
}
