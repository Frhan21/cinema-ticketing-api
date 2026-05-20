package controller

import (
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController interface {
	GetDailyReport(c *gin.Context)
	GetMonthlyReport(c *gin.Context)
}

type reportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) ReportController {
	return &reportController{reportService: reportService}
}

// GetDailyReport godoc
// @Summary      Laporan Penjualan Harian
// @Description  Mengembalikan laporan penjualan tiket pada tanggal tertentu (admin only).
// @Tags         report
// @Produce      json
// @Param        date query string false "Tanggal (YYYY-MM-DD), default: hari ini"
// @Success      200  {object} response.DailyReportResponse
// @Router       /report/daily [get]
func (r *reportController) GetDailyReport(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		// Default ke hari ini (WIB)
		date = time.Now().In(time.FixedZone("WIB", 7*3600)).Format("2006-01-02")
	}

	report, err := r.reportService.GetDailyReport(date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Daily report retrieved successfully", report))
}

// GetMonthlyReport godoc
// @Summary      Laporan Penjualan Bulanan
// @Description  Mengembalikan laporan penjualan tiket pada bulan tertentu (admin only).
// @Tags         report
// @Produce      json
// @Param        year  query string false "Tahun (YYYY), default: tahun ini"
// @Param        month query string false "Bulan (1-12), default: bulan ini"
// @Success      200  {object} response.MonthlyReportResponse
// @Router       /report/monthly [get]
func (r *reportController) GetMonthlyReport(c *gin.Context) {
	now := time.Now().In(time.FixedZone("WIB", 7*3600))

	year := c.DefaultQuery("year", now.Format("2006"))
	month := c.DefaultQuery("month", now.Format("1"))

	report, err := r.reportService.GetMonthlyReport(year, month)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Monthly report retrieved successfully", report))
}
