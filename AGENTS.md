# Cinema Ticketing API

## Runtime and local setup

- This is a Go **1.25.0** Gin/GORM PostgreSQL API. The executable entrypoint is `cmd/api/main.go`; application wiring is in `cmd/setup/`.
- Copy `.env.example` to `.env` and supply PostgreSQL credentials. `.env` is loaded automatically, but absent values fall back to the defaults in `config/config.go`.
- Database connections and migrations require PostgreSQL with `sslmode=require`. Startup automatically applies migrations from `migration/` using the working-directory-relative path `./migration`; run the API from the repository root.
- Startup also creates or updates the configured admin account. Set **both** `ADMIN_EMAIL` and `ADMIN_PASSWORD`, or set neither to skip seeding. Do not use a real shared-admin password in local `.env` files.
- SMTP is optional in configuration, but background jobs start with every API process: pending transactions are auto-cancelled and film reminders are attempted every minute.

## Commands

```bash
go run cmd/api/main.go                       # starts API, migrations, seeding, and schedulers
go build -o cinema-api ./cmd/api             # build server binary
go test ./...                                # compile/test all packages (no *_test.go files currently)
go test ./app/ticket                         # focused package check
go fmt ./...                                 # format Go files
go mod tidy                                  # only after dependency/import changes
~/go/bin/swag init -g cmd/api/main.go --output docs  # regenerate committed Swagger docs after API annotation changes
```

Migration CLI (dispatcher in `cmd/api/main.go`, run from the repository root; reads DB config from `.env`):

```bash
go run cmd/api/main.go create <name>       # create a new migration pair in migration/
go run cmd/api/main.go up                  # apply pending migrations
go run cmd/api/main.go down [n]            # roll back n migrations (default 1)
go run cmd/api/main.go version             # show active migration version
go run cmd/api/main.go force <version>     # force schema version (dirty-state recovery)
```

Running `cmd/api/main.go` with any other argument (or none) starts the API server as usual.

There is no repository Makefile, task runner, CI workflow, or linter configuration.

## Architecture and boundaries

- The codebase uses a **top-level layered layout** (not `internal/` feature packages):
  - `app/<feature>/` owns the feature's repositories and services (for example `app/auth`, `app/user`, `app/movie`, `app/studio`, `app/seat`, `app/schedule`, `app/ticket`, `app/promo`, `app/report`, `app/seed`). Within a feature the flow is `service -> repository -> database`; `cmd/setup/modules.init.go` is the composition root — add new feature dependencies there and wire the handler into `routes.RouteControllers`.
  - `entities/` holds every persistence model **and** the status enums in a single `package entities` (`user.entity.go`, `movie.entity.go`, `enums.go`, ...). All GORM relations reference the shared `entities` package, which removes cross-package relation cycles.
  - `interface/http/` owns the HTTP transport: `handler/` (Gin controllers), `middleware/` (auth, error handling), `routes/` (route registration), and `httpx/` (HTTP helpers such as `GetUserIDFromContext`).
  - `request/` and `response/` hold shared transport DTOs and the response envelope; `config/` holds configuration loading; `job/` owns the `Scheduler`; `migration/` holds paired SQL migrations; `pkg/` holds cross-cutting libraries (`apperror`, `database`, `jwt`, `mailer`, `pagination`, `password`, `utils`).
- Handlers bind Gin input and return HTTP responses; services own business rules; repositories own GORM/database queries; feature DTOs live in `app/<feature>/<feature>.request.go` and `app/<feature>/<feature>.response.go`.
- `app/ticket` also owns the transaction flow (`Transaction`, `TransactionItem`, `TransactionService`). Booking and payment are one feature: `TransactionItem` links `Transaction` and `Ticket`, and both `TicketService.BookTicket` and `TransactionService.PayTransaction` mutate the same aggregate. Do not split these back into two packages without first breaking the cross-package GORM relations.
- `job.Scheduler` owns the background jobs (auto-cancel expired transactions, upcoming-film reminders). Reminders intentionally depend on `schedule.ScheduleRepository` and `ticket.TicketRepository` from the job layer, not from `app/schedule`, to avoid an import cycle (`schedule` is imported by `ticket` models).
- Use the existing `pkg/apperror` constructors for expected service failures and pass errors to Gin with `c.Error(...)`; `interface/http/middleware/error_handler.go` translates `AppError` values to HTTP responses.
- Preserve the established response envelope in `response/`: `message`, optional `data`/`meta`, and optional `error` (not the stale `status` shape described in older prose).
- Keep feature file names in the current dotted convention (for example `ticket.service.go`, `ticket.repository.go`). Sub-feature files inside a package keep their own prefix (for example `transaction.service.go`, `transaction_item.repository.go` in `app/ticket`).
- When a handler imports a feature package whose name collides with a local variable (for example the local `schedule`/`promo`/`seat` variables), import the feature with an alias (`schedulepkg`, `promopkg`, `seatpkg`) and keep the alias **only** for type references — field accesses on local variables must stay unqualified.

## Code conventions

- Dependency flow (must never be reversed): `Routes → Handler → Service → Repository → Database`. Wire new dependencies in `cmd/setup/modules.init.go` and register the handler in `routes.RouteControllers`.
- Handlers, services, and repositories are structs exposed through `NewXxx` constructors (for example `NewUserHandler(userService user.UserService)`). Services are defined by interfaces (`type UserService interface { ... }`) in their feature package.
- Request/response DTO structs are named `CreateXxxRequest`, `UpdateXxxRequest`, `XxxResponse`; persistence structs in `entities/` are named `Xxx`.
- Forbidden: handlers must not touch the database directly; repositories must not call services; services must not accept `*gin.Context`; entities must not contain business logic.
- Use `response.SuccessResponse(message, data)` / `response.ErrorResponse(message)` for HTTP output, `pkg/apperror` constructors for expected service failures, and Gin binding tags for request validation.

## Data and API constraints

- Schema changes go in paired `.up.sql` and `.down.sql` files under `migration/`; do not rely on `AutoMigration` (it is defined but not called at startup).
- IDs are UUIDs. Monetary fields are currently `float64` in models/services despite `DECIMAL` database columns; avoid expanding that pattern for new exact-money behavior without deliberately addressing the representation.
- Routes are registered under `/api/v1`; `/health` and `/swagger/*any` are top-level. Check `interface/http/routes/route.go` rather than README tables for current authentication and role protection.
- Swagger docs in `docs/` are generated and imported blank by `cmd/api/main.go`; regenerate them when changing Swagger annotations or documented contracts.
