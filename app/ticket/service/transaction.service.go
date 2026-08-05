package service

import (
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/mailer"
	"log"
	"time"

	"github.com/google/uuid"
)

type TransactionService interface {
	GetAll() ([]entities.Transaction, error)
	GetByID(userID uuid.UUID, role string, transactionID uuid.UUID) (*entities.Transaction, error)
	CancelTransaction(userID uuid.UUID, transactionID uuid.UUID) error
	AutoCancelExpiredTransactions()
}

type transactionStatusRepository interface {
	FindAll() ([]entities.Transaction, error)
	FindByID(id uuid.UUID) (*entities.Transaction, error)
	FindExpiredPendingTransactions(expiredBefore time.Time) ([]entities.Transaction, error)
	UpdateStatus(id uuid.UUID, status entities.PaymentStatus) error
}

type ticketStatusRepository interface {
	UpdateStatusByIDs(ticketIDs []uuid.UUID, status entities.TicketStatus) error
}

type transactionService struct {
	transactionRepo transactionStatusRepository
	ticketRepo      ticketStatusRepository
	mailer          *mailer.Mailer
	paymentExpiry   time.Duration
}

func NewTransactionService(
	transactionRepo transactionStatusRepository,
	ticketRepo ticketStatusRepository,
	mailer *mailer.Mailer,
	paymentExpiry time.Duration,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		ticketRepo:      ticketRepo,
		mailer:          mailer,
		paymentExpiry:   paymentExpiry,
	}
}

// GetAll mengembalikan semua transaksi (untuk admin).
func (s *transactionService) GetAll() ([]entities.Transaction, error) {
	return s.transactionRepo.FindAll()
}

// GetByID mengembalikan detail satu transaksi beserta items-nya.
func (s *transactionService) GetByID(userID uuid.UUID, role string, transactionID uuid.UUID) (*entities.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return nil, apperror.NewNotFoundError("transaction not found")
	}
	if transaction.UserID != userID && role != string(entities.RoleAdmin) {
		return nil, apperror.NewForbiddenError("transaction does not belong to this user")
	}
	return transaction, nil
}

// CancelTransaction membatalkan transaksi yang statusnya masih pending.
// Setelah berhasil, status transaksi dan semua tiket terkait diubah menjadi "cancelled",
// lalu user mendapat email notifikasi.
func (s *transactionService) CancelTransaction(userID uuid.UUID, transactionID uuid.UUID) error {
	// 1. Validasi transaksi ada
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return apperror.NewNotFoundError("transaction not found")
	}

	// 2. Validasi kepemilikan
	if transaction.UserID != userID {
		return apperror.NewUnauthorizedError("transaction does not belong to this user")
	}

	// 3. Validasi status masih pending
	if transaction.PaymentStatus != entities.PaymentStatusPending {
		return apperror.NewBadRequestError("transaction is already paid or cancelled")
	}

	// 4. Update status transaksi menjadi cancelled
	if err := s.transactionRepo.UpdateStatus(transactionID, entities.PaymentStatusCancelled); err != nil {
		return apperror.NewInternalServerError("failed to cancel transaction")
	}

	// 5. Update semua tiket terkait menjadi cancelled
	var ticketIDs []uuid.UUID
	for _, item := range transaction.Items {
		ticketIDs = append(ticketIDs, item.TicketID)
	}
	if len(ticketIDs) > 0 {
		if err := s.ticketRepo.UpdateStatusByIDs(ticketIDs, entities.TicketStatusCancelled); err != nil {
			return apperror.NewInternalServerError("failed to update ticket statuses")
		}
	}

	// 6. Kirim email notifikasi (non-blocking)
	go func() {
		if err := s.mailer.SendTicketCancelledNotification(
			transaction.User.Email,
			transaction.User.Name,
			"Tiket Bioskop",
		); err != nil {
			log.Printf("[TransactionService] Failed to send cancel email to %s: %v\n", transaction.User.Email, err)
		}
	}()

	return nil
}

// AutoCancelExpiredTransactions implements [TransactionService].
func (s *transactionService) AutoCancelExpiredTransactions() {
	expiredBefore := time.Now().Add(-s.paymentExpiry)

	expiredTransactions, err := s.transactionRepo.FindExpiredPendingTransactions(expiredBefore)
	if err != nil {
		log.Println("[TransactionService] Error fetching expired transactions:", err)
		return
	}

	for _, tx := range expiredTransactions {
		// Update status transaksi menjadi cancelled
		if err := s.transactionRepo.UpdateStatus(tx.ID, entities.PaymentStatusCancelled); err != nil {
			log.Printf("[TransactionService] Failed to cancel transaction %s: %v\n", tx.ID, err)
			continue
		}

		// Update semua tiket terkait menjadi cancelled
		var ticketIDs []uuid.UUID
		for _, item := range tx.Items {
			ticketIDs = append(ticketIDs, item.TicketID)
		}
		if len(ticketIDs) > 0 {
			s.ticketRepo.UpdateStatusByIDs(ticketIDs, entities.TicketStatusCancelled)
		}

		log.Printf("[TransactionService] Auto-cancelled transaction %s (expired after %s)\n", tx.ID, s.paymentExpiry)

		// Kirim notifikasi email ke user (non-blocking)
		go func(user entities.User) {
			if err := s.mailer.SendTicketCancelledNotification(
				user.Email, user.Name, "Tiket Anda",
			); err != nil {
				log.Printf("[TransactionService] Failed to send cancel email to %s: %v\n", user.Email, err)
			}
		}(tx.User)
	}
}
