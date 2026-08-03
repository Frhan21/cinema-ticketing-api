package handler

import (
	"cinema-ticketing-api/app/report"
	"cinema-ticketing-api/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportController interface {
	GetDailyReport(c *gin.Context)
	GetMonthlyReport(c *gin.Context)
}

type reportController struct {
	reportService report.ReportService
}

func NewReportController(reportService report.ReportService) ReportController {
	return &reportController{reportService: reportService}
}

// GetDailyReport godoc
// @Summary      Laporan penjualan harian
// @Description  Admin dapat melihat laporan penjualan tiket pada tanggal tertentu, termasuk total tiket, total pendapatan, dan breakdown per film dan studio.
// @Tags         Report
// @Produce      json
// @Security     BearerAuth
// @Param        date  query  string  false  "Tanggal laporan (format: YYYY-MM-DD). Default: hari ini"  example(2026-05-20)
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Failure      403   {object}  map[string]interface{}
// @Router       /report/daily [get]
func (r *reportController) GetDailyReport(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		date = time.Now().In(time.FixedZone("WIB", 7*3600)).Format("2006-01-02")
	}
	report, err := r.reportService.GetDailyReport(date)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Daily report retrieved successfully", report))
}

// GetMonthlyReport godoc
// @Summary      Laporan penjualan bulanan
// @Description  Admin dapat melihat laporan penjualan tiket selama satu bulan penuh, termasuk breakdown harian, per film, dan per studio.
// @Tags         Report
// @Produce      json
// @Security     BearerAuth
// @Param        year   query  string  false  "Tahun laporan (format: YYYY). Default: tahun ini"   example(2026)
// @Param        month  query  string  false  "Bulan laporan (1-12). Default: bulan ini"           example(5)
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Failure      403    {object}  map[string]interface{}
// @Router       /report/monthly [get]
func (r *reportController) GetMonthlyReport(c *gin.Context) {
	now := time.Now().In(time.FixedZone("WIB", 7*3600))
	year := c.DefaultQuery("year", now.Format("2006"))
	month := c.DefaultQuery("month", now.Format("1"))
	report, err := r.reportService.GetMonthlyReport(year, month)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Monthly report retrieved successfully", report))
}
