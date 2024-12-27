package checkouthdl

import (
	"complete-api/internal/core/domain"
	"complete-api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	checkoutSrv ports.CheckoutService
}

func NewHTTPHandler(checkoutSrv ports.CheckoutService) *HTTPHandler {
	return &HTTPHandler{
		checkoutSrv: checkoutSrv,
	}
}

func (h *HTTPHandler) Checkout(c *gin.Context) {
	var event domain.Event

	if err := c.BindJSON(&event); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.checkoutSrv.Create(event, 100); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Checkout created successfully"})
}
