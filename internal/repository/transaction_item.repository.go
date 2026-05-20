package repository

import (
	"cinema-ticketing-api/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionItemRepository interface {
	CreateMany(items []model.TransactionItem) error
	FindByTransactionId(transactionId uuid.UUID) ([]model.TransactionItem, error)
	FindTicketIdByTransactionId(transactionId uuid.UUID) ([]model.TransactionItem, error)
}

type transactionItemRepository struct {
	DB *gorm.DB
}

// CreateMany implements [TransactionItemRepository].
func (t *transactionItemRepository) CreateMany(items []model.TransactionItem) error {
	return t.DB.Create(&items).Error
}

// FindByTransactionId implements [TransactionItemRepository].
func (t *transactionItemRepository) FindByTransactionId(transactionId uuid.UUID) ([]model.TransactionItem, error) {
	var items []model.TransactionItem

	err := t.DB.Where("transaction_id = ?", transactionId).Find(&items).Error

	return items, err
}

// FindTicketIdByTransactionId implements [TransactionItemRepository].
func (t *transactionItemRepository) FindTicketIdByTransactionId(transactionId uuid.UUID) ([]model.TransactionItem, error) {
	var items []model.TransactionItem

	err := t.DB.Where("transaction_id = ?", transactionId).Pluck("ticket_id", &items).Error

	return items, err
}

func NewTransactionItemRepository(db *gorm.DB) TransactionItemRepository {
	return &transactionItemRepository{DB: db}
}
