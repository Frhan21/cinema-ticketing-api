package request

type PayTransactionRequest struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=credit_card e_wallet bank_transfer"`
}
