# 🎬 Cinema Ticketing API

> REST API backend untuk sistem manajemen tiket bioskop modern. Dibangun dengan **Go**, **Gin**, **GORM**, dan **PostgreSQL**.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Gin](https://img.shields.io/badge/Gin-1.12-informational?style=flat)](https://github.com/gin-gonic/gin)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-blue?style=flat&logo=postgresql)](https://www.postgresql.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

---

## 📋 Daftar Isi

- [Fitur Utama](#-fitur-utama)
- [Arsitektur](#-arsitektur)
- [Tech Stack](#-tech-stack)
- [Prasyarat](#-prasyarat)
- [Instalasi](#-instalasi)
- [Konfigurasi](#-konfigurasi-env)
- [Menjalankan Aplikasi](#-menjalankan-aplikasi)
- [API Endpoints](#-api-endpoints)
- [Struktur Proyek](#-struktur-proyek)
- [Dokumentasi Swagger](#-dokumentasi-swagger)

---

## ✨ Fitur Utama

| Fitur | Deskripsi |
|---|---|
| 🔐 **Autentikasi JWT** | Register, login, dan proteksi endpoint berbasis role (admin/user) |
| 🏢 **Manajemen Studio** | CRUD studio bioskop dengan kapasitas dan fasilitas |
| 🎥 **Manajemen Film** | CRUD data film dengan genre, durasi, dan poster |
| 📅 **Manajemen Jadwal** | Penjadwalan film per studio dengan harga fleksibel |
| 💺 **Manajemen Kursi** | Pengelolaan kursi per studio |
| 🎟️ **Booking Tiket** | Pemesanan multi-kursi dengan validasi double-booking |
| 💳 **Transaksi** | Pembayaran dan pembatalan transaksi dengan notifikasi email |
| 🏷️ **Promo & Diskon** | Kode promo dengan persentase diskon dan masa berlaku |
| 📊 **Laporan Penjualan** | Laporan harian dan bulanan per film dan studio |
| ⏰ **Auto-Cancel** | Tiket yang tidak dibayar dalam 15 menit otomatis dibatalkan |
| 📧 **Notifikasi Email** | Email konfirmasi pembayaran, pembatalan, dan pengingat film H-30 menit |

---

## 🏗️ Arsitektur

Proyek ini menggunakan **Clean Architecture** sederhana dengan alur dependency yang jelas:

```
Client Request
    ↓
Routes
    ↓
Controller   (menerima request, mengirim response)
    ↓
Service      (business logic)
    ↓
Repository   (akses database)
    ↓
Database (PostgreSQL)
```

---

## 🛠️ Tech Stack

| Komponen | Teknologi |
|---|---|
| Language | Go 1.21+ |
| Web Framework | Gin v1.12 |
| ORM | GORM v2 |
| Database | PostgreSQL 16+ |
| Auth | JWT (golang-jwt/jwt v5) |
| Email | gomail.v2 + SMTP |
| API Docs | Swagger (swaggo/swag) |
| Config | godotenv |
| Migration | golang-migrate |

---

## 📦 Prasyarat

Pastikan sudah terinstall:

- [Go 1.21+](https://golang.org/dl/)
- [PostgreSQL 16+](https://www.postgresql.org/download/)
- [golang-migrate](https://github.com/golang-migrate/migrate) (untuk menjalankan migrasi database)

```bash
# Install golang-migrate
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

## 🚀 Instalasi

```bash
# 1. Clone repository
git clone https://github.com/username/cinema-ticketing-api.git
cd cinema-ticketing-api

# 2. Install dependencies
go mod download

# 3. Salin file konfigurasi
cp .env.example .env

# 4. Edit file .env sesuai konfigurasi lokal Anda
nano .env

# 5. Jalankan migrasi database
migrate -path database/migrations \
  -database "postgres://USER:PASSWORD@HOST:PORT/cinema_ticketing?sslmode=disable" up

# 6. Jalankan aplikasi
go run cmd/api/main.go
```

---

## ⚙️ Konfigurasi `.env`

Salin `.env.example` dan sesuaikan nilainya:

```env
# Aplikasi
APP_PORT=:8080
APP_ENV=development

# Database PostgreSQL
DB_HOST=127.0.0.1
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=cinema_ticketing

# JWT
JWT_SECRET=your_super_secret_key_min_32_chars
JWT_EXPIRATION=24h

# SMTP (untuk notifikasi email)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password
SMTP_SENDER=your_email@gmail.com

# Akun Admin Default (di-seed otomatis saat pertama kali jalan)
ADMIN_NAME=Administrator
ADMIN_EMAIL=admin@cinemticket.com
ADMIN_PASSWORD=Admin@12345
```

> **Catatan Gmail**: Gunakan [App Password](https://myaccount.google.com/apppasswords), bukan password akun Google biasa.

---

## ▶️ Menjalankan Aplikasi

```bash
# Development
go run cmd/api/main.go

# Build & run binary
go build -o cinema-api ./cmd/api
./cinema-api
```

Server berjalan di `http://localhost:8080`
Swagger UI tersedia di `http://localhost:8080/swagger/index.html`

---

## 📡 API Endpoints

### 🔑 Authentication — `/api/v1/auth`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `POST` | `/auth/register` | Public | Daftar akun baru |
| `POST` | `/auth/login` | Public | Login dan dapatkan JWT token |

### 👤 User — `/api/v1/user`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `GET` | `/user/profile` | User | Lihat profil sendiri |
| `PUT` | `/user/profile` | User | Update profil |

### 🏢 Studio — `/api/v1/studio`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `GET` | `/studio/` | User | List semua studio |
| `GET` | `/studio/:id` | User | Detail studio |
| `POST` | `/studio/` | Admin | Buat studio baru |
| `PUT` | `/studio/:id` | Admin | Update studio |
| `DELETE` | `/studio/:id` | Admin | Hapus studio |

### 🎥 Film — `/api/v1/movie`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `GET` | `/movie/` | Public | List semua film (pagination) |
| `GET` | `/movie/:id` | Public | Detail film |
| `POST` | `/movie/` | Admin | Tambah film |
| `PUT` | `/movie/:id` | Admin | Update film |
| `DELETE` | `/movie/:id` | Admin | Hapus film |

### 💺 Kursi — `/api/v1/seat`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `GET` | `/seat/` | User | List semua kursi |
| `GET` | `/seat/:id` | User | Detail kursi |
| `GET` | `/seat/studio/:studio_id` | User | Kursi per studio |
| `POST` | `/seat/` | Admin | Tambah kursi |
| `PUT` | `/seat/:id` | Admin | Update kursi |
| `DELETE` | `/seat/:id` | Admin | Hapus kursi |

### 📅 Jadwal — `/api/v1/schedule`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `GET` | `/schedule/` | User | List jadwal |
| `GET` | `/schedule/:id` | User | Detail jadwal |
| `POST` | `/schedule/` | Admin | Buat jadwal |
| `PUT` | `/schedule/:id` | Admin | Update jadwal |
| `DELETE` | `/schedule/:id` | Admin | Hapus jadwal |

### 🎟️ Tiket — `/api/v1/ticket`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `POST` | `/ticket/` | User | Booking tiket (bisa pakai promo) |
| `GET` | `/ticket/history` | User | Riwayat transaksi saya |
| `GET` | `/ticket/available-seats/:schedule_id` | User | Cek kursi tersedia per jadwal |

**Contoh request booking:**
```json
{
  "schedule_id": "uuid-jadwal",
  "seat_id": ["uuid-kursi-1", "uuid-kursi-2"],
  "promo_code": "HEMAT20"
}
```

### 💳 Transaksi — `/api/v1/transaction`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `GET` | `/transaction/` | Admin | Semua transaksi |
| `GET` | `/transaction/:transaction_id` | User | Detail transaksi |
| `POST` | `/transaction/:transaction_id/pay` | User | Bayar transaksi |
| `POST` | `/transaction/:transaction_id/cancel` | User | Batalkan transaksi |

**Metode pembayaran yang tersedia:** `credit_card`, `e_wallet`, `bank_transfer`

### 🏷️ Promo — `/api/v1/promo`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `POST` | `/promo/` | Admin | Buat promo baru |
| `GET` | `/promo/` | Admin | List semua promo |
| `GET` | `/promo/:id` | Admin | Detail promo |
| `PUT` | `/promo/:id` | Admin | Update promo |
| `DELETE` | `/promo/:id` | Admin | Hapus promo |
| `GET` | `/promo/validate/:code?total_price=X` | User | Cek validitas kode promo |

### 📊 Laporan — `/api/v1/report`

| Method | Endpoint | Akses | Deskripsi |
|---|---|---|---|
| `GET` | `/report/daily?date=YYYY-MM-DD` | Admin | Laporan harian |
| `GET` | `/report/monthly?year=YYYY&month=M` | Admin | Laporan bulanan |

---

## 📁 Struktur Proyek

```
cinema-ticketing-api/
├── cmd/
│   └── api/
│       └── main.go              # Entry point aplikasi
│
├── database/
│   └── migrations/              # File SQL migration
│
├── docs/                        # Auto-generated Swagger docs
│
├── internal/
│   ├── config/                  # Konfigurasi aplikasi
│   ├── controller/              # HTTP handlers
│   ├── enums/                   # Konstanta status (Role, TicketStatus, dll)
│   ├── middleware/              # Auth, CORS, Role middleware
│   ├── model/                   # Struct model database (GORM)
│   ├── repository/              # Logic akses database
│   ├── request/                 # DTO request dari client
│   ├── response/                # DTO response ke client
│   ├── routes/                  # Registrasi endpoint
│   └── service/                 # Business logic
│
├── pkg/
│   ├── database/                # Koneksi DB dan seeder
│   ├── jwt/                     # Helper JWT
│   ├── mailer/                  # Helper SMTP & template email
│   ├── pagination/              # Helper pagination
│   ├── password/                # Helper hash password
│   ├── scheduler/               # Background job (auto-cancel, reminder)
│   └── validator/               # Custom validator
│
├── .env.example
├── go.mod
├── go.sum
├── README.md
└── CONTRIBUTING.md
```

---

## 📖 Dokumentasi Swagger

Setelah aplikasi berjalan, buka browser dan akses:

```
http://localhost:8080/swagger/index.html
```

Untuk meregenerasi dokumentasi Swagger setelah ada perubahan:

```bash
~/go/bin/swag init -g cmd/api/main.go --output docs
```

---

## 🔒 Format Response API

Semua response menggunakan format yang konsisten:

**Sukses:**
```json
{
  "status": true,
  "message": "success message",
  "data": {}
}
```

**Error:**
```json
{
  "status": false,
  "message": "error message",
  "data": null
}
```

---

## 📄 Lisensi

Proyek ini menggunakan lisensi [MIT](LICENSE).
