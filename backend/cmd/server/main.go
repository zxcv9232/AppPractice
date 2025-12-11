package main

import (
	"context"
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
	"golang.org/x/sync/errgroup"
)

func init() {
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
	alertMonitor := worker.NewAlertMonitor(redisRepo)
	volumeMonitor := worker.NewVolumeMonitor(redisRepo, priceService)

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return priceFetcher.Start(ctx)
	})

	g.Go(func() error {
		return alertMonitor.Start(ctx)
	})

	g.Go(func() error {
		return volumeMonitor.Start(ctx)
	})

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

	g.Go(func() error {
		log.Info().Str("port", cfg.Port).Msg("Server starting")
		return router.Run(":" + cfg.Port)
	})

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Application stopped")
	}
}
