package controller

import (
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketController interface {
	BookTicket(c *gin.Context)
	GetAvailableSeats(c *gin.Context)
	GetUserHistory(c *gin.Context)
}

type ticketController struct {
	ticketService service.TicketService
}

func NewTicketController(ticketService service.TicketService) TicketController {
	return &ticketController{ticketService: ticketService}
}

// BookTicket godoc
// @Summary      Booking tiket
// @Description  Memesan satu atau lebih kursi untuk jadwal film tertentu. Bisa menyertakan kode promo opsional untuk mendapatkan diskon. Tiket dibuat dengan status 'pending' dan akan otomatis dibatalkan jika tidak dibayar dalam 15 menit.
// @Tags         Ticket
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  request.BookTicketRequest  true  "Detail booking tiket"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /ticket [post]
func (t *ticketController) BookTicket(c *gin.Context) {
	var req request.BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request body"))
		return
	}

	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid user ID format"))
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid user ID"))
		return
	}

	ticket, err := t.ticketService.BookTicket(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Ticket booked successfully", ticket))
}

// GetAvailableSeats godoc
// @Summary      Kursi tersedia
// @Description  Menampilkan daftar kursi yang belum dipesan untuk jadwal film tertentu.
// @Tags         Ticket
// @Produce      json
// @Security     BearerAuth
// @Param        schedule_id  path  string  true  "Schedule ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /ticket/available-seats/{schedule_id} [get]
func (t *ticketController) GetAvailableSeats(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("schedule_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid schedule_id format"))
		return
	}

	availableSeats, err := t.ticketService.GetAvailableSeats(scheduleID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Available seats retrieved successfully", availableSeats))
}

// GetUserHistory godoc
// @Summary      Riwayat transaksi user
// @Description  Menampilkan semua riwayat pembelian tiket milik user yang sedang login.
// @Tags         Ticket
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /ticket/history [get]
func (t *ticketController) GetUserHistory(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid user ID format"))
		return
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid user ID"))
		return
	}

	history, err := t.ticketService.GetUserHistory(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User history retrieved successfully", history))
}

// getUserIDFromContextTicket digunakan internal — alias getUserIDFromContext di ticket controller.
func getUserIDFromContextTicket(c *gin.Context) (uuid.UUID, bool) {
	_ = time.Now() // import guard
	return getUserIDFromContext(c)
}
