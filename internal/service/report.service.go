package service

import (
	"cinema-ticketing-api/internal/repository"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/pkg/apperror"
	"strconv"
)

type ReportService interface {
	GetDailyReport(date string) (*response.DailyReportResponse, error)
	GetMonthlyReport(year, month string) (*response.MonthlyReportResponse, error)
}

type reportService struct {
	reportRepo repository.ReportRepository
}

func NewReportService(reportRepo repository.ReportRepository) ReportService {
	return &reportService{reportRepo: reportRepo}
}

// GetDailyReport mengambil laporan penjualan untuk tanggal tertentu.
// date: format "2006-01-02" (contoh: "2026-05-20")
func (s *reportService) GetDailyReport(date string) (*response.DailyReportResponse, error) {
	if date == "" {
		return nil, apperror.NewBadRequestError("date is required (format: YYYY-MM-DD)")
	}
	return s.reportRepo.GetDailySalesReport(date)
}

// GetMonthlyReport mengambil laporan penjualan untuk bulan tertentu.
func (s *reportService) GetMonthlyReport(yearStr, monthStr string) (*response.MonthlyReportResponse, error) {
	year, err := strconv.Atoi(yearStr)
	if err != nil || year < 2000 {
		return nil, apperror.NewBadRequestError("invalid year (format: YYYY)")
	}

	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return nil, apperror.NewBadRequestError("invalid month (1-12)")
	}

	return s.reportRepo.GetMonthlySalesReport(year, month)
}
