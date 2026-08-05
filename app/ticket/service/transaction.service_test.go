package service

import (
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

type transactionStatusRepositoryStub struct {
	transaction *entities.Transaction
}

func (r *transactionStatusRepositoryStub) FindAll() ([]entities.Transaction, error) {
	return nil, nil
}

func (r *transactionStatusRepositoryStub) FindByID(uuid.UUID) (*entities.Transaction, error) {
	return r.transaction, nil
}

func (r *transactionStatusRepositoryStub) FindExpiredPendingTransactions(time.Time) ([]entities.Transaction, error) {
	return nil, nil
}

func (r *transactionStatusRepositoryStub) UpdateStatus(uuid.UUID, entities.PaymentStatus) error {
	return nil
}

type ticketStatusRepositoryStub struct{}

func (r *ticketStatusRepositoryStub) UpdateStatusByIDs([]uuid.UUID, entities.TicketStatus) error {
	return nil
}

func TestGetTransactionByIDEnforcesOwnership(t *testing.T) {
	ownerID := uuid.New()
	service := NewTransactionService(
		&transactionStatusRepositoryStub{transaction: &entities.Transaction{ID: uuid.New(), UserID: ownerID}},
		&ticketStatusRepositoryStub{},
		nil,
		15*time.Minute,
	)

	_, err := service.GetByID(uuid.New(), string(entities.RoleUser), uuid.New())
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) || appErr.Code != http.StatusForbidden {
		t.Fatalf("expected forbidden error, got %v", err)
	}

	if _, err := service.GetByID(uuid.New(), string(entities.RoleAdmin), uuid.New()); err != nil {
		t.Fatalf("expected admin access, got %v", err)
	}
	if _, err := service.GetByID(ownerID, string(entities.RoleUser), uuid.New()); err != nil {
		t.Fatalf("expected owner access, got %v", err)
	}
}
