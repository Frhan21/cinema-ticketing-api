# Workspace Rules - Go REST API Backend

## 1. Project Context

Project ini adalah REST API backend menggunakan bahasa Go.

Project menggunakan pendekatan Clean Architecture sederhana, tetapi tetap memakai istilah folder yang mudah dipahami.

Gunakan istilah berikut:

- controller
- service
- repository
- model
- request
- response
- routes
- middleware

Jangan gunakan istilah folder berikut:

- usecase
- delivery
- domain

Alasan:

- `service` lebih mudah dipahami daripada `usecase`
- `controller` lebih mudah dipahami daripada `delivery/http`
- `model` lebih mudah dipahami daripada `domain`
- Struktur ini tetap clean, tetapi lebih sederhana untuk project training backend

---

## 2. Target Folder Structure

Struktur folder utama project harus mengikuti bentuk berikut:

```txt
.
├── cmd/
│   └── api/
│       └── main.go
│
├── database/
│   └── migrations/
│
├── internal/
│   ├── config/
│   ├── controller/
│   ├── service/
│   ├── repository/
│   ├── model/
│   ├── request/
│   ├── response/
│   ├── routes/
│   └── middleware/
│
├── pkg/
│   ├── database/
│   ├── jwt/
│   ├── password/
│   └── validator/
│
├── go.mod
├── go.sum
├── .env
└── README.md
```

````

---

## 3. Architecture Flow

Dependency flow wajib mengikuti alur berikut:

```txt
routes -> controller -> service -> repository -> database
```

Response flow mengikuti alur berikut:

```txt
database -> repository -> service -> controller -> client
```

Tidak boleh membuat dependency seperti:

```txt
controller -> database
controller -> repository
repository -> service
repository -> controller
model -> controller
model -> service
request -> repository
response -> repository
```

Aturan utama:

- routes hanya mengenal controller
- controller hanya memanggil service
- service hanya memanggil repository
- repository hanya berurusan dengan database
- model hanya berisi struktur data database
- request hanya berisi payload request
- response hanya berisi payload response dan helper response

---

## 4. Folder Responsibility

### 4.1 cmd/api

Folder:

```txt
cmd/api
```

Berisi entry point aplikasi.

File utama:

```txt
cmd/api/main.go
```

Tanggung jawab:

- load konfigurasi aplikasi
- inisialisasi koneksi database
- inisialisasi repository
- inisialisasi service
- inisialisasi controller
- setup routes
- menjalankan HTTP server

Tidak boleh:

- menaruh business logic
- menaruh query database
- menaruh logic controller
- menaruh logic repository

Contoh isi yang diperbolehkan:

```go
func main() {
    cfg := config.LoadConfig()
    db := database.Connect(cfg)

    router := gin.Default()

    routes.SetupRoutes(router, db)

    router.Run(cfg.AppPort)
}
```

---

### 4.2 database/migrations

Folder:

```txt
database/migrations
```

Berisi file migration database.

Contoh:

```txt
000001_create_users_table.up.sql
000001_create_users_table.down.sql
000002_create_movies_table.up.sql
000002_create_movies_table.down.sql
000003_create_studios_table.up.sql
000003_create_studios_table.down.sql
```

Aturan:

- semua migration SQL disimpan di folder ini
- setiap migration harus punya file `.up.sql` dan `.down.sql`
- migration tidak boleh disimpan di root project
- migration tidak boleh dicampur dengan model Go

---

### 4.3 internal/config

Folder:

```txt
internal/config
```

Berisi konfigurasi aplikasi.

Tanggung jawab:

- membaca environment variable
- membaca file `.env`
- menyimpan konfigurasi database
- menyimpan konfigurasi server
- menyimpan konfigurasi JWT
- menyimpan konfigurasi app environment

Contoh konfigurasi:

```go
type Config struct {
    AppPort     string
    DBHost      string
    DBPort      string
    DBUser      string
    DBPassword  string
    DBName      string
    JWTSecret   string
}
```

Tidak boleh:

- menaruh query database
- menaruh business logic
- menaruh controller
- menaruh response HTTP

---

### 4.4 internal/routes

Folder:

```txt
internal/routes
```

Berisi route registration.

