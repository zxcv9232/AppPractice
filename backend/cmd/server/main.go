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

	// 現有服務
	priceService := service.NewPriceService(redisRepo, cfg.BinanceAPIURL)
	alertService := service.NewAlertService(redisRepo)

	// Telegram 通知服務
	telegramService := service.NewTelegramService(
		cfg.TelegramBotToken,
		cfg.TelegramTestMode,
		cfg.TelegramMyChatID,
	)

	// 訂閱服務
	subscriptionService := service.NewSubscriptionService(redisRepo)

	// 現有 handlers
	priceHandler := handlers.NewPriceHandler(priceService)
	alertHandler := handlers.NewAlertHandler(alertService)

	// 現有 workers
	priceFetcher := worker.NewPriceFetcher(priceService, cfg.PriceFetchInterval)
	alertMonitor := worker.NewAlertMonitor(redisRepo)
	volumeMonitor := worker.NewVolumeMonitor(redisRepo, priceService)

	// 指標監控 worker
	indicatorMonitor := worker.NewIndicatorMonitor(redisRepo, priceService, telegramService)

	// 新增：指標 handler
	indicatorHandler := handlers.NewIndicatorHandler(subscriptionService, indicatorMonitor)

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

	// 新增：啟動指標監控
	g.Go(func() error {
		return indicatorMonitor.Start(ctx)
	})

	router := gin.Default()
	router.Use(middleware.CORS())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		// 現有路由
		api.GET("/prices", priceHandler.GetPrices)
		api.POST("/alerts", alertHandler.CreateAlert)
		api.GET("/alerts/:userId", alertHandler.GetUserAlerts)
		api.DELETE("/alerts/:alertId", alertHandler.DeleteAlert)

		// 新增：指標監控路由
		indicators := api.Group("/indicators")
		{
			indicators.POST("/subscribe", indicatorHandler.CreateSubscription)
			indicators.GET("/subscriptions", indicatorHandler.GetUserSubscriptions)
			indicators.PUT("/subscriptions/:id", indicatorHandler.UpdateSubscription)
			indicators.DELETE("/subscriptions/:id", indicatorHandler.DeleteSubscription)
			indicators.POST("/subscriptions/:id/toggle", indicatorHandler.ToggleSubscription)
			indicators.GET("/:symbol", indicatorHandler.GetIndicatorResult)
		}
	}

	g.Go(func() error {
		log.Info().Str("port", cfg.Port).Msg("Server starting")
		return router.Run(":" + cfg.Port)
	})

	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("Application stopped")
	}
}
