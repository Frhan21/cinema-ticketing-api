package handler

import (
	"cinema-ticketing-api/app/payment/dto"
	"cinema-ticketing-api/entities"
	"cinema-ticketing-api/pkg/apperror"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type paymentServiceStub struct {
	userID        uuid.UUID
	transactionID uuid.UUID
	response      *dto.CreatePaymentResponse
	err           error
}

func (s *paymentServiceStub) CreatePayment(userID uuid.UUID, transactionID uuid.UUID) (*dto.CreatePaymentResponse, error) {
	s.userID = userID
	s.transactionID = transactionID
	return s.response, s.err
}

func (s *paymentServiceStub) HandleNotification(dto.MidtransNotificationRequest) error {
	return nil
}

func (s *paymentServiceStub) FindByID(uuid.UUID) (*entities.Payment, error) {
	return nil, nil
}

func (s *paymentServiceStub) FindByTransactionID(uuid.UUID) (*entities.Payment, error) {
	return nil, nil
}

func (s *paymentServiceStub) FindByOrderID(string) (*entities.Payment, error) {
	return nil, nil
}

func TestPaymentControllerCreatePayment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	transactionID := uuid.New()
	service := &paymentServiceStub{response: &dto.CreatePaymentResponse{
		OrderID:       "cinema-" + transactionID.String(),
		Token:         "snap-token",
		RedirectURL:   "https://app.sandbox.midtrans.com/snap/v3/redirection/snap-token",
		PaymentStatus: "pending",
	}}
	controller := NewPaymentController(service)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/v1/transaction/"+transactionID.String()+"/pay", nil)
	context.Params = gin.Params{{Key: "transaction_id", Value: transactionID.String()}}
	context.Set("user_id", userID.String())

	controller.CreatePayment(context)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if service.userID != userID || service.transactionID != transactionID {
		t.Fatalf("unexpected service arguments: user=%s transaction=%s", service.userID, service.transactionID)
	}
	if !strings.Contains(recorder.Body.String(), service.response.RedirectURL) {
		t.Fatalf("expected redirect URL in response, got %s", recorder.Body.String())
	}
}

func TestPaymentControllerCreatePaymentRejectsInvalidTransactionID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &paymentServiceStub{}
	controller := NewPaymentController(service)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest(http.MethodPost, "/api/v1/transaction/invalid/pay", nil)
	context.Params = gin.Params{{Key: "transaction_id", Value: "invalid"}}
	context.Set("user_id", uuid.NewString())

	controller.CreatePayment(context)

	if len(context.Errors) != 1 {
		t.Fatalf("expected one context error, got %d", len(context.Errors))
	}
	appErr, ok := context.Errors[0].Err.(*apperror.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", context.Errors[0].Err)
	}
	if appErr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, appErr.Code)
	}
}
