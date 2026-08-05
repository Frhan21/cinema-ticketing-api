package service

import (
	"cinema-ticketing-api/app/payment/dto"
	"cinema-ticketing-api/app/payment/gateway"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const midtransOrderIDPrefix = "cinema-"

type paymentService struct {
	paymentRepo     paymentRepository
	transactionRepo transactionRepository
	gateway         gateway.PaymentGateway
	mailer          paymentMailer
}

type paymentMailer interface {
	SendTicketPaidNotification(to, userName, movieTitle, scheduleTime string, totalPrice float64) error
}

type paymentRepository interface {
	Create(payload *entities.Payment) error
	FindByID(id uuid.UUID) (*entities.Payment, error)
	FindByTransactionID(transactionID uuid.UUID) (*entities.Payment, error)
	FindByOrderID(orderID string) (*entities.Payment, error)
	ApplyStatus(payment *entities.Payment, history *entities.PaymentHistory, ticketStatus entities.TicketStatus) (bool, error)
}

type transactionRepository interface {
	FindByID(id uuid.UUID) (*entities.Transaction, error)
}

type PaymentService interface {
	CreatePayment(userID uuid.UUID, transactionID uuid.UUID) (*dto.CreatePaymentResponse, error)
	FindByID(id uuid.UUID) (*entities.Payment, error)
	FindByTransactionID(transactionID uuid.UUID) (*entities.Payment, error)
	FindByOrderID(orderID string) (*entities.Payment, error)
	HandleNotification(req dto.MidtransNotificationRequest) error
}

func NewPaymentService(
	paymentRepo paymentRepository,
	transactionRepo transactionRepository,
	gateway gateway.PaymentGateway,
	mailer paymentMailer,
) PaymentService {
	return &paymentService{
		paymentRepo:     paymentRepo,
		transactionRepo: transactionRepo,
		gateway:         gateway,
		mailer:          mailer,
	}
}

func (s *paymentService) CreatePayment(userID uuid.UUID, transactionID uuid.UUID) (*dto.CreatePaymentResponse, error) {
	transaction, err := s.transactionRepo.FindByID(transactionID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewNotFoundError("transaction not found")
	}
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to get transaction")
	}
	if transaction.UserID != userID {
		return nil, apperror.NewForbiddenError("transaction does not belong to this user")
	}
	if transaction.PaymentStatus != entities.PaymentStatusPending {
		return nil, apperror.NewBadRequestError("transaction is already paid or cancelled")
	}

	existingPayment, err := s.paymentRepo.FindByTransactionID(transactionID)
	if err == nil {
		return newCreatePaymentResponse(existingPayment), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.NewInternalServerError("failed to check existing payment")
	}

	items, expectedGrossAmount, err := buildPaymentItems(transaction)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to calculate payment items")
	}

	orderID := midtransOrderIDPrefix + transaction.ID.String()
	gatewayInput := gateway.CreateTransactionInput{
		OrderID:         orderID,
		CustomerName:    transaction.User.Name,
		CustomerEmail:   transaction.User.Email,
		Items:           items,
		ExpiryStartTime: transaction.CreatedAt,
	}
	grossAmount, err := gatewayInput.GrossAmount()
	if err != nil || grossAmount != expectedGrossAmount {
		return nil, apperror.NewInternalServerError("payment item total does not match transaction total")
	}

	gatewayResult, err := s.gateway.CreateTransaction(gatewayInput)
	if err != nil {
		return nil, apperror.NewInternalServerError("failed to create Midtrans transaction")
	}
	if gatewayResult == nil || gatewayResult.Token == "" || gatewayResult.RedirectURL == "" {
		return nil, apperror.NewInternalServerError("Midtrans returned an incomplete payment response")
	}

	payment := &entities.Payment{
		ID:            uuid.New(),
		TransactionID: transaction.ID,
		OrderID:       orderID,
		SnapToken:     gatewayResult.Token,
		RedirectURL:   gatewayResult.RedirectURL,
		PaymentStatus: entities.PaymentStatusPending,
		GatewayStatus: string(entities.PaymentStatusPending),
		GrossAmount:   grossAmount,
	}
	if err := s.paymentRepo.Create(payment); err != nil {
		return nil, apperror.NewInternalServerError("failed to save payment")
	}

	paymentResponse := &dto.CreatePaymentResponse{
		OrderID:       payment.OrderID,
		Token:         payment.SnapToken,
		RedirectURL:   payment.RedirectURL,
		PaymentStatus: string(payment.PaymentStatus),
	}

	return paymentResponse, nil
}

