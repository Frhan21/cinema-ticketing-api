package enums

// Role mendefinisikan peran user dalam sistem.
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

// TicketStatus mendefinisikan status tiket.
type TicketStatus string

const (
	TicketStatusPending   TicketStatus = "pending"
	TicketStatusPaid      TicketStatus = "paid"
	TicketStatusCancelled TicketStatus = "cancelled"
)

// PaymentStatus mendefinisikan status transaksi/pembayaran.
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusPaid      PaymentStatus = "paid"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PaymentMethod mendefinisikan metode pembayaran yang tersedia.
type PaymentMethod string

const (
	PaymentMethodCash         PaymentMethod = "cash"
	PaymentMethodQRIS         PaymentMethod = "qris"
	PaymentMethodTransfer     PaymentMethod = "transfer"
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
)
