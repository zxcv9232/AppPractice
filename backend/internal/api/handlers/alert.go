package handlers

import (
	"net/http"

	"cryptowatch/internal/models"
	"cryptowatch/internal/service"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	service *service.AlertService
}

func NewAlertHandler(service *service.AlertService) *AlertHandler {
	return &AlertHandler{service: service}
}

func (h *AlertHandler) CreateAlert(c *gin.Context) {
	var req models.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	alert, err := h.service.CreateAlert(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"alertId": alert.AlertID,
		"message": "警報創建成功",
	})
}

func (h *AlertHandler) GetUserAlerts(c *gin.Context) {
	userID := c.Param("userId")
	alerts, err := h.service.GetUserAlerts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": alerts})
}

func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	alertID := c.Param("alertId")
	if err := h.service.DeleteAlert(alertID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "警報已刪除"})
}

