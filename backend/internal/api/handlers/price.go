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

func (h *PriceHandler) GetPrices(c *gin.Context) {
	prices, err := h.service.GetAllPrices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": prices})
}

