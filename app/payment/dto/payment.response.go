package dto

type CreatePaymentResponse struct {
	OrderID       string `json:"order_id"`
	Token         string `json:"token"`
	RedirectURL   string `json:"redirect_url"`
	PaymentStatus string `json:"payment_status"`
}
