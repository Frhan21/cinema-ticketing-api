package service

import (
	"cinema-ticketing-api/internal/enums"
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/pkg/mailer"
	"log"
	"time"

	"github.com/google/uuid"
)

type TransactionService interface {
	GetAll() ([]model.Transaction, error)
	GetByID(transactionID uuid.UUID) (*model.Transaction, error)
	PayTransaction(userID uuid.UUID, transactionID uuid.UUID, req request.PayTransactionRequest) error
	CancelTransaction(userID uuid.UUID, transactionID uuid.UUID) error
	AutoCancelExpiredTransactions()
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	ticketRepo      repository.TicketRepository
	mailer          *mailer.Mailer
}

func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	ticketRepo repository.TicketRepository,
	mailer *mailer.Mailer,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		ticketRepo:      ticketRepo,
		mailer:          mailer,
	}
}

// GetAll mengembalikan semua transaksi (untuk admin).
func (s *transactionService) GetAll() ([]model.Transaction, error) {
	return s.transactionRepo.FindAll()
}

// GetByID mengembalikan detail satu transaksi beserta items-nya.
func (s *transactionService) GetByID(transactionID uuid.UUID) (*model.Transaction, error) {
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return nil, apperror.NewNotFoundError("transaction not found")
	}
	return transaction, nil
}
// PayTransaction membayar transaksi yang statusnya masih pending.
// Setelah berhasil, status transaksi dan semua tiket terkait diubah menjadi "paid",
// lalu user mendapat email konfirmasi.
func (s *transactionService) PayTransaction(userID uuid.UUID, transactionID uuid.UUID, req request.PayTransactionRequest) error {
	// 1. Validasi transaksi ada + preload user & items
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if err != nil {
		return apperror.NewNotFoundError("transaction not found")
	}

	// 2. Validasi kepemilikan
	if transaction.UserID != userID {
		return apperror.NewUnauthorizedError("transaction does not belong to this user")
	}

	// 3. Validasi status masih pending
	if transaction.PaymentStatus != enums.PaymentStatusPending {
		return apperror.NewBadRequestError("transaction is already paid or cancelled")
	}

	// 4. Update status transaksi + payment method
	if err := s.transactionRepo.UpdateStatusAndMethod(transactionID, enums.PaymentStatusPaid, enums.PaymentMethod(req.PaymentMethod)); err != nil {
		return apperror.NewInternalServerError("failed to update transaction status")
	}

	// 5. Update semua tiket terkait menjadi paid
	var ticketIDs []uuid.UUID
	for _, item := range transaction.Items {
		ticketIDs = append(ticketIDs, item.TicketID)
	}
	if len(ticketIDs) > 0 {
		if err := s.ticketRepo.UpdateStatusByIDs(ticketIDs, enums.TicketStatusPaid); err != nil {
			return apperror.NewInternalServerError("failed to update ticket statuses")
		}
	}

	// 6. Kirim email konfirmasi (non-blocking)
	go func() {
		if err := s.mailer.SendTicketPaidNotification(
			transaction.User.Email,
			transaction.User.Name,
			"Tiket Bioskop",
			"",
			transaction.TotalPrice,
		); err != nil {
			log.Printf("[TransactionService] Failed to send paid email to %s: %v\n", transaction.User.Email, err)
		}
	}()

	return nil
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
	if transaction.PaymentStatus != enums.PaymentStatusPending {
		return apperror.NewBadRequestError("transaction is already paid or cancelled")
	}

	// 4. Update status transaksi menjadi cancelled
	if err := s.transactionRepo.UpdateStatus(transactionID, enums.PaymentStatusCancelled); err != nil {
		return apperror.NewInternalServerError("failed to cancel transaction")
	}

	// 5. Update semua tiket terkait menjadi cancelled
	var ticketIDs []uuid.UUID
	for _, item := range transaction.Items {
		ticketIDs = append(ticketIDs, item.TicketID)
	}
	if len(ticketIDs) > 0 {
		if err := s.ticketRepo.UpdateStatusByIDs(ticketIDs, enums.TicketStatusCancelled); err != nil {
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
	expiredBefore := time.Now().Add(-15 * time.Minute)

	expiredTransactions, err := s.transactionRepo.FindExpiredPendingTransactions(expiredBefore)
	if err != nil {
		log.Println("[TransactionService] Error fetching expired transactions:", err)
		return
	}

	for _, tx := range expiredTransactions {
		// Update status transaksi menjadi cancelled
		if err := s.transactionRepo.UpdateStatus(tx.ID, enums.PaymentStatusCancelled); err != nil {
			log.Printf("[TransactionService] Failed to cancel transaction %s: %v\n", tx.ID, err)
			continue
		}

		// Update semua tiket terkait menjadi cancelled
		var ticketIDs []uuid.UUID
		for _, item := range tx.Items {
			ticketIDs = append(ticketIDs, item.TicketID)
		}
		if len(ticketIDs) > 0 {
			s.ticketRepo.UpdateStatusByIDs(ticketIDs, enums.TicketStatusCancelled)
		}

		log.Printf("[TransactionService] Auto-cancelled transaction %s (expired 15 min)\n", tx.ID)

		// Kirim notifikasi email ke user (non-blocking)
		go func(user model.User) {
			if err := s.mailer.SendTicketCancelledNotification(
				user.Email, user.Name, "Tiket Anda",
			); err != nil {
				log.Printf("[TransactionService] Failed to send cancel email to %s: %v\n", user.Email, err)
			}
		}(tx.User)
	}
}
