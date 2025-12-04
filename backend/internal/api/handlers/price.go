package handlers

import (
	"net/http"

	"cryptowatch/internal/service"

	"github.com/gin-gonic/gin"
)

type PriceHandler struct {
	service *service.PriceService
}

func NewPriceHandler(service *service.PriceService) *PriceHandler {
	return &PriceHandler{service: service}
}

// GetPrices godoc
// @Summary      獲取所有加密貨幣價格
// @Description  返回 BTC, ETH, BNB, SOL, XRP 等幣種的即時價格
// @Tags         prices
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /prices [get]
func (h *PriceHandler) GetPrices(c *gin.Context) {
	prices, err := h.service.GetAllPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": prices})
}
