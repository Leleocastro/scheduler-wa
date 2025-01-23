package checkouthdl

import (
	"complete-api/internal/core/domain"
	"complete-api/internal/core/ports"
	"log"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	checkoutSrv ports.CheckoutService
	paymentSrv  ports.PaymentService
}

func NewHTTPHandler(checkoutSrv ports.CheckoutService, paymentSrv ports.PaymentService) *HTTPHandler {
	return &HTTPHandler{
		checkoutSrv: checkoutSrv,
		paymentSrv:  paymentSrv,
	}
}

func (h *HTTPHandler) Checkout(c *gin.Context) {
	var event domain.Event

	if err := c.BindJSON(&event); err != nil {
		log.Printf("Erro ao fazer o bind do JSON: %v", err)
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	plan, err := h.paymentSrv.GetPlanByPriceID(event.Data.Object.Metadata.PriceID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := h.checkoutSrv.Create(event, plan); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Checkout created successfully"})
}
