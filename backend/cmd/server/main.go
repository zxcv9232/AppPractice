package main

import (
	"log"

	"cryptowatch/config"
	_ "cryptowatch/docs"
	"cryptowatch/internal/api/handlers"
	"cryptowatch/internal/api/middleware"
	"cryptowatch/internal/repository"
	"cryptowatch/internal/service"
	"cryptowatch/internal/worker"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           CryptoWatch API
// @version         1.0
// @description     加密貨幣價格監控與警報系統 API
// @host            localhost:8080
// @BasePath        /api

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

	volumeMonitor := worker.NewVolumeMonitor(redisRepo, priceService)
	go volumeMonitor.Start()

	router := gin.Default()
	router.Use(middleware.CORS())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
