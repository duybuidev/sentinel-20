package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/duybuidev/sentinel-20/api/internal/repository"
	"github.com/duybuidev/sentinel-20/api/internal/service"
)

type LogHandler struct {
	svc *service.LogService
}

func NewLogHandler(svc *service.LogService) *LogHandler {
	return &LogHandler{svc: svc}
}

func (h *LogHandler) GetAll(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	logs, err := h.svc.GetAll(repository.LogFilter{
		Level:   c.Query("level"),
		Service: c.Query("service"),
		Limit:   limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": logs, "total": len(logs)})
}

func (h *LogHandler) GetStats(c *gin.Context) {
	stats, err := h.svc.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
