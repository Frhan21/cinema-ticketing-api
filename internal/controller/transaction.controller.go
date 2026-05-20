package controller

import (
	"cinema-ticketing-api/internal/request"
	"cinema-ticketing-api/internal/response"
	"cinema-ticketing-api/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionController interface {
	PayTransaction(c *gin.Context)
	CancelTransaction(c *gin.Context)
}

type transactionController struct {
	transactionService service.TransactionService
}

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionController{transactionService: transactionService}
}

// getUserIDFromContext adalah helper untuk mengambil dan mem-parse user_id dari Gin context.
func getUserIDFromContext(c *gin.Context) (uuid.UUID, bool) {
	raw, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}

	str, ok := raw.(string)
	if !ok {
		return uuid.Nil, false
	}

	parsed, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, false
	}

	return parsed, true
}

// CancelTransaction membatalkan transaksi berdasarkan transaction_id.
func (t *transactionController) CancelTransaction(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}

	transactionID, err := uuid.Parse(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid transaction_id format"))
		return
	}

	if err := t.transactionService.CancelTransaction(userID, transactionID); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction cancelled successfully", nil))
}

// PayTransaction membayar transaksi berdasarkan transaction_id.
func (t *transactionController) PayTransaction(c *gin.Context) {
	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Unauthorized"))
		return
	}

	transactionID, err := uuid.Parse(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid transaction_id format"))
		return
	}

	var req request.PayTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid request body"))
		return
	}

	if err := t.transactionService.PayTransaction(userID, transactionID, req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Transaction paid successfully", nil))
}
