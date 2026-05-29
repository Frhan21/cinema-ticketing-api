package response

import "github.com/google/uuid"

type TicketResponse struct {
	TransactionId uuid.UUID   `json:"transaction_id"`
	TicketIds     []uuid.UUID `json:"ticket_ids"`
	TotalPrice    float64     `json:"total_price"`
	PaymentStatus string      `json:"payment_status"`
}

type TransactionHistoryResponse struct {
	ID            string  `json:"id"`
	TotalPrice    float64 `json:"total_price"`
	PaymentStatus string  `json:"payment_status"`
	CreatedAt     string  `json:"created_at"`
}