func (s *paymentService) FindByID(id uuid.UUID) (*entities.Payment, error) {
	return s.paymentRepo.FindByID(id)
}

func (s *paymentService) FindByTransactionID(transactionID uuid.UUID) (*entities.Payment, error) {
	return s.paymentRepo.FindByTransactionID(transactionID)
}

func (s *paymentService) FindByOrderID(orderID string) (*entities.Payment, error) {
	return s.paymentRepo.FindByOrderID(orderID)
}

func (s *paymentService) HandleNotification(req dto.MidtransNotificationRequest) error {
	payment, err := s.paymentRepo.FindByOrderID(req.OrderID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperror.NewNotFoundError("payment not found")
	}
	if err != nil {
		return apperror.NewInternalServerError("failed to get payment")
	}

	verifiedStatus, err := s.gateway.CheckTransaction(req.OrderID)
	if err != nil {
		return apperror.NewInternalServerError("failed to verify Midtrans transaction")
	}
	if verifiedStatus == nil || verifiedStatus.OrderID != payment.OrderID {
		return apperror.NewBadRequestError("payment order ID mismatch")
	}

	verifiedAmount, err := parseIDRAmount(verifiedStatus.GrossAmount)
	if err != nil || verifiedAmount != payment.GrossAmount {
		return apperror.NewBadRequestError("payment amount mismatch")
	}

	paymentStatus, ticketStatus, paidAt, handled := mapGatewayStatus(verifiedStatus)
	if !handled {
		return nil
	}

	payment.PaymentStatus = paymentStatus
	payment.GatewayStatus = verifiedStatus.TransactionStatus
	if verifiedStatus.MidtransTransactionID != "" {
		payment.MidtransTransactionID = stringPointer(verifiedStatus.MidtransTransactionID)
	}
	if verifiedStatus.PaymentType != "" {
		payment.PaymentType = stringPointer(verifiedStatus.PaymentType)
	}
	if paidAt != nil && payment.PaidAt == nil {
		payment.PaidAt = paidAt
	}

	history := &entities.PaymentHistory{
		ID:            uuid.New(),
		PaymentID:     payment.ID,
		PaymentStatus: payment.PaymentStatus,
		GatewayStatus: payment.GatewayStatus,
	}
	statusChanged, err := s.paymentRepo.ApplyStatus(payment, history, ticketStatus)
	if err != nil {
		return apperror.NewInternalServerError("failed to update payment status")
	}
	if statusChanged && payment.PaymentStatus == entities.PaymentStatusPaid {
		s.sendPaidNotification(payment)
	}

	return nil
}

func newCreatePaymentResponse(payment *entities.Payment) *dto.CreatePaymentResponse {
	return &dto.CreatePaymentResponse{
		OrderID:       payment.OrderID,
		Token:         payment.SnapToken,
		RedirectURL:   payment.RedirectURL,
		PaymentStatus: string(payment.PaymentStatus),
	}
}

