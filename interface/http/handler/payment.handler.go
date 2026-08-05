package handler

import (
	"cinema-ticketing-api/app/payment/dto"
	paymentservice "cinema-ticketing-api/app/payment/service"
	"cinema-ticketing-api/interface/http/httpx"
	"cinema-ticketing-api/pkg/apperror"
	"cinema-ticketing-api/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentController interface {
	CreatePayment(c *gin.Context)
	HandleNotification(c *gin.Context)
}

type paymentController struct {
	paymentService paymentservice.PaymentService
}

func NewPaymentController(paymentService paymentservice.PaymentService) PaymentController {
	return &paymentController{paymentService: paymentService}
}

// CreatePayment godoc
// @Summary      Buat pembayaran Midtrans
// @Description  Membuat Snap transaction dan mengembalikan URL pembayaran Midtrans.
// @Tags         Payment
// @Produce      json
// @Security     BearerAuth
// @Param        transaction_id  path  string  true  "Transaction ID (UUID)"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /transaction/{transaction_id}/pay [post]
func (p *paymentController) CreatePayment(c *gin.Context) {
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

	payment, err := p.paymentService.CreatePayment(userID, transactionID)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Payment created successfully", payment))
}

// HandleNotification godoc
// @Summary      Midtrans payment notification
// @Description  Receives a Midtrans webhook and verifies its status directly with Midtrans.
// @Tags         Payment
// @Accept       json
// @Produce      json
// @Param        body  body  dto.MidtransNotificationRequest  true  "Midtrans notification"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /payment/notification [post]
func (p *paymentController) HandleNotification(c *gin.Context) {
	var req dto.MidtransNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(apperror.NewBadRequestError("Invalid notification payload"))
		c.Abort()
		return
	}

	if err := p.paymentService.HandleNotification(req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, response.SuccessResponse("Payment notification processed", nil))
}
