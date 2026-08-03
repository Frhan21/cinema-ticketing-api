# 🎬 Cinema Ticketing API

REST API untuk pengelolaan bioskop, jadwal tayang, booking kursi, pembayaran tiket, promo, dan laporan penjualan.

![Go](https://img.shields.io/badge/Go-1.25.0-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-1.12-009688)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-required-4169E1?logo=postgresql&logoColor=white)

## Highlights

- JWT authentication dengan role `user` dan `admin`
- CRUD studio, film, kursi, jadwal, dan promo
- Booking beberapa kursi dalam satu transaksi
- Pembayaran dan pembatalan transaksi
- Auto-cancel transaksi pending setelah 15 menit
- Pengingat film 30 menit sebelum jadwal tayang
- Dokumentasi interaktif Swagger

## Tech Stack

| Area | Teknologi |
| --- | --- |
| API | Go 1.25, Gin |
| Database | PostgreSQL, GORM |
| Migrasi | golang-migrate |
| Auth | JWT |
| Dokumentasi | swaggo / Swagger |
| Notifikasi | SMTP via gomail |

## Prasyarat

- Go **1.25.0**
- PostgreSQL yang dapat diakses dengan `sslmode=require`

## Menjalankan Secara Lokal

1. Install dependency dan buat konfigurasi lokal.

   ```bash
   go mod download
   cp .env.example .env
   ```

2. Isi `.env` dengan kredensial PostgreSQL. Konfigurasi dibaca dari `.env`, lalu fallback ke nilai default di `config/config.go`.

   ```env
   APP_PORT=:8080
   DB_HOST=127.0.0.1
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=cinema_ticketing
   JWT_SECRET=replace-with-a-secure-secret
   JWT_EXPIRATION=24h
   ```

3. Jalankan API dari root repository.

   ```bash
   go run cmd/api/main.go
   ```

   Startup akan menjalankan migrasi dari `migration/`, memastikan akun admin, dan memulai scheduler. API tersedia di `http://localhost:8080`.

> **Catatan:** koneksi dan migrasi database menggunakan `sslmode=require`; sesuaikan konfigurasi PostgreSQL lokal Anda. Jangan commit `.env` atau memakai password admin bersama/produksi di mesin lokal.

### Admin Seed

Saat startup, admin dibuat atau diperbarui jika **kedua** variabel berikut diisi:

```env
ADMIN_NAME=Administrator
ADMIN_EMAIL=admin@example.com
ADMIN_PASSWORD=change-this-password
```

Kosongkan `ADMIN_EMAIL` dan `ADMIN_PASSWORD` untuk melewati proses seeding. Jika hanya salah satunya terisi, aplikasi akan berhenti saat startup.

### SMTP (Opsional)

SMTP digunakan untuk email pembayaran, pembatalan, dan pengingat film.

```env
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=your_username
SMTP_PASSWORD=your_password
SMTP_SENDER=no-reply@example.com
```

### Migrasi Database

Migrasi otomatis dijalankan saat startup API, tetapi ada juga CLI migrasi di `main.go` (root repository) untuk membuat dan mengelola migrasi manual. Jalankan dari root repository:

```bash
go run cmd/api/main.go create <nama>     # buat file migrasi baru (up/down) di migration/
go run cmd/api/main.go up                # jalankan semua migrasi yang belum diterapkan
go run cmd/api/main.go down [n]          # rollback n migrasi terakhir (default 1)
go run cmd/api/main.go version           # tampilkan versi migrasi aktif
go run cmd/api/main.go force <versi>     # paksa set versi migrasi (pemulihan state dirty)
```

Koneksi database dibaca dari `.env` (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`).

## Dokumentasi API

| Resource | URL |
| --- | --- |
| Health check | `GET http://localhost:8080/health` |
| Swagger UI | `http://localhost:8080/swagger/index.html` |
| API base URL | `http://localhost:8080/api/v1` |

Endpoint yang memerlukan autentikasi menerima header berikut:

```http
Authorization: Bearer <token>
```

Gunakan `POST /api/v1/auth/register` atau `POST /api/v1/auth/login` untuk memperoleh token.

### Ringkasan Endpoint

| Resource | Endpoint | Akses |
| --- | --- | --- |
| Auth | `POST /auth/register`, `POST /auth/login` | Public |
| User | `GET`, `PUT /user/profile` | Authenticated |
| Studio | `GET /studio/`, `GET /studio/:id`; write endpoints | Authenticated; admin untuk write |
| Movie | `GET /movie/`, `GET /movie/:id`; write endpoints | Public read; admin untuk write |
| Seat | `/seat/`, `/seat/:id`, `/seat/studio/:studio_id` | Authenticated; admin untuk write |
| Schedule | `/schedule/`, `/schedule/:id` | Authenticated; admin untuk write |
| Ticket | `POST /ticket/`, `GET /ticket/history`, `GET /ticket/available-seats/:schedule_id` | Authenticated |
| Transaction | `GET /transaction/:transaction_id`, `POST /pay`, `POST /cancel` | Authenticated |
| Transaction list | `GET /transaction/` | Admin |
| Promo | `GET /promo/validate/:code`; CRUD | Authenticated; admin untuk CRUD |
| Report | `GET /report/daily`, `GET /report/monthly` | Admin |

Detail request, query parameter, dan response tersedia di Swagger UI.

## Contoh Alur Booking

1. Login dan simpan token JWT.
2. Lihat jadwal dan kursi yang tersedia.
3. Booking satu atau lebih kursi.
4. Bayar transaksi yang masih `pending` dalam 15 menit.

```bash
# Booking
curl -X POST http://localhost:8080/api/v1/ticket/ \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "schedule_id": "<schedule-uuid>",
    "seat_id": ["<seat-uuid-1>", "<seat-uuid-2>"],
    "promo_code": "OPTIONAL"
  }'

# Pembayaran
curl -X POST http://localhost:8080/api/v1/transaction/<transaction-uuid>/pay \
  -H 'Authorization: Bearer <token>' \
  -H 'Content-Type: application/json' \
  -d '{"payment_method":"bank_transfer"}'
```

Metode pembayaran yang diterima: `credit_card`, `e_wallet`, dan `bank_transfer`.

## Format Response

Respons sukses dan error menggunakan envelope yang konsisten:

```json
{
  "message": "Ticket booked successfully",
  "data": {}
}
```

```json
{
  "message": "An error occurred",
  "error": "deskripsi error"
}
```

Respons paginasi juga dapat menyertakan `meta`.

## Struktur Proyek

```text
cmd/api/main.go         # entrypoint aplikasi + CLI migrasi (create/up/down/version/force)
cmd/setup/           # composition root, infrastruktur, dan scheduler
migration/           # migrasi SQL berpasangan up/down
entities/            # semua model persistence + enums status (package entities)
app/
  auth/              # autentikasi (register, login)
  user/              # profil user
  movie/             # manajemen film
  studio/            # manajemen studio
  seat/              # manajemen kursi
  schedule/          # jadwal tayang
  ticket/            # booking tiket, transaksi, pembayaran, dan item transaksi
  promo/             # kode promo
  report/            # laporan penjualan admin
  seed/              # seeding admin saat startup
interface/http/
  handler/           # Gin controllers
  middleware/        # auth, error handling
  routes/            # registrasi endpoint
  httpx/             # helper HTTP (misal GetUserIDFromContext)
request/             # DTO input bersama
response/            # response envelope dan DTO bersama
config/              # konfigurasi dari .env
job/                 # background job (auto-cancel, reminder film)
pkg/                 # database, JWT, mailer, error, dan utilitas bersama
```

Setiap feature package mengikuti alur `service -> repository -> database`; handler di `interface/http/handler` menerima request HTTP dan memanggil service.

## Development Commands

```bash
go fmt ./...                                  # format Go files
go test ./...                                 # compile/test seluruh package
go test ./app/ticket                        # focused package check
go build -o cinema-api ./cmd/api              # build binary
~/go/bin/swag init -g cmd/api/main.go --output docs  # update Swagger setelah mengubah annotation API
```

Repository ini belum memiliki test file, Makefile, task runner, CI workflow, atau konfigurasi linter.

## Kontribusi

### Cara Menambah Fitur Baru

Alur dependency: `Routes → Handler → Service → Repository → Database`. Contoh menambah fitur `Review`:

1. **Buat migration** — `go run cmd/api/main.go create create_reviews_table`, lalu isi SQL di `migration/<timestamp>_create_reviews_table.up.sql` dan `.down.sql`.
2. **Buat entity** — `entities/review.entity.go` (`package entities`, tipe `Review`).
3. **Buat feature package** — `app/review/` berisi `review.request.go` (DTO), `review.repository.go`, dan `review.service.go` (interface + implementasi).
4. **Buat handler** — `interface/http/handler/review.handler.go` (`package handler`) dengan annotation Swagger, memanggil service.
5. **Daftarkan route** — di `interface/http/routes/route.go`, pakai `ctrl.Review`.
6. **Wire** — tambahkan di `cmd/setup/modules.init.go`: `reviewRepo`, `reviewService`, lalu `handler.NewReviewController(reviewService)`.
7. **Regenerate Swagger** — `~/go/bin/swag init -g cmd/api/main.go --output docs`.

Konvensi penamaan: `NewXxxHandler/NewXxxService/NewXxxRepository` sebagai constructor; DTO `CreateXxxRequest`/`UpdateXxxRequest`/`XxxResponse`; service didefinisikan sebagai interface. Service tidak boleh menerima `*gin.Context`, handler tidak boleh query database langsung.

### Proses Pull Request

1. Fork repository dan buat branch: `git checkout -b feat/nama-fitur`.
2. Commit dengan format jelas:
   ```
   feat: tambah fitur review film
   fix: perbaiki validasi email duplikat
   refactor: pisahkan logic diskon ke promo service
   docs: update README endpoint terbaru
   ```
3. Pastikan build lolos: `go fmt ./...`, `go mod tidy`, `go build ./...`.
4. Push dan buat Pull Request ke branch `main`.

### Catatan Kontribusi

- Tambahkan perubahan schema sebagai pasangan migration `.up.sql` dan `.down.sql` di `migration/`; aplikasi tidak memanggil `AutoMigration` saat startup. Gunakan `go run cmd/api/main.go create <nama>` untuk membuat pasangan file migrasi baru.
- ID menggunakan UUID.
- Untuk kontrak endpoint dan proteksi role terkini, gunakan `interface/http/routes/route.go` sebagai sumber kebenaran.
- Regenerasi `docs/` saat mengubah Swagger annotation atau kontrak API yang terdokumentasi.
