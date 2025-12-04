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

// CreateAlert godoc
// @Summary      創建警報
// @Description  創建價格警報或成交量警報
// @Tags         alerts
// @Accept       json
// @Produce      json
// @Param        alert  body      models.CreateAlertRequest  true  "警報資訊"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Router       /alerts [post]
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

// GetUserAlerts godoc
// @Summary      獲取用戶警報
// @Description  根據用戶ID獲取所有警報
// @Tags         alerts
// @Produce      json
// @Param        userId  path      string  true  "用戶ID"
// @Success      200     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Router       /alerts/{userId} [get]
func (h *AlertHandler) GetUserAlerts(c *gin.Context) {
	userID := c.Param("userId")
	alerts, err := h.service.GetUserAlerts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": alerts})
}

// DeleteAlert godoc
// @Summary      刪除警報
// @Description  根據警報ID刪除警報
// @Tags         alerts
// @Produce      json
// @Param        alertId  path      string  true  "警報ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /alerts/{alertId} [delete]
func (h *AlertHandler) DeleteAlert(c *gin.Context) {
	alertID := c.Param("alertId")
	if err := h.service.DeleteAlert(alertID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "警報已刪除"})
}
