// Package main adalah entry point dari Cinema Ticketing API.
//
// @title           Cinema Ticketing API
// @version         1.0
// @description     REST API untuk sistem manajemen tiket bioskop. Mendukung manajemen studio, film, jadwal tayang, kursi, booking tiket, transaksi, promo diskon, dan laporan penjualan.
// @termsOfService  http://swagger.io/terms/
//
// @contact.name    Cinema Ticketing Support
// @contact.email   support@cinemticket.com
//
// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT
//
// @host      localhost:8080
// @BasePath  /api/v1
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"cinema-ticketing-api/cmd/setup"
	"cinema-ticketing-api/config"
	_ "cinema-ticketing-api/docs"
	"cinema-ticketing-api/pkg/database"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// migrationDir adalah path relatif ke folder migrasi, diukur dari root repository.
// CLI migrasi dan auto-migration saat startup sama-sama memakai path ini.
const migrationDir = "./migration"

var invalidFileNameChars = regexp.MustCompile(`[^a-z0-9]+`)

// main menjalankan CLI migrasi ketika argumen pertama adalah subcommand migrasi;
// selain itu menjalankan API server seperti biasa.
//
// Penggunaan CLI (dari root repository):
//
//	go run cmd/api/main.go create <nama>   // buat file migrasi baru (up/down)
//	go run cmd/api/main.go up              // jalankan semua migrasi pending
//	go run cmd/api/main.go down [n]        // rollback n migrasi terakhir (default 1)
//	go run cmd/api/main.go version         // tampilkan versi migrasi aktif
//	go run cmd/api/main.go force <versi>   // paksa set versi (pemulihan state dirty)
func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "create", "up", "down", "version", "force":
			runMigrationCommand(os.Args[1], os.Args[2:])
			return
		case "help", "--help", "-h":
			migrationUsage()
			return
		}
	}

	setup.InitApp()
}

// runMigrationCommand mendispatch subcommand migrasi ke handler masing-masing.
func runMigrationCommand(command string, args []string) {
	switch command {
	case "create":
		runCreate(args)
	case "up":
		runUp()
	case "down":
		runDown(args)
	case "version":
		runVersion()
	case "force":
		runForce(args)
	}
}

// runCreate membuat pasangan file migrasi <timestamp>_<nama>.up.sql dan .down.sql.
func runCreate(args []string) {
	if len(args) < 1 || strings.TrimSpace(args[0]) == "" {
		fmt.Fprintln(os.Stderr, "usage: go run cmd/api/main.go create <nama_migrasi>")
		os.Exit(1)
	}

	name := sanitizeFileName(args[0])
	version := time.Now().Format("20060102150405") // YYYYMMDDHHMMSS, konsisten dengan migrasi yang ada
	base := filepath.Join(migrationDir, fmt.Sprintf("%s_%s", version, name))
	upPath := base + ".up.sql"
	downPath := base + ".down.sql"

	if err := os.MkdirAll(migrationDir, 0o755); err != nil {
		log.Fatalf("failed to create migration directory: %v", err)
	}
	if fileExists(upPath) || fileExists(downPath) {
		log.Fatalf("migration already exists: %s", base)
	}

	template := "-- Migration: %s\n"
	if err := os.WriteFile(upPath, []byte(fmt.Sprintf(template, name)), 0o644); err != nil {
		log.Fatalf("failed to create %s: %v", upPath, err)
	}
	if err := os.WriteFile(downPath, []byte(fmt.Sprintf(template, name)), 0o644); err != nil {
		log.Fatalf("failed to create %s: %v", downPath, err)
	}

	fmt.Printf("created:\n  %s\n  %s\n", upPath, downPath)
}

// runUp menerapkan semua migrasi yang belum dijalankan.
func runUp() {
	m := newMigrate()
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migration to apply")
			return
		}
		log.Fatalf("migration up failed: %v", err)
	}
	fmt.Println("all migrations applied")
}

// runDown me-rollback n migrasi terakhir (default 1).
func runDown(args []string) {
	steps := 1
	if len(args) > 0 {
		parsed, err := strconv.Atoi(args[0])
		if err != nil || parsed < 1 {
			fmt.Fprintf(os.Stderr, "invalid step count: %q (must be a positive integer)\n", args[0])
			os.Exit(1)
		}
		steps = parsed
	}

	m := newMigrate()
	if err := m.Steps(-steps); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migration to roll back")
			return
		}
		log.Fatalf("migration down failed: %v", err)
	}
	fmt.Printf("rolled back %d migration(s)\n", steps)
}

// runVersion menampilkan versi migrasi aktif dan status dirty.
func runVersion() {
	m := newMigrate()
	version, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		fmt.Println("no migration applied yet")
		return
	}
	if err != nil {
		log.Fatalf("failed to read migration version: %v", err)
	}
	fmt.Printf("current version: %d (dirty: %v)\n", version, dirty)
}

// runForce memaksa set versi migrasi tanpa menjalankan migrasi (untuk recovery).
func runForce(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: go run cmd/api/main.go force <versi>")
		os.Exit(1)
	}
	version, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		log.Fatalf("invalid version %q: must be a non-negative integer", args[0])
	}

	m := newMigrate()
	if err := m.Force(int(version)); err != nil {
		log.Fatalf("force version failed: %v", err)
	}
	fmt.Printf("forced schema version to %d\n", version)
}

// newMigrate membuat instance golang-migrate dengan konfigurasi dari .env.
func newMigrate() *migrate.Migrate {
	cfg := config.Load()
	m, err := migrate.New("file://"+migrationDir, database.PostgresURL(cfg.DB))
	if err != nil {
		log.Fatalf("failed to initialize migration: %v", err)
	}
	return m
}

func sanitizeFileName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = invalidFileNameChars.ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	if name == "" {
		return "migration"
	}
	return name
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func migrationUsage() {
	fmt.Print(`Migration CLI untuk Cinema Ticketing API

Penggunaan (dari root repository):
  go run cmd/api/main.go create <nama>     buat file migrasi baru (up/down) di migration/
  go run cmd/api/main.go up                jalankan semua migrasi yang belum diterapkan
  go run cmd/api/main.go down [n]          rollback n migrasi terakhir (default 1)
  go run cmd/api/main.go version           tampilkan versi migrasi aktif
  go run cmd/api/main.go force <versi>     paksa set versi migrasi (untuk pemulihan state dirty)

Tanpa argumen, aplikasi berjalan sebagai API server.
Koneksi database dibaca dari .env (DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME).
`)
}
