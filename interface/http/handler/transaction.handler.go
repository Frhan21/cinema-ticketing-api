package handler

import (
	ticketservice "cinema-ticketing-api/app/ticket/service"
	"cinema-ticketing-api/interface/http/httpx"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionController interface {
	GetAll(c *gin.Context)
	GetByID(c *gin.Context)
	CancelTransaction(c *gin.Context)
}

type transactionController struct {
	transactionService ticketservice.TransactionService
}

func NewTransactionController(transactionService ticketservice.TransactionService) TransactionController {
	return &transactionController{transactionService: transactionService}
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
		c.Error(err)
		c.Abort()
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
	userID, ok := httpx.GetUserIDFromContext(c)
	if !ok {
		c.Error(apperror.NewUnauthorizedError("Unauthorized"))
		c.Abort()
		return
	}
	role, _ := httpx.GetRoleFromContext(c)

	transactionID, err := uuid.Parse(c.Param("transaction_id"))
	if err != nil {
		c.Error(apperror.NewBadRequestError("Invalid transaction_id format"))
		c.Abort()
		return
	}
	transaction, err := t.transactionService.GetByID(userID, role, transactionID)
	if err != nil {
		c.Error(err)
		c.Abort()
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
	userID, ok := httpx.GetUserIDFromContext(c)
	if !ok {
		c.Error(apperror.NewUnauthorizedError("Unauthorized"))
		c.Abort()
		return
	}
	transactionID, err := uuid.Parse(c.Param("transaction_id"))
	if err != nil {
		c.Error(apperror.NewBadRequestError("Invalid transaction_id format"))
		c.Abort()
		return
	}
	if err := t.transactionService.CancelTransaction(userID, transactionID); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, response.SuccessResponse("Transaction cancelled successfully", nil))
}
