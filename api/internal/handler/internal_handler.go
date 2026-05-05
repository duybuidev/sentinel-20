package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/duybuidev/sentinel-20/api/internal/model"
	"github.com/duybuidev/sentinel-20/api/internal/service"
)

type InternalHandler struct {
	containerSvc *service.ContainerService
	eventSvc     *service.EventService
}

func NewInternalHandler(containerSvc *service.ContainerService, eventSvc *service.EventService) *InternalHandler {
	return &InternalHandler{containerSvc: containerSvc, eventSvc: eventSvc}
}

func (h *InternalHandler) SyncContainer(c *gin.Context) {
	var payload struct {
		Name       string  `json:"name"`
		Image      string  `json:"image"`
		Status     string  `json:"status"`
		CPUPercent float64 `json:"cpu_percent"`
		MemUsage   int64   `json:"mem_usage"`
		MemLimit   int64   `json:"mem_limit"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	container := &model.Container{
		Name:       payload.Name,
		Image:      payload.Image,
		Status:     payload.Status,
		CPUPercent: payload.CPUPercent,
		MemUsage:   payload.MemUsage,
		MemLimit:   payload.MemLimit,
	}

	if err := h.containerSvc.Upsert(container); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

func (h *InternalHandler) CreateEvent(c *gin.Context) {
	var payload struct {
		ContainerName string `json:"container_name"`
		Type          string `json:"type"`
		ExitCode      int    `json:"exit_code"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Tìm container ID từ tên
	container, err := h.containerSvc.GetByName(payload.ContainerName)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
		return
	}

	event := &model.Event{
		ContainerID: container.ID,
		Type:        payload.Type,
		ExitCode:    payload.ExitCode,
		Description: payload.ContainerName + " " + payload.Type,
	}

	if err := h.eventSvc.Insert(event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
