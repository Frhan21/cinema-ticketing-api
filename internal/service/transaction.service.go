package service

import (
	"cinema-ticketing-api/internal/enums"
	"cinema-ticketing-api/internal/model"
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/pkg/mailer"
	"errors"
	"log"

	"github.com/google/uuid"
)

type TransactionService interface {
	GetAll() ([]model.Transaction, error)
	GetByID(transactionID uuid.UUID) (*model.Transaction, error)
	PayTransaction(userID uuid.UUID, transactionID uuid.UUID, req request.PayTransactionRequest) error
	CancelTransaction(userID uuid.UUID, transactionID uuid.UUID) error
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
		return nil, errors.New("transaction not found")
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
		return errors.New("transaction not found")
	}

	// 2. Validasi kepemilikan
	if transaction.UserID != userID {
		return errors.New("transaction does not belong to this user")
	}

	// 3. Validasi status masih pending
	if transaction.PaymentStatus != enums.PaymentStatusPending {
		return errors.New("transaction is already paid or cancelled")
	}

	// 4. Update status transaksi + payment method
	if err := s.transactionRepo.UpdateStatusAndMethod(transactionID, enums.PaymentStatusPaid, enums.PaymentMethod(req.PaymentMethod)); err != nil {
		return errors.New("failed to update transaction status")
	}

	// 5. Update semua tiket terkait menjadi paid
	var ticketIDs []uuid.UUID
	for _, item := range transaction.Items {
		ticketIDs = append(ticketIDs, item.TicketID)
	}
	if len(ticketIDs) > 0 {
		if err := s.ticketRepo.UpdateStatusByIDs(ticketIDs, enums.TicketStatusPaid); err != nil {
			return errors.New("failed to update ticket statuses")
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
		return errors.New("transaction not found")
	}

	// 2. Validasi kepemilikan
	if transaction.UserID != userID {
		return errors.New("transaction does not belong to this user")
	}

	// 3. Validasi status masih pending
	if transaction.PaymentStatus != enums.PaymentStatusPending {
		return errors.New("transaction is already paid or cancelled")
	}

	// 4. Update status transaksi menjadi cancelled
	if err := s.transactionRepo.UpdateStatus(transactionID, enums.PaymentStatusCancelled); err != nil {
		return errors.New("failed to cancel transaction")
	}

	// 5. Update semua tiket terkait menjadi cancelled
	var ticketIDs []uuid.UUID
	for _, item := range transaction.Items {
		ticketIDs = append(ticketIDs, item.TicketID)
	}
	if len(ticketIDs) > 0 {
		if err := s.ticketRepo.UpdateStatusByIDs(ticketIDs, enums.TicketStatusCancelled); err != nil {
			return errors.New("failed to update ticket statuses")
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
