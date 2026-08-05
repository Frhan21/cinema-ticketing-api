# Payment and Booking Specification

Status: Implemented with known limitations  
Last updated: 2026-08-05  
Product: Cinema Ticketing API  
Payment provider: Midtrans Snap Redirect

## 1. Purpose

This document is the source of truth for the ticket booking and payment flow. Changes to payment routes, state transitions, amount calculation, persistence, or Midtrans integration must update this specification in the same change.

The selected product flow is:

```text
select schedule and seats
-> create tickets and transaction
-> create Midtrans payment
-> redirect customer to Midtrans
-> process Midtrans notification
-> confirm or cancel tickets
```

Seat selection intentionally happens before payment. This prevents a customer from paying without having specific seats reserved.

## 2. Scope

### In scope

- Authenticated booking for one schedule and one or more seats.
- Optional percentage promo during booking.
- One local payment for one transaction.
- Midtrans Snap Redirect checkout.
- Server-generated item details and gross amount.
- Server-to-server payment notification verification.
- Atomic payment, transaction, ticket, and payment-history status updates.
- Transaction cancellation and automatic expiry of pending transactions.

### Out of scope

- Refund and partial-refund processing.
- Chargeback workflows.
- Multiple payment attempts for one transaction.
- Payment-method-specific promo campaigns.
- Split payments.
- Multi-currency payments; all amounts are IDR.
- Reliable email delivery through a transactional outbox.

## 3. Terms

| Term | Meaning |
| --- | --- |
| Booking | Selection of a schedule and concrete seats before payment. |
| Transaction | The local aggregate containing the user, total price, status, and ticket references. |
| Transaction item | Link between a transaction and a ticket, including the server-side item price. |
| Payment | Local Midtrans payment record associated one-to-one with a transaction. |
| Order ID | Merchant-generated Midtrans order ID in the form `cinema-{transaction_uuid}`. |
| Gateway status | Original Midtrans status such as `pending`, `settlement`, or `expire`. |
| Payment status | Simplified local status: `pending`, `paid`, or `cancelled`. |

## 4. Functional Requirements

### PAY-001 — Booking precedes payment

An authenticated user must select a schedule and one or more seats before initiating payment.

Booking request:

```http
POST /api/v1/ticket/
Authorization: Bearer <token>
Content-Type: application/json
```

```json
{
  "schedule_id": "<schedule-uuid>",
  "seat_id": ["<seat-uuid-1>", "<seat-uuid-2>"],
  "promo_code": "OPTIONAL"
}
```

The server must obtain schedule price and seat information from the database. The client must not submit item prices or a final payable amount.

### PAY-002 — Booking output

A successful booking creates:

- One pending ticket for each selected seat.
- One pending transaction owned by the current user.
- One transaction item for each ticket.

The response returns `transaction_id`, `ticket_ids`, `total_price`, and `payment_status`.

### PAY-003 — Payment initiation

An authenticated owner initiates payment using:

```http
POST /api/v1/transaction/{transaction_id}/pay
Authorization: Bearer <token>
```

The request has no JSON body. Payment method selection belongs to the Midtrans checkout page.

The server must reject payment initiation when:

- The transaction does not exist.
- The transaction belongs to another user.
- The transaction status is not `pending`.
- The transaction contains no items.
- An item or final amount is not a positive whole IDR value.
- Item totals do not reconcile with the transaction total.

### PAY-004 — Idempotent initiation

Only one payment may exist for one transaction. If a payment already exists, the endpoint returns its existing Snap token and redirect URL without creating another Midtrans order.

### PAY-005 — Item details and gross amount

The server builds Midtrans item details from persisted transaction items. Ticket items with the same price are grouped:

```json
{
  "id": "cinema-ticket-1",
  "name": "Cinema Ticket",
  "price": 50000,
  "quantity": 2
}
```

`gross_amount` is calculated as:

```text
sum(item.price * item.quantity)
```

It is never accepted from the payment-initiation request.

### PAY-006 — Simple promo representation

Booking currently stores the discounted final total but does not persist a promo snapshot on the transaction. Payment therefore derives a generic discount from:

```text
discount = transaction item subtotal - transaction total price
```

When the discount is greater than zero, the server appends a negative Midtrans item:

```json
{
  "id": "promo-discount",
  "name": "Promo Discount",
  "price": -20000,
  "quantity": 1
}
```

The final item sum must exactly equal the local transaction total.

### PAY-007 — Snap response

A successful initiation returns:

```json
{
  "message": "Payment created successfully",
  "data": {
    "order_id": "cinema-<transaction-uuid>",
    "token": "<snap-token>",
    "redirect_url": "https://app.sandbox.midtrans.com/...",
    "payment_status": "pending"
  }
}
```

The frontend redirects the customer to `redirect_url`.

### PAY-008 — Notification endpoint

Midtrans sends status changes to:

```http
POST /api/v1/payment/notification
Content-Type: application/json
```

This endpoint is public because it is called by Midtrans and must not require JWT authentication.

### PAY-009 — Notification authenticity

The server must not trust notification status or amount directly. It uses the notification `order_id` to find a local payment, then calls Midtrans GET Status through the configured server key.