func buildPaymentItems(transaction *entities.Transaction) ([]gateway.ItemDetail, int64, error) {
	if len(transaction.Items) == 0 {
		return nil, 0, errors.New("transaction has no items")
	}

	items := make([]gateway.ItemDetail, 0, len(transaction.Items)+1)
	for _, transactionItem := range transaction.Items {
		price, err := wholeIDRAmount(transactionItem.Price)
		if err != nil {
			return nil, 0, err
		}

		grouped := false
		for i := range items {
			if items[i].Name == "Cinema Ticket" && items[i].Price == price {
				if items[i].Quantity == math.MaxInt32 {
					return nil, 0, errors.New("payment item quantity exceeds the supported range")
				}
				items[i].Quantity++
				grouped = true
				break
			}
		}
		if !grouped {
			items = append(items, gateway.ItemDetail{
				ID:       "cinema-ticket-" + strconv.Itoa(len(items)+1),
				Name:     "Cinema Ticket",
				Price:    price,
				Quantity: 1,
			})
		}
	}

	subtotal, err := (gateway.CreateTransactionInput{Items: items}).GrossAmount()
	if err != nil {
		return nil, 0, err
	}
	finalTotal, err := wholeIDRAmount(transaction.TotalPrice)
	if err != nil {
		return nil, 0, err
	}
	if finalTotal > subtotal {
		return nil, 0, errors.New("transaction total exceeds item subtotal")
	}
	if discountAmount := subtotal - finalTotal; discountAmount > 0 {
		items = append(items, gateway.ItemDetail{
			ID:       "promo-discount",
			Name:     "Promo Discount",
			Price:    -discountAmount,
			Quantity: 1,
		})
	}

	return items, finalTotal, nil
}

func wholeIDRAmount(value float64) (int64, error) {
	if value <= 0 || math.Trunc(value) != value || value > math.MaxInt64 {
		return 0, errors.New("amount must be a positive whole IDR value")
	}
	return int64(value), nil
}

func parseIDRAmount(amount string) (int64, error) {
	whole, fraction, hasFraction := strings.Cut(strings.TrimSpace(amount), ".")
	if hasFraction && strings.Trim(fraction, "0") != "" {
		return 0, errors.New("IDR amount contains a fractional value")
	}

	value, err := strconv.ParseInt(whole, 10, 64)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid IDR amount")
	}
	return value, nil
}

func mapGatewayStatus(status *gateway.TransactionStatus) (entities.PaymentStatus, entities.TicketStatus, *time.Time, bool) {
	switch status.TransactionStatus {
	case "settlement":
		return entities.PaymentStatusPaid, entities.TicketStatusPaid, settlementTime(status.SettlementTime), true
	case "capture":
		if status.FraudStatus == "accept" {
			return entities.PaymentStatusPaid, entities.TicketStatusPaid, settlementTime(status.SettlementTime), true
		}
		return entities.PaymentStatusPending, entities.TicketStatusPending, nil, true
	case "pending":
		return entities.PaymentStatusPending, entities.TicketStatusPending, nil, true
	case "deny", "cancel", "expire", "failure":
		return entities.PaymentStatusCancelled, entities.TicketStatusCancelled, nil, true
	default:
		return "", "", nil, false
	}
}

func settlementTime(value string) *time.Time {
	if value != "" {
		location := time.FixedZone("WIB", 7*60*60)
		if parsed, err := time.ParseInLocation("2006-01-02 15:04:05", value, location); err == nil {
			return &parsed
		}
	}

	now := time.Now()
	return &now
}

func stringPointer(value string) *string {
	return &value
}

func (s *paymentService) sendPaidNotification(payment *entities.Payment) {
	if s.mailer == nil || payment.Transaction.User.Email == "" || len(payment.Transaction.Items) == 0 {
		return
	}

	ticket := payment.Transaction.Items[0].Ticket
	movieTitle := ticket.Schedule.Movie.Title
	if movieTitle == "" {
		movieTitle = "Tiket Bioskop"
	}
	scheduleTime := ticket.Schedule.StartTime.
		In(time.FixedZone("WIB", 7*60*60)).
		Format("02 Jan 2006, 15:04 WIB")

	go func() {
		if err := s.mailer.SendTicketPaidNotification(
			payment.Transaction.User.Email,
			payment.Transaction.User.Name,
			movieTitle,
			scheduleTime,
			float64(payment.GrossAmount),
		); err != nil {
			log.Printf("[PaymentService] Failed to send paid email: %v", err)
		}
	}()
}
