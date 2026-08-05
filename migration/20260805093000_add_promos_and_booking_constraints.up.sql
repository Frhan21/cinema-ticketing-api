-- Migration: add_promos_and_booking_constraints
CREATE TABLE promos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    discount DECIMAL(5,2) NOT NULL CHECK (discount > 0 AND discount <= 100),
    max_usage INTEGER NOT NULL DEFAULT 0 CHECK (max_usage >= 0),
    used_count INTEGER NOT NULL DEFAULT 0 CHECK (used_count >= 0),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CHECK (end_date >= start_date)
);

CREATE UNIQUE INDEX uq_active_ticket_schedule_seat
    ON tickets(schedule_id, seat_id)
    WHERE status <> 'cancelled';
