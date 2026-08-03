package report

import "gorm.io/gorm"

type ReportRepository interface {
	GetDailySalesReport(date string) (*DailyReportResponse, error)
	GetMonthlySalesReport(year, month int) (*MonthlyReportResponse, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// GetDailySalesReport mengambil laporan penjualan berdasarkan tanggal.
// date format: "2006-01-02"
func (r *reportRepository) GetDailySalesReport(date string) (*DailyReportResponse, error) {
	// Total tiket & revenue untuk hari ini
	type total struct {
		Count   int
		Revenue float64
	}
	var t total
	r.db.Raw(`
		SELECT COUNT(*) as count, COALESCE(SUM(t.price), 0) as revenue
		FROM tickets t
		WHERE t.status = 'paid'
		AND DATE(t.created_at AT TIME ZONE 'Asia/Jakarta') = ?
	`, date).Scan(&t)

	// Per film
	var byMovie []MovieSalesReport
	r.db.Raw(`
		SELECT
			m.id::text as movie_id,
			m.title as movie_title,
			COUNT(t.id) as tickets_sold,
			COALESCE(SUM(t.price), 0) as revenue
		FROM tickets t
		JOIN schedules s ON s.id = t.schedule_id
		JOIN movies m ON m.id = s.movie_id
		WHERE t.status = 'paid'
		AND DATE(t.created_at AT TIME ZONE 'Asia/Jakarta') = ?
		GROUP BY m.id, m.title
		ORDER BY revenue DESC
	`, date).Scan(&byMovie)

	// Per studio
	var byStudio []StudioSalesReport
	r.db.Raw(`
		SELECT
			st.id::text as studio_id,
			st.name as studio_name,
			COUNT(t.id) as tickets_sold,
			COALESCE(SUM(t.price), 0) as revenue
		FROM tickets t
		JOIN schedules s ON s.id = t.schedule_id
		JOIN studios st ON st.id = s.studio_id
		WHERE t.status = 'paid'
		AND DATE(t.created_at AT TIME ZONE 'Asia/Jakarta') = ?
		GROUP BY st.id, st.name
		ORDER BY revenue DESC
	`, date).Scan(&byStudio)

	return &DailyReportResponse{
		Date:             date,
		TotalTicketsSold: t.Count,
		TotalRevenue:     t.Revenue,
		ByMovie:          byMovie,
		ByStudio:         byStudio,
	}, nil
}

// GetMonthlySalesReport mengambil laporan penjualan untuk bulan tertentu.
func (r *reportRepository) GetMonthlySalesReport(year, month int) (*MonthlyReportResponse, error) {
	// Total
	type total struct {
		Count   int
		Revenue float64
	}
	var t total
	r.db.Raw(`
		SELECT COUNT(*) as count, COALESCE(SUM(price), 0) as revenue
		FROM tickets
		WHERE status = 'paid'
		AND EXTRACT(YEAR FROM created_at AT TIME ZONE 'Asia/Jakarta') = ?
		AND EXTRACT(MONTH FROM created_at AT TIME ZONE 'Asia/Jakarta') = ?
	`, year, month).Scan(&t)

	// Per film
	var byMovie []MovieSalesReport
	r.db.Raw(`
		SELECT
			m.id::text as movie_id,
			m.title as movie_title,
			COUNT(t.id) as tickets_sold,
			COALESCE(SUM(t.price), 0) as revenue
		FROM tickets t
		JOIN schedules s ON s.id = t.schedule_id
		JOIN movies m ON m.id = s.movie_id
		WHERE t.status = 'paid'
		AND EXTRACT(YEAR FROM t.created_at AT TIME ZONE 'Asia/Jakarta') = ?
		AND EXTRACT(MONTH FROM t.created_at AT TIME ZONE 'Asia/Jakarta') = ?
		GROUP BY m.id, m.title
		ORDER BY revenue DESC
	`, year, month).Scan(&byMovie)

	// Per studio
	var byStudio []StudioSalesReport
	r.db.Raw(`
		SELECT
			st.id::text as studio_id,
			st.name as studio_name,
			COUNT(t.id) as tickets_sold,
			COALESCE(SUM(t.price), 0) as revenue
		FROM tickets t
		JOIN schedules s ON s.id = t.schedule_id
		JOIN studios st ON st.id = s.studio_id
		WHERE t.status = 'paid'
		AND EXTRACT(YEAR FROM t.created_at AT TIME ZONE 'Asia/Jakarta') = ?
		AND EXTRACT(MONTH FROM t.created_at AT TIME ZONE 'Asia/Jakarta') = ?
		GROUP BY st.id, st.name
		ORDER BY revenue DESC
	`, year, month).Scan(&byStudio)

	// Breakdown harian
	var daily []DailySales
	r.db.Raw(`
		SELECT
			DATE(created_at AT TIME ZONE 'Asia/Jakarta')::text as date,
			COUNT(*) as tickets_sold,
			COALESCE(SUM(price), 0) as revenue
		FROM tickets
		WHERE status = 'paid'
		AND EXTRACT(YEAR FROM created_at AT TIME ZONE 'Asia/Jakarta') = ?
		AND EXTRACT(MONTH FROM created_at AT TIME ZONE 'Asia/Jakarta') = ?
		GROUP BY DATE(created_at AT TIME ZONE 'Asia/Jakarta')
		ORDER BY date ASC
	`, year, month).Scan(&daily)

	return &MonthlyReportResponse{
		Year:             year,
		Month:            month,
		TotalTicketsSold: t.Count,
		TotalRevenue:     t.Revenue,
		ByMovie:          byMovie,
		ByStudio:         byStudio,
		DailyBreakdown:   daily,
	}, nil
}
