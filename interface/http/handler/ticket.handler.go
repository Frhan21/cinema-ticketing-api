package handler

import (
	"cinema-ticketing-api/app/seat"
	"cinema-ticketing-api/app/ticket"
	"cinema-ticketing-api/interface/http/httpx"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/response"
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
	ticketService ticket.TicketService
}

func NewTicketController(ticketService ticket.TicketService) TicketController {
	return &ticketController{ticketService: ticketService}
}

// BookTicket godoc
// @Summary      Booking tiket
// @Description  Memesan satu atau lebih kursi untuk jadwal film tertentu. Bisa menyertakan kode promo opsional untuk mendapatkan diskon. Tiket dibuat dengan status 'pending' dan akan otomatis dibatalkan jika tidak dibayar dalam 15 menit.
// @Tags         Ticket
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  ticket.BookTicketRequest  true  "Detail booking tiket"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /ticket [post]
func (t *ticketController) BookTicket(c *gin.Context) {
	var req ticket.BookTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	userID, ok := httpx.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}

	transaction, ticketIDs, err := t.ticketService.BookTicket(userID, &req)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	ticketRes := ticket.TicketResponse{
		TransactionId: transaction.ID,
		TicketIds:     ticketIDs,
		TotalPrice:    transaction.TotalPrice,
		PaymentStatus: string(transaction.PaymentStatus),
	}

	c.JSON(http.StatusCreated, response.SuccessResponse("Ticket booked successfully", ticketRes))
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
		c.Error(apperror.NewBadRequestError("Invalid schedule_id format"))
		c.Abort()
		return
	}

	availableSeats, err := t.ticketService.GetAvailableSeats(scheduleID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var res []seat.SeatAvailabilityResponse
	for _, s := range availableSeats {
		res = append(res, seat.SeatAvailabilityResponse{
			ID:         s.ID.String(),
			SeatNumber: s.SeatNumber,
		})
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Available seats retrieved successfully", res))
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
	userID, ok := httpx.GetUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}

	history, err := t.ticketService.GetUserHistory(userID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var res []ticket.TransactionHistoryResponse
	for _, h := range history {
		res = append(res, ticket.TransactionHistoryResponse{
			ID:            h.ID.String(),
			TotalPrice:    h.TotalPrice,
			PaymentStatus: string(h.PaymentStatus),
			CreatedAt:     h.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, response.SuccessResponse("User history retrieved successfully", res))
}
