FROM golang:1.25.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/api/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/migration /app/migration
# Copy .env jika digunakan (tidak recommended untuk production)
# COPY --from=builder /app/.env /app/.env 
EXPOSE 8000
CMD ["/app/main"]