Tanggung jawab:

- mendaftarkan endpoint API
- menghubungkan endpoint dengan controller
- memasang middleware pada route tertentu
- mengelompokkan route berdasarkan fitur

Contoh:

```go
func SetupRoutes(router *gin.Engine, userController *controller.UserController) {
    api := router.Group("/api/v1")

    users := api.Group("/users")
    {
        users.GET("", userController.FindAll)
        users.POST("", userController.Create)
        users.GET("/:id", userController.FindByID)
        users.PUT("/:id", userController.Update)
        users.DELETE("/:id", userController.Delete)
    }
}
```

Tidak boleh:

- menaruh business logic
- menaruh query database
- bind request body
- membuat response JSON langsung kecuali untuk health check sederhana

---

### 4.5 internal/controller

Folder:

```txt
internal/controller
```

Berisi HTTP controller atau handler.

Tanggung jawab:

- menerima HTTP request dari Gin
- mengambil parameter URL
- mengambil query parameter
- bind JSON request body
- validasi request dasar
- memanggil service
- mengirim response JSON ke client

Controller boleh:

- menggunakan `c.ShouldBindJSON()`
- membaca `c.Param()`
- membaca `c.Query()`
- memanggil method dari service
- memakai helper response

Controller tidak boleh:

- melakukan query database langsung
- memanggil GORM langsung
- memanggil repository langsung
- menyimpan business logic kompleks
- melakukan hash password langsung
- generate JWT langsung jika logic tersebut sudah ada di service/helper
- menghitung transaksi kompleks
- mengatur relasi database

Contoh alur controller:

```go
func (ctrl *UserController) Create(c *gin.Context) {
    var req request.CreateUserRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    result, err := ctrl.userService.Create(req)
    if err != nil {
        response.Error(c, http.StatusBadRequest, err.Error())
        return
    }

    response.Success(c, http.StatusCreated, "user created successfully", result)
}
```

---

### 4.6 internal/service

Folder:

```txt
internal/service
```

Berisi business logic aplikasi.

Tanggung jawab:

- menjalankan aturan bisnis aplikasi
- validasi proses bisnis
- mengolah request dari controller
- memanggil repository
- mengubah data model menjadi response
- menjalankan proses seperti register, login, booking tiket, cek kursi, hitung harga, update status transaksi

Service boleh:

- memanggil repository
- memakai helper dari `pkg`
- melakukan hash password
- generate JWT melalui helper
- validasi role
- validasi status transaksi
- menghitung total pembayaran
- mengecek ketersediaan kursi
- menggabungkan data dari beberapa repository

Service tidak boleh:

- menerima `*gin.Context`
- mengirim response JSON langsung
- menjalankan route
- membaca HTTP request langsung
- terlalu banyak berisi query GORM detail
- mengatur migration database

Contoh service:

```go
func (s *UserService) Register(req request.RegisterRequest) (*response.UserResponse, error) {
    existingUser, _ := s.userRepository.FindByEmail(req.Email)
    if existingUser != nil {
        return nil, errors.New("email already registered")
    }

    hashedPassword, err := password.Hash(req.Password)
    if err != nil {
        return nil, err
    }

    user := model.User{
        ID:       uuid.NewString(),
        Name:     req.Name,
        Email:    req.Email,
        Password: hashedPassword,
    }

    savedUser, err := s.userRepository.Create(user)
    if err != nil {
        return nil, err
    }

    result := response.UserResponse{
        ID:    savedUser.ID,
        Name:  savedUser.Name,
        Email: savedUser.Email,
    }

    return &result, nil
}
```

---

### 4.7 internal/repository

Folder:

```txt
internal/repository
```

Berisi logic akses database.

Tanggung jawab:

- create data
- find data
- update data
- delete data
- query database menggunakan GORM atau SQL
- mengambil data berdasarkan filter
- menjalankan transaction jika diperlukan

Repository boleh:

- menerima dependency database seperti `*gorm.DB`
- menggunakan GORM
- menggunakan SQL query
- menerima dan mengembalikan `model`

Repository tidak boleh:

