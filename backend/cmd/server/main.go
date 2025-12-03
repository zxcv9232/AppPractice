package main

import (
	"log"

	"cryptowatch/config"
	"cryptowatch/internal/api/handlers"
	"cryptowatch/internal/api/middleware"
	"cryptowatch/internal/repository"
	"cryptowatch/internal/service"
	"cryptowatch/internal/worker"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	redisRepo := repository.NewRedisRepository(cfg.RedisURL)

	priceService := service.NewPriceService(redisRepo, cfg.BinanceAPIURL)
	alertService := service.NewAlertService(redisRepo)

	priceHandler := handlers.NewPriceHandler(priceService)
	alertHandler := handlers.NewAlertHandler(alertService)

	priceFetcher := worker.NewPriceFetcher(priceService, cfg.PriceFetchInterval)
	go priceFetcher.Start()

	alertMonitor := worker.NewAlertMonitor(redisRepo)
	go alertMonitor.Start()

	volumeMonitor := worker.NewVolumeMonitor(redisRepo)
	go volumeMonitor.Start()

	router := gin.Default()
	router.Use(middleware.CORS())

	api := router.Group("/api")
	{
		api.GET("/prices", priceHandler.GetPrices)
		api.POST("/alerts", alertHandler.CreateAlert)
		api.GET("/alerts/:userId", alertHandler.GetUserAlerts)
		api.DELETE("/alerts/:alertId", alertHandler.DeleteAlert)
	}

	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
