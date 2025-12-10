package main

import (
	"os"

	"cryptowatch/config"
	_ "cryptowatch/docs"
	"cryptowatch/internal/api/handlers"
	"cryptowatch/internal/api/middleware"
	"cryptowatch/internal/repository"
	"cryptowatch/internal/service"
	"cryptowatch/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	// 設定 zerolog 全域 logger
	// 正式環境使用 JSON 格式，開發環境可改用 ConsoleWriter
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "cryptowatch").
		Logger()
}

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

	log.Info().Str("port", cfg.Port).Msg("Server starting")
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
