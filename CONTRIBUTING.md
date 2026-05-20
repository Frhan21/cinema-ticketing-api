# 🤝 Contributing Guide — Cinema Ticketing API

Terima kasih sudah tertarik untuk berkontribusi! Panduan ini menjelaskan cara setup environment pengembangan, konvensi kode, dan proses pull request.

---

## 📋 Daftar Isi

- [Persiapan Environment](#-persiapan-environment)
- [Struktur Folder](#-struktur-folder)
- [Konvensi Kode](#-konvensi-kode)
- [Alur Arsitektur](#-alur-arsitektur)
- [Cara Menambah Fitur Baru](#-cara-menambah-fitur-baru)
- [Enums & Status](#-enums--status)
- [Validasi & Response](#-validasi--response)
- [Proses Pull Request](#-proses-pull-request)
- [Perintah Berguna](#-perintah-berguna)

---

## 🛠️ Persiapan Environment

```bash
# 1. Fork dan clone repository
git clone https://github.com/username/cinema-ticketing-api.git
cd cinema-ticketing-api

# 2. Install dependencies
go mod download

# 3. Setup konfigurasi
cp .env.example .env
# Edit .env sesuai environment lokal Anda

# 4. Jalankan migrasi
migrate -path database/migrations \
  -database "postgres://USER:PASS@localhost:5432/cinema_ticketing?sslmode=disable" up

# 5. Jalankan aplikasi
go run cmd/api/main.go
```

---

## 📁 Struktur Folder

Proyek menggunakan **Clean Architecture** sederhana. Setiap folder memiliki tanggung jawab yang jelas:

| Folder | Tanggung Jawab |
|---|---|
| `internal/config` | Load konfigurasi dari `.env` |
| `internal/controller` | Menerima request HTTP, memanggil service, mengirim response |
| `internal/service` | Business logic aplikasi |
| `internal/repository` | Semua query database (GORM) |
| `internal/model` | Struct model database |
| `internal/request` | DTO input dari client |
| `internal/response` | DTO output ke client |
| `internal/routes` | Registrasi semua endpoint |
| `internal/middleware` | Auth, role, CORS |
| `internal/enums` | Konstanta status (TicketStatus, PaymentStatus, dll) |
| `pkg/database` | Koneksi DB dan seeder |
| `pkg/mailer` | Kirim email via SMTP |
| `pkg/scheduler` | Background job |
| `pkg/jwt` | Helper JWT |
| `pkg/password` | Helper hash password |

---

## 📐 Konvensi Kode

### Penamaan File

Gunakan **snake_case** untuk nama file:

```
user_controller.go
user_service.go
user_repository.go
user.model.go
auth.request.go
ticket.response.go
```

### Penamaan Package

Package mengikuti nama folder:

```go
package controller
package service
package repository
package model
package request
package response
package enums
```

### Penamaan Struct

```go
type UserController struct {}
type UserService struct {}
type UserRepository struct {}
type CreateUserRequest struct {}
type UserResponse struct {}
```

### Constructor Function

Selalu gunakan constructor function:

```go
func NewUserController(userService service.UserService) UserController
func NewUserService(userRepo repository.UserRepository) UserService
func NewUserRepository(db *gorm.DB) UserRepository
```

### Interface

Setiap layer (kecuali Controller yang sudah struct) menggunakan interface:

```go
type UserService interface {
    GetProfile(userID uuid.UUID) (*response.UserResponse, error)
    UpdateProfile(userID uuid.UUID, req request.UpdateUserRequest) (*response.UserResponse, error)
}
```

---

## 🔄 Alur Arsitektur

**WAJIB** diikuti — jangan melanggar dependency flow ini:

```
Routes → Controller → Service → Repository → Database
```

**Larangan:**
- ❌ Controller tidak boleh query database langsung
- ❌ Repository tidak boleh memanggil service
- ❌ Service tidak boleh menerima `*gin.Context`
- ❌ Model tidak boleh berisi business logic

---

## ➕ Cara Menambah Fitur Baru

Contoh menambah fitur `Review` (ulasan film):

### 1. Buat Migration

```sql
-- database/migrations/000010_create_reviews_table.up.sql
CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id),
    movie_id UUID NOT NULL REFERENCES movies(id),
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

```sql
-- database/migrations/000010_create_reviews_table.down.sql
DROP TABLE IF EXISTS reviews;
```

### 2. Buat Model

```go
// internal/model/review.model.go
package model

type Review struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    UserID    uuid.UUID `gorm:"not null" json:"user_id"`
    MovieID   uuid.UUID `gorm:"not null" json:"movie_id"`
    Rating    int       `gorm:"not null" json:"rating"`
    Comment   string    `json:"comment"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    User  User  `gorm:"foreignKey:UserID" json:"user"`
    Movie Movie `gorm:"foreignKey:MovieID" json:"movie"`
}
```

### 3. Buat Request & Response

```go
// internal/request/review.request.go
type CreateReviewRequest struct {
    MovieID uuid.UUID `json:"movie_id" binding:"required"`
    Rating  int       `json:"rating" binding:"required,min=1,max=5"`
    Comment string    `json:"comment"`
}
```

### 4. Buat Repository

```go
// internal/repository/review.repository.go
type ReviewRepository interface {
    Create(review *model.Review) error
    FindByMovieID(movieID uuid.UUID) ([]model.Review, error)
}
```

### 5. Buat Service

```go
// internal/service/review.service.go
type ReviewService interface {
    Create(userID uuid.UUID, req request.CreateReviewRequest) error
    GetByMovie(movieID uuid.UUID) ([]model.Review, error)
}
```

### 6. Buat Controller (dengan Swagger annotation)

```go
// internal/controller/review.controller.go

// Create godoc
// @Summary      Buat ulasan film
// @Tags         Review
// @Accept       json
// @Security     BearerAuth
// @Param        body body request.CreateReviewRequest true "Data ulasan"
// @Success      201  {object} map[string]interface{}
// @Router       /review [post]
func (r *reviewController) Create(c *gin.Context) {
    // ...
}
```

### 7. Daftarkan ke Routes

```go
// internal/routes/route.go
reviewRoute := api.Group("/review")
reviewRoute.Use(middleware.AuthMiddleware(jwtCfg))
{
    reviewRoute.POST("/", ctrl.Review.Create)
    reviewRoute.GET("/movie/:movie_id", ctrl.Review.GetByMovie)
}
```

### 8. Wire di main.go

```go
reviewRepo := repository.NewReviewRepository(db)
reviewService := service.NewReviewService(reviewRepo)
reviewController := controller.NewReviewController(reviewService)
```

### 9. Regenerate Swagger

```bash
~/go/bin/swag init -g cmd/api/main.go --output docs
```

---

## 📌 Enums & Status

Semua konstanta status **harus** diambil dari `internal/enums/enums.go`:

```go
enums.TicketStatusPending    // "pending"
enums.TicketStatusPaid       // "paid"
enums.TicketStatusCancelled  // "cancelled"

enums.PaymentStatusPending   // "pending"
enums.PaymentStatusPaid      // "paid"
enums.PaymentStatusCancelled // "cancelled"

enums.PaymentMethodCash      // "cash"
enums.PaymentMethodQRIS      // "qris"
enums.PaymentMethodTransfer  // "transfer"
enums.PaymentMethodCreditCard // "credit_card"

enums.RoleAdmin              // "admin"
enums.RoleUser               // "user"
```

---

## ✅ Validasi & Response

### Format Response

Gunakan **selalu** helper dari `internal/response`:

```go
// Sukses
c.JSON(http.StatusOK, response.SuccessResponse("message", data))

// Error
c.JSON(http.StatusBadRequest, response.ErrorResponse("error message"))
```

### Format JSON yang dihasilkan

```json
// Sukses
{ "status": true, "message": "...", "data": {} }

// Error
{ "status": false, "message": "...", "data": null }
```

### Validasi Request

Gunakan binding tag Gin:

```go
type CreateRequest struct {
    Name  string `json:"name" binding:"required,min=3"`
    Email string `json:"email" binding:"required,email"`
    Price float64 `json:"price" binding:"required,min=0"`
}
```

---

## 🔀 Proses Pull Request

1. **Fork** repository ini
2. Buat branch baru: `git checkout -b feat/nama-fitur`
3. Commit dengan format yang jelas:
   ```
   feat: tambah fitur review film
   fix: perbaiki validasi email duplikat
   refactor: pisahkan logic diskon ke promo service
   docs: update README endpoint terbaru
   ```
4. Pastikan build berhasil:
   ```bash
   go fmt ./...
   go mod tidy
   go build ./...
   ```
5. Push dan buat Pull Request ke branch `main`

---

## ⚡ Perintah Berguna

```bash
# Format kode
go fmt ./...

# Tidy dependencies
go mod tidy

# Build check
go build ./...

# Run tests
go test ./...

# Run aplikasi
go run cmd/api/main.go

# Generate/update Swagger docs
~/go/bin/swag init -g cmd/api/main.go --output docs

# Jalankan migrasi naik
migrate -path database/migrations \
  -database "postgres://USER:PASS@HOST:PORT/DB_NAME?sslmode=disable" up

# Rollback migrasi terakhir
migrate -path database/migrations \
  -database "postgres://USER:PASS@HOST:PORT/DB_NAME?sslmode=disable" down 1
```

---

## ❓ Ada Pertanyaan?

Buka [GitHub Issue](https://github.com/username/cinema-ticketing-api/issues) atau hubungi maintainer di email yang tercantum di README.
