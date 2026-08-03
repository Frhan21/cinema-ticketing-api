package ticket

import (
	"cinema-ticketing-api/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionItemRepository interface {
	CreateMany(items []entities.TransactionItem) error
	FindByTransactionId(transactionId uuid.UUID) ([]entities.TransactionItem, error)
	FindTicketIdByTransactionId(transactionId uuid.UUID) ([]entities.TransactionItem, error)
}

type transactionItemRepository struct {
	DB *gorm.DB
}

// CreateMany implements [TransactionItemRepository].
func (t *transactionItemRepository) CreateMany(items []entities.TransactionItem) error {
	return t.DB.Create(&items).Error
}

// FindByTransactionId implements [TransactionItemRepository].
func (t *transactionItemRepository) FindByTransactionId(transactionId uuid.UUID) ([]entities.TransactionItem, error) {
	var items []entities.TransactionItem

	err := t.DB.Where("transaction_id = ?", transactionId).Find(&items).Error

	return items, err
}

// FindTicketIdByTransactionId implements [TransactionItemRepository].
func (t *transactionItemRepository) FindTicketIdByTransactionId(transactionId uuid.UUID) ([]entities.TransactionItem, error) {
	var items []entities.TransactionItem

	err := t.DB.Where("transaction_id = ?", transactionId).Pluck("ticket_id", &items).Error

	return items, err
}

func NewTransactionItemRepository(db *gorm.DB) TransactionItemRepository {
	return &transactionItemRepository{DB: db}
}
