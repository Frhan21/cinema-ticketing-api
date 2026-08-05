package dto

type MidtransNotificationRequest struct {
	OrderID           string `json:"order_id" binding:"required"`
	TransactionID     string `json:"transaction_id"`
	TransactionStatus string `json:"transaction_status" binding:"required"`
	PaymentType       string `json:"payment_type"`
	GrossAmount       string `json:"gross_amount" binding:"required"`
	StatusCode        string `json:"status_code" binding:"required"`
	SignatureKey      string `json:"signature_key" binding:"required"`
	FraudStatus       string `json:"fraud_status"`
	SettlementTime    string `json:"settlement_time"`
}
