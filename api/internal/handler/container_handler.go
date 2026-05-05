package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/duybuidev/sentinel-20/api/internal/service"
)

type ContainerHandler struct {
	svc *service.ContainerService
}

func NewContainerHandler(svc *service.ContainerService) *ContainerHandler {
	return &ContainerHandler{svc: svc}
}

func (h *ContainerHandler) GetAll(c *gin.Context) {
	containers, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": containers, "total": len(containers)})
}

func (h *ContainerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	container, err := h.svc.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "container not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": container})
}

func (h *ContainerHandler) GetSummary(c *gin.Context) {
	summary, err := h.svc.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": summary})
}
