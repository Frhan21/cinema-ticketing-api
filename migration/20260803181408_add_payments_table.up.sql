-- Migration: add_payments_table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID NOT NULL UNIQUE,
    order_id VARCHAR(50) NOT NULL UNIQUE,
    midtrans_transaction_id VARCHAR(100) UNIQUE,
    snap_token TEXT NOT NULL,
    redirect_url TEXT NOT NULL,
    payment_type VARCHAR(50),
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    gateway_status VARCHAR(30) NOT NULL DEFAULT 'pending',
    gross_amount BIGINT NOT NULL CHECK (gross_amount > 0),
    paid_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
);

CREATE TABLE payment_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payment_id UUID NOT NULL,
    payment_status VARCHAR(20) NOT NULL,
    gateway_status VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (payment_id) REFERENCES payments(id) ON DELETE CASCADE
);

CREATE INDEX idx_payment_histories_payment_id ON payment_histories(payment_id);