- menerima `*gin.Context`
- membuat response JSON
- membaca request HTTP
- memanggil controller
- memanggil service
- mengatur validasi business logic
- generate JWT
- hash password kecuali sangat diperlukan dan disepakati

Contoh repository:

```go
type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(user model.User) (*model.User, error) {
    if err := r.db.Create(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User

    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
}
```

---

### 4.8 internal/model

Folder:

```txt
internal/model
```

Berisi struct model database.

Tanggung jawab:

- mendefinisikan struktur tabel
- mendefinisikan relasi antar tabel
- menyimpan tag GORM
- menyimpan field database

Model boleh:

- berisi tag `gorm`
- berisi tag `json` jika diperlukan
- berisi relasi antar model

Model tidak boleh:

- berisi business logic
- berisi HTTP response
- berisi request validation
- memanggil repository
- memanggil service
- memanggil controller

Contoh model:

```go
type User struct {
    ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
    Name      string    `gorm:"type:varchar(255)" json:"name"`
    Email     string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`
    Password  string    `gorm:"type:varchar(255)" json:"-"`
    Role      string    `gorm:"type:varchar(50)" json:"role"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

---

### 4.9 internal/request

Folder:

```txt
internal/request
```

Berisi DTO request dari client.

Tanggung jawab:

- mendefinisikan payload request
- menyimpan tag JSON
- menyimpan tag validation atau binding
- memisahkan bentuk input dari model database

Request tidak boleh:

- berisi logic database
- berisi response
- berisi business logic
- berisi GORM tag
- dipakai sebagai model database

Contoh request:

```go
type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=6"`
}
```

---

### 4.10 internal/response

Folder:

```txt
internal/response
```

Berisi DTO response dan helper response JSON.

Tanggung jawab:

- mendefinisikan bentuk response API
- membuat helper response sukses
- membuat helper response error
- memastikan format response konsisten

Contoh format sukses:

```json
{
  "status": true,
  "message": "success",
  "data": {}
}
```

Contoh format error:

```json
{
  "status": false,
  "message": "error message",
  "data": null
}
```

Contoh helper:

```go
func Success(c *gin.Context, code int, message string, data any) {
    c.JSON(code, gin.H{
        "status":  true,
        "message": message,
        "data":    data,
    })
}

