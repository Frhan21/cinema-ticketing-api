package gateway

import (
	"cinema-ticketing-api/config"
	"fmt"
	"strings"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type midtransGateway struct {
	snapClient    snap.Client
	coreClient    coreapi.Client
	expiryMinutes int64
}

func NewMidtransGateway(cfg config.MidtransConfig) (PaymentGateway, error) {
	serverKey := strings.TrimSpace(cfg.ServerKey)
	if serverKey == "" {
		return nil, fmt.Errorf("midtrans server key is required")
	}
	if cfg.ExpiryMinutes <= 0 {
		return nil, fmt.Errorf("payment expiry minutes must be greater than zero")
	}

	var environment midtrans.EnvironmentType
	switch strings.ToLower(strings.TrimSpace(cfg.Environment)) {
	case "sandbox":
		environment = midtrans.Sandbox
	case "production":
		environment = midtrans.Production
	default:
		return nil, fmt.Errorf("invalid Midtrans environment %q", cfg.Environment)
	}

	var snapClient snap.Client
	snapClient.New(serverKey, environment)

	var coreClient coreapi.Client
	coreClient.New(serverKey, environment)

	return &midtransGateway{
		snapClient:    snapClient,
		coreClient:    coreClient,
		expiryMinutes: int64(cfg.ExpiryMinutes),
	}, nil
}

func (g *midtransGateway) CreateTransaction(input CreateTransactionInput) (*CreateTransactionResult, error) {
	if strings.TrimSpace(input.OrderID) == "" {
		return nil, fmt.Errorf("Midtrans order ID is required")
	}

	grossAmount, err := input.GrossAmount()
	if err != nil {
		return nil, err
	}

	items := make([]midtrans.ItemDetails, 0, len(input.Items))
	for _, item := range input.Items {
		items = append(items, midtrans.ItemDetails{
			ID:    item.ID,
			Name:  item.Name,
			Price: item.Price,
			Qty:   item.Quantity,
		})
	}
	expiryStartTime := input.ExpiryStartTime
	if expiryStartTime.IsZero() {
		expiryStartTime = time.Now()
	}
	expiryStartTime = expiryStartTime.In(time.FixedZone("WIB", 7*60*60))

	response, midtransErr := g.snapClient.CreateTransaction(&snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  input.OrderID,
			GrossAmt: grossAmount,
		},
		Items: &items,
		CustomerDetail: &midtrans.CustomerDetails{
			FName: input.CustomerName,
			Email: input.CustomerEmail,
		},
		CreditCard: &snap.CreditCardDetails{Secure: true},
		Expiry: &snap.ExpiryDetails{
			StartTime: expiryStartTime.Format("2006-01-02 15:04:05 -0700"),
			Unit:      "minute",
			Duration:  g.expiryMinutes,
		},
	})
	if midtransErr != nil {
		return nil, midtransErr
	}

	return &CreateTransactionResult{
		Token:       response.Token,
		RedirectURL: response.RedirectURL,
	}, nil
}

func (g *midtransGateway) CheckTransaction(orderID string) (*TransactionStatus, error) {
	if strings.TrimSpace(orderID) == "" {
		return nil, fmt.Errorf("Midtrans order ID is required")
	}

	response, err := g.coreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, err
	}

	return &TransactionStatus{
		OrderID:               response.OrderID,
		MidtransTransactionID: response.TransactionID,
		GrossAmount:           response.GrossAmount,
		PaymentType:           response.PaymentType,
		TransactionStatus:     response.TransactionStatus,
		FraudStatus:           response.FraudStatus,
		StatusCode:            response.StatusCode,
		SettlementTime:        response.SettlementTime,
	}, nil
}
