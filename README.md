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

2. Isi `.env` dengan kredensial PostgreSQL. Konfigurasi dibaca dari `.env`, lalu fallback ke nilai default di `internal/config/config.go`.

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

   Startup akan menjalankan migrasi dari `database/migrations`, memastikan akun admin, dan memulai scheduler. API tersedia di `http://localhost:8080`.

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
cmd/api/             # entrypoint aplikasi
cmd/setup/           # composition root, infrastruktur, dan scheduler
database/migrations/ # migrasi SQL berpasangan up/down
internal/
  controller/        # HTTP handlers
  service/           # aturan bisnis
  repository/        # query GORM/database
  model/             # model persistensi
  request/           # DTO input
  response/          # DTO dan response envelope
  routes/            # registrasi endpoint
pkg/                 # database, JWT, mailer, error, dan utilitas bersama
```

Alur request: `routes -> controller -> service -> repository -> database`.

## Development Commands

```bash
go fmt ./...                                  # format Go files
go test ./...                                 # compile/test seluruh package
go test ./internal/service                    # focused package check
go build -o cinema-api ./cmd/api              # build binary
~/go/bin/swag init -g cmd/api/main.go --output docs  # update Swagger setelah mengubah annotation API
```

Repository ini belum memiliki test file, Makefile, task runner, CI workflow, atau konfigurasi linter.

## Catatan Kontribusi

- Tambahkan perubahan schema sebagai pasangan migration `.up.sql` dan `.down.sql` di `database/migrations`; aplikasi tidak memanggil `AutoMigration` saat startup.
- ID menggunakan UUID.
- Untuk kontrak endpoint dan proteksi role terkini, gunakan `internal/routes/route.go` sebagai sumber kebenaran.
- Regenerasi `docs/` saat mengubah Swagger annotation atau kontrak API yang terdokumentasi.
