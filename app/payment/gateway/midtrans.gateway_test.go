package gateway

import (
	"cinema-ticketing-api/config"
	"testing"

	"github.com/midtrans/midtrans-go"
)

func TestNewMidtransGateway(t *testing.T) {
	tests := []struct {
		name    string
		config  config.MidtransConfig
		wantEnv midtrans.EnvironmentType
		wantErr bool
	}{
		{
			name: "sandbox",
			config: config.MidtransConfig{
				ServerKey:     "sandbox-server-key",
				Environment:   "sandbox",
				ExpiryMinutes: 15,
			},
			wantEnv: midtrans.Sandbox,
		},
		{
			name: "production",
			config: config.MidtransConfig{
				ServerKey:     "production-server-key",
				Environment:   "production",
				ExpiryMinutes: 15,
			},
			wantEnv: midtrans.Production,
		},
		{
			name: "missing server key",
			config: config.MidtransConfig{
				Environment: "sandbox",
			},
			wantErr: true,
		},
		{
			name: "invalid environment",
			config: config.MidtransConfig{
				ServerKey:     "server-key",
				Environment:   "staging",
				ExpiryMinutes: 15,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gateway, err := NewMidtransGateway(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			client := gateway.(*midtransGateway)
			if client.snapClient.Env != tt.wantEnv {
				t.Fatalf("expected environment %v, got %v", tt.wantEnv, client.snapClient.Env)
			}
			if client.coreClient.Env != tt.wantEnv {
				t.Fatalf("expected environment %v, got %v", tt.wantEnv, client.coreClient.Env)
			}
		})
	}
}

func TestMidtransGatewayRejectsInvalidCreateInput(t *testing.T) {
	gateway, err := NewMidtransGateway(config.MidtransConfig{
		ServerKey:     "sandbox-server-key",
		Environment:   "sandbox",
		ExpiryMinutes: 15,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	validItems := []ItemDetail{{ID: "ticket", Name: "Cinema Ticket", Price: 10000, Quantity: 1}}
	if _, err := gateway.CreateTransaction(CreateTransactionInput{Items: validItems}); err == nil {
		t.Fatal("expected missing order ID error")
	}
	if _, err := gateway.CreateTransaction(CreateTransactionInput{OrderID: "order-1"}); err == nil {
		t.Fatal("expected missing items error")
	}
}

func TestCreateTransactionInputGrossAmountIncludesDiscount(t *testing.T) {
	input := CreateTransactionInput{Items: []ItemDetail{
		{ID: "ticket", Name: "Cinema Ticket", Price: 50000, Quantity: 3},
		{ID: "discount", Name: "Promo Discount", Price: -30000, Quantity: 1},
	}}

	grossAmount, err := input.GrossAmount()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if grossAmount != 120000 {
		t.Fatalf("expected gross amount 120000, got %d", grossAmount)
	}
}
