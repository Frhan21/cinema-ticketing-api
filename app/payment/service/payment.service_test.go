package service

import (
	paymentdto "cinema-ticketing-api/app/payment/dto"
	paymentgateway "cinema-ticketing-api/app/payment/gateway"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentRepositoryStub struct {
	payment              *entities.Payment
	findByTransactionErr error
	findByOrderErr       error
	createdPayment       *entities.Payment
	appliedPayment       *entities.Payment
	appliedHistory       *entities.PaymentHistory
	appliedTicketStatus  entities.TicketStatus
	applyCalls           int
	applyErr             error
}

func (r *paymentRepositoryStub) Create(payment *entities.Payment) error {
	r.createdPayment = payment
	return nil
}

func (r *paymentRepositoryStub) FindByID(uuid.UUID) (*entities.Payment, error) {
	return r.payment, nil
}

func (r *paymentRepositoryStub) FindByTransactionID(uuid.UUID) (*entities.Payment, error) {
	return r.payment, r.findByTransactionErr
}

func (r *paymentRepositoryStub) FindByOrderID(string) (*entities.Payment, error) {
	return r.payment, r.findByOrderErr
}

func (r *paymentRepositoryStub) ApplyStatus(payment *entities.Payment, history *entities.PaymentHistory, ticketStatus entities.TicketStatus) (bool, error) {
	r.applyCalls++
	r.appliedPayment = payment
	r.appliedHistory = history
	r.appliedTicketStatus = ticketStatus
	return true, r.applyErr
}

type transactionRepositoryStub struct {
	transaction *entities.Transaction
	err         error
}

func (r *transactionRepositoryStub) FindByID(uuid.UUID) (*entities.Transaction, error) {
	return r.transaction, r.err
}

type paymentGatewayStub struct {
	createCalls int
	input       paymentgateway.CreateTransactionInput
	result      *paymentgateway.CreateTransactionResult
	err         error
	checkOrder  string
	checkResult *paymentgateway.TransactionStatus
	checkErr    error
}

type paymentMailerStub struct {
	sent chan struct{}
}

func (m *paymentMailerStub) SendTicketPaidNotification(string, string, string, string, float64) error {
	m.sent <- struct{}{}
	return nil
}

func (g *paymentGatewayStub) CreateTransaction(input paymentgateway.CreateTransactionInput) (*paymentgateway.CreateTransactionResult, error) {
	g.createCalls++
	g.input = input
	return g.result, g.err
}

func (g *paymentGatewayStub) CheckTransaction(orderID string) (*paymentgateway.TransactionStatus, error) {
	g.checkOrder = orderID
	return g.checkResult, g.checkErr
}

func TestCreatePaymentCreatesSnapTransaction(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	paymentRepo := &paymentRepositoryStub{findByTransactionErr: gorm.ErrRecordNotFound}
	transactionRepo := &transactionRepositoryStub{transaction: &entities.Transaction{
		ID:            transactionID,
		UserID:        userID,
		TotalPrice:    50000,
		PaymentStatus: entities.PaymentStatusPending,
		User: entities.User{
			Name:  "Cinema User",
			Email: "user@example.com",
		},
		Items: transactionItems(25000, 25000),
	}}
	gateway := &paymentGatewayStub{result: &paymentgateway.CreateTransactionResult{
		Token:       "snap-token",
		RedirectURL: "https://app.sandbox.midtrans.com/snap/v3/redirection/snap-token",
	}}
	service := NewPaymentService(paymentRepo, transactionRepo, gateway, nil)

	response, err := service.CreatePayment(userID, transactionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantOrderID := midtransOrderIDPrefix + transactionID.String()
	if gateway.createCalls != 1 {
		t.Fatalf("expected one gateway call, got %d", gateway.createCalls)
	}
	grossAmount, grossErr := gateway.input.GrossAmount()
	if gateway.input.OrderID != wantOrderID || grossErr != nil || grossAmount != 50000 {
		t.Fatalf("unexpected gateway input: %+v", gateway.input)
	}
	if len(gateway.input.Items) != 1 || gateway.input.Items[0].Price != 25000 || gateway.input.Items[0].Quantity != 2 {
		t.Fatalf("unexpected payment items: %+v", gateway.input.Items)
	}
	if paymentRepo.createdPayment == nil {
		t.Fatal("expected payment to be persisted")
	}
	if paymentRepo.createdPayment.OrderID != wantOrderID {
		t.Fatalf("expected order ID %q, got %q", wantOrderID, paymentRepo.createdPayment.OrderID)
	}
	if response.Token != "snap-token" || response.RedirectURL != gateway.result.RedirectURL {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestCreatePaymentReturnsExistingPayment(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	paymentRepo := &paymentRepositoryStub{payment: &entities.Payment{
		TransactionID: transactionID,
		OrderID:       midtransOrderIDPrefix + transactionID.String(),
		SnapToken:     "existing-token",
		RedirectURL:   "https://example.com/existing-token",
		PaymentStatus: entities.PaymentStatusPending,
	}}
	transactionRepo := &transactionRepositoryStub{transaction: &entities.Transaction{
		ID:            transactionID,
		UserID:        userID,
		TotalPrice:    50000,
		PaymentStatus: entities.PaymentStatusPending,
	}}
	gateway := &paymentGatewayStub{}
	service := NewPaymentService(paymentRepo, transactionRepo, gateway, nil)

	response, err := service.CreatePayment(userID, transactionID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gateway.createCalls != 0 {
		t.Fatalf("expected no gateway call, got %d", gateway.createCalls)
	}
	if response.Token != "existing-token" {
		t.Fatalf("expected existing token, got %q", response.Token)
	}
}

func TestCreatePaymentRejectsOtherUserTransaction(t *testing.T) {
	transactionID := uuid.New()
	service := NewPaymentService(
		&paymentRepositoryStub{},
		&transactionRepositoryStub{transaction: &entities.Transaction{
			ID:            transactionID,
			UserID:        uuid.New(),
			PaymentStatus: entities.PaymentStatusPending,
		}},
		&paymentGatewayStub{},
		nil,
	)

	_, err := service.CreatePayment(uuid.New(), transactionID)
	assertAppErrorCode(t, err, http.StatusForbidden)
}

func TestCreatePaymentHandlesGatewayFailure(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	service := NewPaymentService(
		&paymentRepositoryStub{findByTransactionErr: gorm.ErrRecordNotFound},
		&transactionRepositoryStub{transaction: &entities.Transaction{
			ID:            transactionID,
			UserID:        userID,
			TotalPrice:    50000,
			PaymentStatus: entities.PaymentStatusPending,
			Items:         transactionItems(25000, 25000),
		}},
		&paymentGatewayStub{err: errors.New("Midtrans unavailable")},
		nil,
	)

	_, err := service.CreatePayment(userID, transactionID)
	assertAppErrorCode(t, err, http.StatusInternalServerError)
}

func TestCreatePaymentAddsPromoDiscountItem(t *testing.T) {
	userID := uuid.New()
	transactionID := uuid.New()
	paymentRepo := &paymentRepositoryStub{findByTransactionErr: gorm.ErrRecordNotFound}
	transactionRepo := &transactionRepositoryStub{transaction: &entities.Transaction{
		ID:            transactionID,
		UserID:        userID,
		TotalPrice:    80000,
		PaymentStatus: entities.PaymentStatusPending,
		Items:         transactionItems(50000, 50000),
	}}
	gateway := &paymentGatewayStub{result: &paymentgateway.CreateTransactionResult{
		Token:       "snap-token",
		RedirectURL: "https://example.com/snap-token",
	}}
	service := NewPaymentService(paymentRepo, transactionRepo, gateway, nil)

	if _, err := service.CreatePayment(userID, transactionID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gateway.input.Items) != 2 {
		t.Fatalf("expected ticket and discount items, got %+v", gateway.input.Items)
	}
	discount := gateway.input.Items[1]
	if discount.Name != "Promo Discount" || discount.Price != -20000 || discount.Quantity != 1 {
		t.Fatalf("unexpected discount item: %+v", discount)
	}
	grossAmount, err := gateway.input.GrossAmount()
	if err != nil || grossAmount != 80000 {
		t.Fatalf("expected gross amount 80000, got %d (err=%v)", grossAmount, err)
	}
	if paymentRepo.createdPayment.GrossAmount != 80000 {
		t.Fatalf("expected persisted gross amount 80000, got %d", paymentRepo.createdPayment.GrossAmount)
	}
}

func TestHandleNotificationAppliesSettlement(t *testing.T) {
	paymentID := uuid.New()
	transactionID := uuid.New()
	orderID := midtransOrderIDPrefix + transactionID.String()
	paymentRepo := &paymentRepositoryStub{payment: &entities.Payment{
		ID:            paymentID,
		TransactionID: transactionID,
		OrderID:       orderID,
		PaymentStatus: entities.PaymentStatusPending,
		GatewayStatus: "pending",
		GrossAmount:   50000,
		Transaction: entities.Transaction{
			User: entities.User{Name: "Cinema User", Email: "user@example.com"},
			Items: []entities.TransactionItem{{
				Ticket: entities.Ticket{Schedule: entities.Schedule{
					StartTime: time.Now().Add(time.Hour),
					Movie:     entities.Movie{Title: "Test Movie"},
				}},
			}},
		},
	}}
	gateway := &paymentGatewayStub{checkResult: &paymentgateway.TransactionStatus{
		OrderID:               orderID,
		MidtransTransactionID: "midtrans-transaction-id",
		GrossAmount:           "50000.00",
		PaymentType:           "bank_transfer",
		TransactionStatus:     "settlement",
		SettlementTime:        "2026-08-05 10:16:00",
	}}
	mailer := &paymentMailerStub{sent: make(chan struct{}, 1)}
	service := NewPaymentService(paymentRepo, &transactionRepositoryStub{}, gateway, mailer)

	err := service.HandleNotification(paymentdto.MidtransNotificationRequest{OrderID: orderID})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gateway.checkOrder != orderID {
		t.Fatalf("expected status check for %q, got %q", orderID, gateway.checkOrder)
	}
	if paymentRepo.applyCalls != 1 {
		t.Fatalf("expected one status update, got %d", paymentRepo.applyCalls)
	}
	if paymentRepo.appliedPayment.PaymentStatus != entities.PaymentStatusPaid {
		t.Fatalf("expected paid status, got %q", paymentRepo.appliedPayment.PaymentStatus)
	}
	if paymentRepo.appliedPayment.GatewayStatus != "settlement" {
		t.Fatalf("expected settlement gateway status, got %q", paymentRepo.appliedPayment.GatewayStatus)
	}
	if paymentRepo.appliedTicketStatus != entities.TicketStatusPaid {
		t.Fatalf("expected paid ticket status, got %q", paymentRepo.appliedTicketStatus)
	}
	if paymentRepo.appliedPayment.PaidAt == nil {
		t.Fatal("expected paid_at to be populated")
	}
	if paymentRepo.appliedHistory.PaymentID != paymentID {
		t.Fatalf("expected history payment ID %s, got %s", paymentID, paymentRepo.appliedHistory.PaymentID)
	}
	select {
	case <-mailer.sent:
	case <-time.After(time.Second):
		t.Fatal("expected paid notification email")
	}
}

func TestHandleNotificationRejectsAmountMismatch(t *testing.T) {
	orderID := midtransOrderIDPrefix + uuid.NewString()
	paymentRepo := &paymentRepositoryStub{payment: &entities.Payment{
		OrderID:     orderID,
		GrossAmount: 50000,
	}}
	gateway := &paymentGatewayStub{checkResult: &paymentgateway.TransactionStatus{
		OrderID:           orderID,
		GrossAmount:       "40000.00",
		TransactionStatus: "settlement",
	}}
	service := NewPaymentService(paymentRepo, &transactionRepositoryStub{}, gateway, nil)

	err := service.HandleNotification(paymentdto.MidtransNotificationRequest{OrderID: orderID})
	assertAppErrorCode(t, err, http.StatusBadRequest)
	if paymentRepo.applyCalls != 0 {
		t.Fatalf("expected no status update, got %d", paymentRepo.applyCalls)
	}
}

func TestHandleNotificationIgnoresUnsupportedStatus(t *testing.T) {
	orderID := midtransOrderIDPrefix + uuid.NewString()
	paymentRepo := &paymentRepositoryStub{payment: &entities.Payment{
		OrderID:     orderID,
		GrossAmount: 50000,
	}}
	gateway := &paymentGatewayStub{checkResult: &paymentgateway.TransactionStatus{
		OrderID:           orderID,
		GrossAmount:       "50000.00",
		TransactionStatus: "refund",
	}}
	service := NewPaymentService(paymentRepo, &transactionRepositoryStub{}, gateway, nil)

	if err := service.HandleNotification(paymentdto.MidtransNotificationRequest{OrderID: orderID}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if paymentRepo.applyCalls != 0 {
		t.Fatalf("expected no status update, got %d", paymentRepo.applyCalls)
	}
}

func assertAppErrorCode(t *testing.T, err error, wantCode int) {
	t.Helper()
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != wantCode {
		t.Fatalf("expected status %d, got %d", wantCode, appErr.Code)
	}
}

func transactionItems(prices ...float64) []entities.TransactionItem {
	items := make([]entities.TransactionItem, 0, len(prices))
	for _, price := range prices {
		items = append(items, entities.TransactionItem{Price: price})
	}
	return items
}
