package statshdl

import (
	"complete-api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	statsSrv ports.StatsService
}

func NewHTTPHandler(statsSrv ports.StatsService) *HTTPHandler {
	return &HTTPHandler{
		statsSrv: statsSrv,
	}
}

func (h *HTTPHandler) GetUsageByConsumer(c *gin.Context) {
	username := c.Param("consumer")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if username == "" {
		c.JSON(400, gin.H{"error": "username is required"})
	}

	res, err := h.statsSrv.GetUsageByConsumer(username, startDate, endDate)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": res})
}
