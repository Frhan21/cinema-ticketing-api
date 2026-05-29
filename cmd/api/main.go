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
	"cinema-ticketing-api/cmd/setup"
	_ "cinema-ticketing-api/docs"
)

func main() {
	setup.InitApp()
}
