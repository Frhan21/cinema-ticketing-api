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
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
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

// GetAll godoc
// @Summary      Lihat semua transaksi
// @Description  Mengembalikan seluruh data transaksi. Hanya bisa diakses oleh admin.
// @Tags         Transaction
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /transaction [get]
func (t *transactionController) GetAll(c *gin.Context) {
	transactions, err := t.transactionService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Transactions retrieved successfully", transactions))
}

// GetByID godoc
// @Summary      Detail transaksi
// @Description  Mengembalikan detail satu transaksi berdasarkan ID-nya.
// @Tags         Transaction
// @Produce      json
// @Security     BearerAuth
// @Param        transaction_id  path  string  true  "Transaction ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /transaction/{transaction_id} [get]
func (t *transactionController) GetByID(c *gin.Context) {
	transactionID, err := uuid.Parse(c.Param("transaction_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("Invalid transaction_id format"))
		return
	}
	transaction, err := t.transactionService.GetByID(transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Transaction retrieved successfully", transaction))
}

// CancelTransaction godoc
// @Summary      Batalkan transaksi
// @Description  Membatalkan transaksi yang masih berstatus pending. Semua tiket terkait juga ikut dibatalkan dan email notifikasi dikirim ke user.
// @Tags         Transaction
// @Produce      json
// @Security     BearerAuth
// @Param        transaction_id  path  string  true  "Transaction ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /transaction/{transaction_id}/cancel [post]
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

// PayTransaction godoc
// @Summary      Bayar transaksi
// @Description  Melakukan pembayaran untuk transaksi yang masih berstatus pending. Tiket berubah menjadi 'paid' dan email konfirmasi dikirim.
// @Tags         Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        transaction_id  path  string                        true  "Transaction ID (UUID)"
// @Param        body            body  request.PayTransactionRequest  true  "Metode pembayaran"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /transaction/{transaction_id}/pay [post]
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
