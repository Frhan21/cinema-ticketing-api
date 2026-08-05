package gateway

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type ItemDetail struct {
	ID       string
	Name     string
	Price    int64
	Quantity int32
}

type CreateTransactionInput struct {
	OrderID         string
	CustomerName    string
	CustomerEmail   string
	Items           []ItemDetail
	ExpiryStartTime time.Time
}

func (i CreateTransactionInput) GrossAmount() (int64, error) {
	if len(i.Items) == 0 {
		return 0, fmt.Errorf("at least one payment item is required")
	}

	var total int64
	for _, item := range i.Items {
		if strings.TrimSpace(item.Name) == "" {
			return 0, fmt.Errorf("payment item name is required")
		}
		if item.Quantity <= 0 {
			return 0, fmt.Errorf("payment item quantity must be greater than zero")
		}
		if item.Price == 0 {
			return 0, fmt.Errorf("payment item price must not be zero")
		}

		quantity := int64(item.Quantity)
		if item.Price > 0 && item.Price > math.MaxInt64/quantity {
			return 0, fmt.Errorf("payment item amount exceeds the supported range")
		}
		if item.Price < 0 && item.Price < math.MinInt64/quantity {
			return 0, fmt.Errorf("payment item amount exceeds the supported range")
		}
		lineTotal := item.Price * quantity
		if lineTotal > 0 && total > math.MaxInt64-lineTotal {
			return 0, fmt.Errorf("payment gross amount exceeds the supported range")
		}
		if lineTotal < 0 && total < math.MinInt64-lineTotal {
			return 0, fmt.Errorf("payment gross amount exceeds the supported range")
		}
		total += lineTotal
	}

	if total <= 0 {
		return 0, fmt.Errorf("payment gross amount must be greater than zero")
	}
	return total, nil
}

type CreateTransactionResult struct {
	Token       string
	RedirectURL string
}

type TransactionStatus struct {
	OrderID               string
	MidtransTransactionID string
	GrossAmount           string
	PaymentType           string
	TransactionStatus     string
	FraudStatus           string
	StatusCode            string
	SettlementTime        string
}

type PaymentGateway interface {
	CreateTransaction(input CreateTransactionInput) (*CreateTransactionResult, error)
	CheckTransaction(orderID string) (*TransactionStatus, error)
}
