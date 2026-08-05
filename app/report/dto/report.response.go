package dto

// DailyReportResponse adalah struktur laporan penjualan harian.
type DailyReportResponse struct {
	Date             string              `json:"date"`
	TotalTicketsSold int                 `json:"total_tickets_sold"`
	TotalRevenue     float64             `json:"total_revenue"`
	ByMovie          []MovieSalesReport  `json:"by_movie"`
	ByStudio         []StudioSalesReport `json:"by_studio"`
}

// MonthlyReportResponse adalah struktur laporan penjualan bulanan.
type MonthlyReportResponse struct {
	Year             int                 `json:"year"`
	Month            int                 `json:"month"`
	TotalTicketsSold int                 `json:"total_tickets_sold"`
	TotalRevenue     float64             `json:"total_revenue"`
	ByMovie          []MovieSalesReport  `json:"by_movie"`
	ByStudio         []StudioSalesReport `json:"by_studio"`
	DailyBreakdown   []DailySales        `json:"daily_breakdown"`
}

// MovieSalesReport ringkasan penjualan per film.
type MovieSalesReport struct {
	MovieID     string  `json:"movie_id"`
	MovieTitle  string  `json:"movie_title"`
	TicketsSold int     `json:"tickets_sold"`
	Revenue     float64 `json:"revenue"`
}

// StudioSalesReport ringkasan penjualan per studio.
type StudioSalesReport struct {
	StudioID    string  `json:"studio_id"`
	StudioName  string  `json:"studio_name"`
	TicketsSold int     `json:"tickets_sold"`
	Revenue     float64 `json:"revenue"`
}

// DailySales ringkasan penjualan per hari (untuk laporan bulanan).
type DailySales struct {
	Date        string  `json:"date"`
	TicketsSold int     `json:"tickets_sold"`
	Revenue     float64 `json:"revenue"`
}