func Error(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{
        "status":  false,
        "message": message,
        "data":    nil,
    })
}
```

Response tidak boleh:

- berisi query database
- berisi business logic
- memanggil repository
- memanggil service

---

### 4.11 internal/middleware

Folder:

```txt
internal/middleware
```

Berisi middleware aplikasi.

Tanggung jawab:

- authentication middleware
- authorization middleware
- role middleware
- CORS middleware
- logger middleware
- recovery middleware
- request-level validation jika diperlukan

Middleware boleh:

- membaca header authorization
- validasi token
- mengambil user ID dari token
- mengecek role user
- menghentikan request jika tidak valid

Middleware tidak boleh:

- berisi business logic utama
- query database kompleks
- memproses booking
- memproses pembayaran
- mengatur response selain response error middleware

Contoh middleware:

```go
func AuthMiddleware(jwtService *jwt.Service) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")

        if token == "" {
            response.Error(c, http.StatusUnauthorized, "unauthorized")
            c.Abort()
            return
        }

        claims, err := jwtService.ValidateToken(token)
        if err != nil {
            response.Error(c, http.StatusUnauthorized, "invalid token")
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

---

### 4.12 pkg

Folder:

```txt
pkg
```

Berisi helper reusable yang tidak terikat ke satu fitur tertentu.

Subfolder yang digunakan:

```txt
pkg/database
pkg/jwt
pkg/password
pkg/validator
```

Aturan:

- `pkg/database` untuk koneksi database
- `pkg/jwt` untuk generate dan validasi JWT
- `pkg/password` untuk hash dan compare password
- `pkg/validator` untuk custom validation helper

Tidak boleh:

- menaruh controller di pkg
- menaruh service bisnis utama di pkg
- menaruh repository fitur di pkg
- menaruh model database fitur di pkg

---

## 5. Naming Rules

Gunakan nama package sesuai nama folder.

Contoh:

```txt
internal/controller -> package controller
internal/service -> package service
internal/repository -> package repository
internal/model -> package model
internal/request -> package request
internal/response -> package response
internal/routes -> package routes
internal/middleware -> package middleware
internal/config -> package config
pkg/database -> package database
pkg/jwt -> package jwt
pkg/password -> package password
pkg/validator -> package validator
```

Gunakan nama file dengan snake_case.

Contoh:

```txt
user_controller.go
user_service.go
user_repository.go
user_model.go
auth_request.go
movie_request.go
booking_response.go
default_response.go
auth_middleware.go
role_middleware.go
```

Gunakan nama struct yang jelas.

Contoh:

```go
type UserController struct {}
type UserService struct {}
type UserRepository struct {}
type CreateUserRequest struct {}
type UserResponse struct {}
```

Gunakan constructor function.

Contoh:

```go
func NewUserController(userService *service.UserService) *UserController
func NewUserService(userRepository *repository.UserRepository) *UserService
func NewUserRepository(db *gorm.DB) *UserRepository
```

---

## 6. Coding Rules

Aturan coding umum:

- Jangan mengubah business logic tanpa alasan.
- Jangan membuat abstraction berlebihan.
- Jangan membuat interface jika belum dibutuhkan.
- Jangan membuat folder baru tanpa alasan jelas.
- Jangan mencampur controller, service, dan repository dalam satu file.
- Jangan menaruh query database di controller.
- Jangan menaruh response HTTP di repository.
- Jangan hardcode konfigurasi sensitif.
- Gunakan `.env` untuk konfigurasi.
- Gunakan UUID untuk primary key jika project sudah memakai UUID.
- Gunakan error handling yang jelas.
- Gunakan response JSON yang konsisten.
- Gunakan nama function yang mudah dipahami.
- Gunakan package import yang rapi.
- Jalankan `go fmt` setelah mengubah file Go.
- Jalankan `go mod tidy` setelah mengubah import atau dependency.

---

## 7. Database Rules

Aturan database:

- Migration file berada di `database/migrations`
- Database connection berada di `pkg/database`
- Model berada di `internal/model`
- Repository adalah satu-satunya layer yang boleh berinteraksi langsung dengan database query
- Gunakan GORM atau SQL hanya di repository
- Jangan query database langsung dari controller
- Jangan query database langsung dari routes
- Jangan query database langsung dari response
- Jangan query database langsung dari request

Jika memakai UUID:

- ID utama menggunakan UUID
- UUID dibuat di service sebelum disimpan melalui repository
- Jangan memakai auto increment integer jika project sudah disepakati memakai UUID

Contoh UUID:

```go
ID: uuid.NewString()
```

---

## 8. API Response Rules

Gunakan format response yang konsisten.

Response sukses:

```json
{
  "status": true,
  "message": "success",
  "data": {}
}
```

Response error:

```json
{
  "status": false,
  "message": "error message",
  "data": null
}
```

Aturan:

- Semua response sukses menggunakan helper dari `internal/response`
- Semua response error menggunakan helper dari `internal/response`
- Controller bertanggung jawab mengirim response
- Service hanya mengembalikan data dan error
- Repository hanya mengembalikan model dan error

---

## 9. Error Handling Rules

Aturan error:

- Repository mengembalikan error asli dari database
- Service boleh mengubah error database menjadi error bisnis
- Controller mengubah error dari service menjadi HTTP response
- Jangan panic untuk error normal
- Jangan expose error database mentah ke client jika mengandung detail sensitif

Contoh:

```go
if err != nil {
    return nil, errors.New("user not found")
}
```

---

## 10. Refactor Rules

Saat melakukan refactor:

1. Pahami struktur project terlebih dahulu.
2. Buat struktur folder target.
3. Pindahkan file sesuai layer.
4. Update package name.
5. Update import path.
6. Jalankan `go fmt`.
7. Jalankan `go mod tidy`.
8. Jalankan compile check.
9. Pastikan aplikasi bisa dijalankan.

Command validasi:

```bash
go fmt ./...
go mod tidy
go test ./...
go run cmd/api/main.go
```

Jika ada error:

- baca pesan error
- cari file yang bermasalah
- perbaiki package name
- perbaiki import path
- perbaiki dependency injection
- ulangi pengecekan sampai compile berhasil

---

## 11. Dependency Injection Rules

Dependency dibuat dari bawah ke atas:

```txt
database
repository
service
controller
routes
server
```

Contoh:

```go
db := database.Connect(cfg)

userRepository := repository.NewUserRepository(db)
userService := service.NewUserService(userRepository)
userController := controller.NewUserController(userService)

routes.SetupRoutes(router, userController)
```

Aturan:

- repository menerima database
- service menerima repository
- controller menerima service
- routes menerima controller
- main mengatur semua dependency

---

## 12. Feature Rules

Untuk setiap fitur baru, buat file sesuai layer.

Contoh fitur `movie`:

```txt
internal/model/movie_model.go
internal/request/movie_request.go
internal/response/movie_response.go
internal/repository/movie_repository.go
internal/service/movie_service.go
internal/controller/movie_controller.go
```

Jika fitur membutuhkan middleware:

```txt
internal/middleware/auth_middleware.go
internal/middleware/role_middleware.go
```

Jika fitur membutuhkan route:

```txt
internal/routes/routes.go
```

Alur implementasi fitur:

1. buat model
2. buat request
3. buat response
4. buat repository
5. buat service
6. buat controller
7. daftarkan route
8. test endpoint

---

## 13. Authentication Rules

Jika project memiliki authentication:

- logic register berada di service
- logic login berada di service
- password hashing berada di `pkg/password`
- JWT helper berada di `pkg/jwt`
- auth middleware berada di `internal/middleware`
- controller hanya menerima request login/register dan mengirim response

Tidak boleh:

- hash password langsung di controller
- generate JWT langsung di controller jika sudah ada helper
- validasi token di controller
- menyimpan secret JWT hardcoded di source code

---

## 14. Middleware Rules

Middleware digunakan untuk request-level logic.

Contoh middleware:

- authentication
- authorization
- role checking
- CORS
- logger
- recovery

Aturan:

- middleware boleh membaca header
- middleware boleh membaca token
- middleware boleh menyimpan value ke context Gin
- middleware boleh menghentikan request dengan `c.Abort()`
- middleware tidak boleh menjalankan business logic utama fitur

---

## 15. Environment Rules

Gunakan `.env` untuk konfigurasi.

Contoh:

```env
APP_PORT=:8080
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=movie_ticket_db

JWT_SECRET=secret
JWT_EXPIRED_HOUR=24
```

Aturan:

- jangan hardcode credential database
- jangan hardcode JWT secret
- jangan commit `.env` jika berisi data sensitif
- sediakan `.env.example` jika diperlukan

---

## 16. Migration Rules

Gunakan migration untuk perubahan struktur database.

Aturan:

- jangan membuat tabel manual tanpa migration jika project sudah memakai migration
- setiap perubahan tabel harus dibuatkan migration baru
- file migration disimpan di `database/migrations`
- setiap file `.up.sql` harus punya pasangan `.down.sql`

Contoh:

```txt
000001_create_users_table.up.sql
000001_create_users_table.down.sql
```

---

## 17. Git Rules

Saat melakukan perubahan besar:

- jangan ubah banyak logic sekaligus tanpa alasan
- refactor struktur dulu
- pastikan compile
- baru lanjut tambah fitur
- jangan hapus file tanpa memahami fungsinya
- jangan rename business concept sembarangan

---

## 18. Validation Commands

Setelah perubahan, jalankan:

```bash
go fmt ./...
go mod tidy
go test ./...
go run cmd/api/main.go
```

Jika project belum punya test, `go test ./...` tetap dijalankan untuk compile check.

---

## 19. Preferred Architecture Summary

Request flow:

```txt
Client Request
    ↓
Routes
    ↓
Controller
    ↓
Service
    ↓
Repository
    ↓
Database
```

Response flow:

```txt
Database
    ↓
Repository
    ↓
Service
    ↓
Controller
    ↓
Client Response
```

---

## 20. Final Rule

Prioritaskan struktur yang:

- mudah dibaca
- mudah dikembangkan
- tidak terlalu abstrak
- tidak mencampur tanggung jawab layer
- cocok untuk project training backend Go
- tetap bisa berkembang menjadi project production
````
