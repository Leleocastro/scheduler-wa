package scheduleshdl

import (
	"complete-api/internal/core/domain"
	"complete-api/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	schedulesSrv ports.SchedulesService
}

func NewHTTPHandler(schedulesSrv ports.SchedulesService) *HTTPHandler {
	return &HTTPHandler{
		schedulesSrv: schedulesSrv,
	}
}

func (h *HTTPHandler) GetSchedules(c *gin.Context) {

	username := c.Request.Header.Get("X-Consumer-Username")

	schedules, err := h.schedulesSrv.GetSchedules(username)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, schedules)

}

func (h *HTTPHandler) CreateSchedule(c *gin.Context) {
	var schedule domain.ScheduleMessage

	if err := c.BindJSON(&schedule); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Regras de validação
	if schedule.Repeats > 0 && schedule.CronExpr == "" {
		c.JSON(400, gin.H{"error": "Field 'cronExpr' is required when 'repeats' is set"})
		return
	}

	if schedule.Until > 0 && schedule.CronExpr == "" {
		c.JSON(400, gin.H{"error": "Field 'cronExpr' is required when 'until' is set"})
		return
	}

	if schedule.Repeats > 0 && schedule.Until > 0 {
		c.JSON(400, gin.H{"error": "Fields 'repeats' and 'until' cannot be used together"})
		return
	}

	username := c.Request.Header.Get("X-Consumer-Username")
	if username == "" {
		c.JSON(400, gin.H{"error": "Missing X-Consumer-Username header"})
		return
	}

	newSchedule, err := h.schedulesSrv.CreateSchedule(username, schedule)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, newSchedule)
}

func (h *HTTPHandler) UpdateSchedule(c *gin.Context) {
	scheduleID := c.Param("scheduleID")
	var schedule domain.ScheduleMessage

	if err := c.BindJSON(&schedule); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON format"})
		return
	}

	username := c.Request.Header.Get("X-Consumer-Username")

	updatedSchedule, err := h.schedulesSrv.UpdateSchedule(username, scheduleID, schedule)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, updatedSchedule)
}

func (h *HTTPHandler) DeleteSchedule(c *gin.Context) {
	scheduleID := c.Param("scheduleID")

	username := c.Request.Header.Get("X-Consumer-Username")

	if err := h.schedulesSrv.DeleteSchedule(username, scheduleID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, "Schedule deleted successfully")
}