The verified response must match:

- Local payment order ID.
- Local payment gross amount.

An invalid or mismatched payment must not update local state.

### PAY-010 — Status mapping

| Midtrans transaction status | Additional rule | Local payment/transaction status | Ticket status |
| --- | --- | --- | --- |
| `pending` | — | `pending` | `pending` |
| `capture` | fraud status `accept` | `paid` | `paid` |
| `capture` | fraud status other than `accept` | `pending` | `pending` |
| `settlement` | — | `paid` | `paid` |
| `deny` | — | `cancelled` | `cancelled` |
| `cancel` | — | `cancelled` | `cancelled` |
| `expire` | — | `cancelled` | `cancelled` |
| `failure` | — | `cancelled` | `cancelled` |

Unsupported statuses such as `refund` are currently acknowledged without a state mutation and require a future specification update.

### PAY-011 — Atomic status update

Notification processing must update the following in one database transaction:

- Payment status and gateway metadata.
- Payment history when status changes.
- Transaction payment status and payment method.
- All tickets linked through transaction items.

A row lock on the payment prevents duplicate notifications from inserting duplicate history records. Email or other network side effects must not run inside the database transaction.

### PAY-012 — Frontend confirmation

Midtrans redirect parameters are not proof of payment. After returning from Midtrans, the frontend must query the local transaction and display:

- Pending screen for `pending`.
- Success screen for `paid`.
- Failed or cancelled screen for `cancelled`.

## 5. Data Invariants

- IDs are UUIDs, except Midtrans order and provider transaction IDs.
- `payments.transaction_id` is unique.
- `payments.order_id` is unique.
- Payment gross amount is a positive whole IDR `BIGINT`.
- Snap token is not serialized from the persistence entity; it is returned through a response DTO.
- Midtrans item IDs in one request are unique.
- Every selected seat must belong to the schedule studio.
- A non-cancelled ticket occupies one seat for one schedule.
- `Seat.IsAvailable` is not a schedule-specific occupancy flag and must not be toggled when a ticket is paid.

## 6. Configuration

Required for payment-enabled application startup:

```env
MIDTRANS_SERVER_KEY="your-midtrans-server-key"
MIDTRANS_ENV="sandbox"
PAYMENT_EXPIRY_MINUTES=15
```

Allowed environments are `sandbox` and `production`. `PAYMENT_EXPIRY_MINUTES` must be positive and is shared by Snap expiry and the local auto-cancel scheduler. The Snap expiry starts at the persisted transaction creation time. The server key is backend-only and must never be exposed to the frontend or committed to source control.

Configure the Midtrans Payment Notification URL as:

```text
https://<public-api-host>/api/v1/payment/notification
```

Local development requires a public HTTPS tunnel because Midtrans cannot call localhost.

## 7. Failure and Retry Behavior

- Repeating payment initiation returns the existing local payment.
- Repeating the same notification is a no-op after the first committed status transition.
- Database or Midtrans verification failures return an error so notification delivery can be retried.
- Email failure must not roll back a confirmed payment.
- A transaction pending beyond `PAYMENT_EXPIRY_MINUTES` is auto-cancelled by the local scheduler and uses the same start time and duration as Snap.

## 8. Security Requirements

- The authenticated user ID comes from JWT context, never from the request body.
- Payment initiation verifies transaction ownership.
- Prices, quantity, promo adjustment, and gross amount are derived from persisted server data.
- Notification status is verified with Midtrans GET Status.
- Server keys, Snap tokens, and sensitive customer data must not be logged.
- SQL operations use GORM parameter binding; raw string interpolation is prohibited.

## 9. Acceptance Criteria

- Creating a payment for two IDR 50,000 tickets produces quantity `2` and gross amount IDR 100,000.
- A transaction discounted from IDR 100,000 to IDR 80,000 includes a negative IDR 20,000 promo item and gross amount IDR 80,000.
- A user cannot initiate another user's payment.
- A repeated initiation does not call Midtrans again.
- A verified settlement marks payment, transaction, and tickets paid atomically.
- A mismatched Midtrans amount does not mutate local data.
- A repeated settlement does not create duplicate payment history.
- The notification route is reachable without JWT authentication.
- Promo usage, transaction, tickets, and transaction items are committed atomically.
- Concurrent booking cannot create two active tickets for the same schedule and seat.
- Transaction detail is visible only to its owner or an admin.
- A paid transition schedules one confirmation email after the database commit.

## 10. Known Limitations

- Promo identity and discount snapshot are not persisted on the transaction; Midtrans displays a generic `Promo Discount` item.
- Booking and promo preview currently duplicate promo-validation rules; future changes must keep both paths consistent.
- Local cancellation does not cancel or expire the corresponding Midtrans transaction.
- Manual cancellation and auto-cancellation still update transaction and ticket rows through separate repository calls rather than one database transaction.
- Payment email is best-effort asynchronous delivery; there is no transactional outbox or retry queue.

## 11. Verification

Run after modifying this flow:

```bash
go test ./...
git diff --check
```

When Swagger alias resolution is repaired, also run:

```bash
~/go/bin/swag init -g cmd/api/main.go --output docs
```
