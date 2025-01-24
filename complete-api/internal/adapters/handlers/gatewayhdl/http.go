package gatewayhdl

import (
	"complete-api/internal/core/domain"
	"complete-api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	gatewaySrv ports.GatewayService
}

func NewHTTPHandler(gatewaySrv ports.GatewayService) *HTTPHandler {
	return &HTTPHandler{
		gatewaySrv: gatewaySrv,
	}
}

func (h *HTTPHandler) CreateConsumer(c *gin.Context) {
	var customer domain.Consumer

	if err := c.BindJSON(&customer); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	if err := h.gatewaySrv.CreateConsumer(customer.Username, customer.CustomID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Consumer created successfully"})
}

func (h *HTTPHandler) GetAPIKey(c *gin.Context) {
	username := c.Query("username")

	apiKey, err := h.gatewaySrv.GetAPIKey(username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"api_key": apiKey})
}
