package controller

import (
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"net/http"

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
// @Summary      Book Ticket
// @Description  Membooking kursi untuk jadwal tertentu. Satu request bisa booking beberapa kursi sekaligus.
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        request body request.BookTicketRequest true "Booking details"
// @Success      201  {object} response.TicketResponse
// @Failure      400  {object} map[string]interface{}
// @Failure      401  {object} map[string]interface{}
// @Router       /tickets [post]
func (t *ticketController) BookTicket(c *gin.Context) {
	var req request.BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request body"))
		return
	}

	// Ambil user_id yang di-set oleh AuthMiddleware
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}

	// user_id disimpan sebagai string di context (dari JWT claims)
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
// @Summary      Get Available Seats
// @Description  Menampilkan daftar kursi yang masih tersedia untuk jadwal tertentu.
// @Tags         ticket
// @Accept       json
// @Produce      json
// @Param        schedule_id path string true "Schedule ID"
// @Success      200  {object} []response.SeatAvailabilityResponse
// @Failure      400  {object} map[string]interface{}
// @Router       /tickets/available-seats/:schedule_id [get]
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
