package handlers

import (
	"net/http"

	"cryptowatch/internal/models"
	"cryptowatch/internal/service"
	"cryptowatch/internal/worker"

	"github.com/gin-gonic/gin"
)

// IndicatorHandler 指標 API 處理器
type IndicatorHandler struct {
	subscriptionService *service.SubscriptionService
	indicatorMonitor    *worker.IndicatorMonitor
}

// NewIndicatorHandler 創建指標處理器
func NewIndicatorHandler(
	subscriptionService *service.SubscriptionService,
	indicatorMonitor *worker.IndicatorMonitor,
) *IndicatorHandler {
	return &IndicatorHandler{
		subscriptionService: subscriptionService,
		indicatorMonitor:    indicatorMonitor,
	}
}

// CreateSubscription 創建訂閱
// @Summary      創建指標監控訂閱
// @Description  訂閱特定幣種的 LRC 突破警報
// @Tags         indicators
// @Accept       json
// @Produce      json
// @Param        request body models.CreateSubscriptionRequest true "訂閱資料"
// @Success      201 {object} models.IndicatorSubscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /indicators/subscribe [post]
func (h *IndicatorHandler) CreateSubscription(c *gin.Context) {
	var req models.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.subscriptionService.CreateSubscription(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

// GetUserSubscriptions 獲取用戶訂閱
// @Summary      獲取用戶的所有訂閱
// @Description  獲取指定用戶的所有指標監控訂閱
// @Tags         indicators
// @Produce      json
// @Param        userId query string true "用戶 ID"
// @Success      200 {array} models.IndicatorSubscription
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /indicators/subscriptions [get]
func (h *IndicatorHandler) GetUserSubscriptions(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	subs, err := h.subscriptionService.GetUserSubscriptions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

// UpdateSubscription 更新訂閱
// @Summary      更新訂閱設定
// @Description  更新指標監控訂閱的設定（開關、通知間隔、成交量條件等）
// @Tags         indicators
// @Accept       json
// @Produce      json
// @Param        id path string true "訂閱 ID"
// @Param        request body models.UpdateSubscriptionRequest true "更新資料"
// @Success      200 {object} models.IndicatorSubscription
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /indicators/subscriptions/{id} [put]
func (h *IndicatorHandler) UpdateSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription id is required"})
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.subscriptionService.UpdateSubscription(subscriptionID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// DeleteSubscription 刪除訂閱
// @Summary      刪除訂閱
// @Description  刪除指標監控訂閱
// @Tags         indicators
// @Param        id path string true "訂閱 ID"
// @Success      204
// @Failure      500 {object} map[string]string
// @Router       /indicators/subscriptions/{id} [delete]
func (h *IndicatorHandler) DeleteSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription id is required"})
		return
	}

	if err := h.subscriptionService.DeleteSubscription(subscriptionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ToggleSubscription 切換訂閱開關
// @Summary      切換訂閱開關
// @Description  快速切換訂閱的啟用/停用狀態
// @Tags         indicators
// @Produce      json
// @Param        id path string true "訂閱 ID"
// @Success      200 {object} models.IndicatorSubscription
// @Failure      500 {object} map[string]string
// @Router       /indicators/subscriptions/{id}/toggle [post]
func (h *IndicatorHandler) ToggleSubscription(c *gin.Context) {
	subscriptionID := c.Param("id")
	if subscriptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subscription id is required"})
		return
	}

	sub, err := h.subscriptionService.ToggleSubscription(subscriptionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

// GetIndicatorResult 獲取指標結果
// @Summary      獲取幣種當前指標值
// @Description  獲取指定幣種的 LRC 指標計算結果
// @Tags         indicators
// @Produce      json
// @Param        symbol path string true "幣種代號"
// @Success      200 {object} models.IndicatorResult
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /indicators/{symbol} [get]
func (h *IndicatorHandler) GetIndicatorResult(c *gin.Context) {
	symbol := c.Param("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol is required"})
		return
	}

	result, err := h.indicatorMonitor.GetIndicatorResult(symbol)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

