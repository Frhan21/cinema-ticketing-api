# Cinema Ticketing API

## Runtime and local setup

- This is a Go **1.25.0** Gin/GORM PostgreSQL API. The executable entrypoint is `cmd/api/main.go`; application wiring is in `cmd/setup/`.
- Copy `.env.example` to `.env` and supply PostgreSQL credentials. `.env` is loaded automatically, but absent values fall back to the defaults in `internal/config/config.go`.
- Database connections and migrations require PostgreSQL with `sslmode=require`. Startup automatically applies migrations from `database/migrations` using the working-directory-relative path `./database/migrations`; run the API from the repository root.
- Startup also creates or updates the configured admin account. Set **both** `ADMIN_EMAIL` and `ADMIN_PASSWORD`, or set neither to skip seeding. Do not use a real shared-admin password in local `.env` files.
- SMTP is optional in configuration, but background jobs start with every API process: pending transactions are auto-cancelled and film reminders are attempted every minute.

## Commands

```bash
go run cmd/api/main.go                       # starts API, migrations, seeding, and schedulers
go build -o cinema-api ./cmd/api             # build server binary
go test ./...                                # compile/test all packages (no *_test.go files currently)
go test ./internal/service                   # focused package check
go fmt ./...                                 # format Go files
go mod tidy                                  # only after dependency/import changes
~/go/bin/swag init -g cmd/api/main.go --output docs  # regenerate committed Swagger docs after API annotation changes
```

There is no repository Makefile, task runner, CI workflow, or linter configuration.

## Architecture and boundaries

- Request flow is `routes -> controller -> service -> repository -> database`. `cmd/setup/modules.init.go` is the composition root; add new feature dependencies there and include the controller in `routes.RouteControllers`.
- Controllers bind Gin input and return HTTP responses; services own business rules; repositories own GORM/database queries; `model` is persistence data; `request` and `response` are transport DTOs.
- Use the existing `pkg/apperror` constructors for expected service failures and pass errors to Gin with `c.Error(...)`; `internal/middleware/error_handler.go` translates `AppError` values to HTTP responses.
- Preserve the established response envelope in `internal/response`: `message`, optional `data`/`meta`, and optional `error` (not the stale `status` shape described in older prose).
- Keep feature file names in the current dotted convention (for example `ticket.service.go`, `ticket.repository.go`), even though older instructions may mention snake_case.

## Data and API constraints

- Schema changes go in paired `.up.sql` and `.down.sql` files under `database/migrations`; do not rely on `AutoMigration` (it is defined but not called at startup).
- IDs are UUIDs. Monetary fields are currently `float64` in models/services despite `DECIMAL` database columns; avoid expanding that pattern for new exact-money behavior without deliberately addressing the representation.
- Routes are registered under `/api/v1`; `/health` and `/swagger/*any` are top-level. Check `internal/routes/route.go` rather than README tables for current authentication and role protection.
- Swagger docs in `docs/` are generated and imported blank by `cmd/api/main.go`; regenerate them when changing Swagger annotations or documented contracts.
