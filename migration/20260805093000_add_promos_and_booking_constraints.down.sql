-- Migration: add_promos_and_booking_constraints
DROP INDEX IF EXISTS uq_active_ticket_schedule_seat;
DROP TABLE IF EXISTS promos;
