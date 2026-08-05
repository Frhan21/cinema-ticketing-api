package repository

import (
	"cinema-ticketing-api/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	FindByID(id uuid.UUID) (*entities.Transaction, error)
	FindByUserID(userID uuid.UUID) ([]entities.Transaction, error)
	FindAll() ([]entities.Transaction, error)
	FindExpiredPendingTransactions(expiredBefore time.Time) ([]entities.Transaction, error)
	UpdateStatus(id uuid.UUID, status entities.PaymentStatus) error
	Delete(id uuid.UUID) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

// FindByID implements [TransactionRepository].
func (r *transactionRepository) FindByID(id uuid.UUID) (*entities.Transaction, error) {
	var transaction entities.Transaction
	err := r.db.Preload("User").
		Preload("Items.Ticket.Schedule.Movie").
		Preload("Items.Ticket.Seat").
		Where("id = ?", id).
		First(&transaction).Error
	return &transaction, err
}

// FindByUserID implements [TransactionRepository].
func (r *transactionRepository) FindByUserID(userID uuid.UUID) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := r.db.Preload("Items.Ticket").Where("user_id = ?", userID).Find(&transactions).Error
	return transactions, err
}

// FindAll implements [TransactionRepository].
func (r *transactionRepository) FindAll() ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := r.db.Preload("User").Preload("Items.Ticket").Find(&transactions).Error
	return transactions, err
}

// FindExpiredPendingTransactions implements [TransactionRepository].
func (r *transactionRepository) FindExpiredPendingTransactions(expiredBefore time.Time) ([]entities.Transaction, error) {
	var transactions []entities.Transaction
	err := r.db.Preload("Items").Preload("User").
		Where("payment_status = ? AND created_at < ?", entities.PaymentStatusPending, expiredBefore).
		Find(&transactions).Error
	return transactions, err
}

// UpdateStatus implements [TransactionRepository].
func (r *transactionRepository) UpdateStatus(id uuid.UUID, status entities.PaymentStatus) error {
	return r.db.Model(&entities.Transaction{}).Where("id = ?", id).Update("payment_status", status).Error
}

// Delete implements [TransactionRepository].
func (r *transactionRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&entities.Transaction{}).Error
}
