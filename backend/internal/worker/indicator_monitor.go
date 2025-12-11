package worker

import (
	"context"
	"fmt"
	"time"

	"cryptowatch/internal/indicators"
	"cryptowatch/internal/models"
	"cryptowatch/internal/repository"
	"cryptowatch/internal/service"

	"github.com/rs/zerolog/log"
)

// IndicatorMonitor 技術指標監控器
type IndicatorMonitor struct {
	repo            *repository.RedisRepository
	priceService    *service.PriceService
	telegramService *service.TelegramService
	config          models.IndicatorConfig
}

// NewIndicatorMonitor 創建指標監控器
func NewIndicatorMonitor(
	repo *repository.RedisRepository,
	priceService *service.PriceService,
	telegramService *service.TelegramService,
) *IndicatorMonitor {
	return &IndicatorMonitor{
		repo:            repo,
		priceService:    priceService,
		telegramService: telegramService,
		config:          models.DefaultIndicatorConfig(),
	}
}

// Start 啟動監控
func (w *IndicatorMonitor) Start(ctx context.Context) error {
	// 每 30 秒檢查一次
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	log.Info().Msg("Indicator Monitor Worker started")

	// 啟動時先執行一次
	w.checkAndNotify()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Indicator Monitor Worker stopped")
			return ctx.Err()
		case <-ticker.C:
			w.checkAndNotify()
		}
	}
}

// checkAndNotify 檢查指標並發送通知
func (w *IndicatorMonitor) checkAndNotify() {
	// 獲取系統配置
	config, err := w.repo.GetIndicatorConfig()
	if err != nil {
		log.Error().Err(err).Msg("Error getting indicator config")
		config = &w.config
	}

	for _, symbol := range config.Symbols {
		// 計算指標
		result, err := w.calculateIndicators(symbol, config)
		if err != nil {
			log.Error().Err(err).Str("symbol", symbol).Msg("Error calculating indicators")
			continue
		}

		// 快取結果
		w.repo.SetIndicatorResult(result)

		// 檢查是否突破 LRC（必要條件）
		if !result.IsAboveUpper && !result.IsBelowLower {
			continue
		}

		// 獲取該幣種的所有訂閱者
		subscriptions, err := w.repo.GetSubscriptionsBySymbol(symbol)
		if err != nil {
			log.Error().Err(err).Str("symbol", symbol).Msg("Error getting subscriptions")
			continue
		}

		for _, sub := range subscriptions {
			// 檢查開關
			if !sub.Enabled {
				continue
			}

			// 檢查冷卻時間
			if w.isInCooldown(sub.SubscriptionID, sub.NotifyIntervalMin) {
				continue
			}

			// 檢查成交量條件（如果啟用）
			if sub.EnableVolumeCheck {
				if !w.checkVolumeCondition(result, sub) {
					continue
				}
			}

			// 發送通知
			w.sendNotification(sub, result)

			// 記錄通知時間
			w.recordNotification(sub.SubscriptionID)
		}
	}
}

// calculateIndicators 計算指標
func (w *IndicatorMonitor) calculateIndicators(symbol string, config *models.IndicatorConfig) (*models.IndicatorResult, error) {
	// 嘗試從快取獲取
	cached, err := w.repo.GetIndicatorResult(symbol)
	if err == nil && cached != nil {
		// 快取有效
		return cached, nil
	}

	// 獲取 4H K 線數據（用於 LRC 計算）
	lrcKlines, err := w.priceService.FetchKlines(symbol, config.LRCInterval, config.LRCLength+5)
	if err != nil {
		return nil, fmt.Errorf("error fetching LRC klines: %v", err)
	}

	// 計算 LRC
	closePrices := service.GetClosePrices(lrcKlines)
	lrc, err := indicators.CalculateLRC(closePrices, config.LRCLength, config.LRCDevMultiplier)
	if err != nil {
		return nil, fmt.Errorf("error calculating LRC: %v", err)
	}

	// 獲取當前價格
	currentPrice, err := w.priceService.FetchCurrentPrice(symbol)
	if err != nil {
		return nil, fmt.Errorf("error getting current price: %v", err)
	}

	// 獲取 1 分 K 成交量數據
	volumeKlines, err := w.priceService.FetchKlines(symbol, "1m", config.DefaultVolumeAvgPeriod+5)
	if err != nil {
		log.Warn().Err(err).Str("symbol", symbol).Msg("Error fetching volume klines, skipping volume calculation")
		// 成交量獲取失敗不影響主要功能
	}

	var volumeResult indicators.VolumeResult
	if len(volumeKlines) > 0 {
		volumes := service.GetVolumes(volumeKlines)
		volumeResult = indicators.CalculateVolumeStats(volumes, config.DefaultVolumeAvgPeriod)
	}

	result := &models.IndicatorResult{
		Symbol:        symbol,
		UpperBand:     lrc.UpperBand,
		LowerBand:     lrc.LowerBand,
		CenterLine:    lrc.CenterLine,
		Slope:         lrc.Slope,
		Deviation:     lrc.Deviation,
		CurrentPrice:  currentPrice,
		CurrentVolume: volumeResult.CurrentVolume,
		AvgVolume:     volumeResult.AvgVolume,
		VolumeRatio:   volumeResult.VolumeRatio,
		IsAboveUpper:  currentPrice > lrc.UpperBand,
		IsBelowLower:  currentPrice < lrc.LowerBand,
		CalculatedAt:  time.Now(),
	}

	return result, nil
}

// checkVolumeCondition 檢查成交量條件
func (w *IndicatorMonitor) checkVolumeCondition(result *models.IndicatorResult, sub *models.IndicatorSubscription) bool {
	config := indicators.VolumeConfig{
		Enabled:    sub.EnableVolumeCheck,
		Mode:       indicators.VolumeCheckMode(sub.VolumeCheckMode),
		FixedValue: sub.VolumeFixedValue,
		Multiplier: sub.VolumeMultiplier,
		AvgPeriod:  sub.VolumeAvgPeriod,
	}

	volumeResult := indicators.VolumeResult{
		CurrentVolume: result.CurrentVolume,
		AvgVolume:     result.AvgVolume,
		VolumeRatio:   result.VolumeRatio,
	}

	return indicators.CheckVolumeCondition(volumeResult, config)
}

// isInCooldown 檢查是否在冷卻時間內
func (w *IndicatorMonitor) isInCooldown(subscriptionID string, intervalMin int) bool {
	key := fmt.Sprintf("indicator_notify:%s", subscriptionID)
	lastNotify, err := w.repo.GetLastNotifyTime(key)
	if err != nil {
		return false // 如果取不到，表示還沒通知過
	}

	cooldown := time.Duration(intervalMin) * time.Minute
	return time.Since(lastNotify) < cooldown
}

// recordNotification 記錄通知時間
func (w *IndicatorMonitor) recordNotification(subscriptionID string) {
	key := fmt.Sprintf("indicator_notify:%s", subscriptionID)
	w.repo.SetLastNotifyTime(key)
}

// sendNotification 發送通知
func (w *IndicatorMonitor) sendNotification(sub *models.IndicatorSubscription, result *models.IndicatorResult) {
	direction := "突破上軌 📈"
	alertType := "above_upper"
	if result.IsBelowLower {
		direction = "跌破下軌 📉"
		alertType = "below_lower"
	}

	payload := models.AlertPayload{
		Title:        fmt.Sprintf("🚨 %s %s", result.Symbol, direction),
		Body:         fmt.Sprintf("價格 %.2f | 上軌 %.2f | 下軌 %.2f", result.CurrentPrice, result.UpperBand, result.LowerBand),
		Symbol:       result.Symbol,
		Type:         alertType,
		CurrentPrice: result.CurrentPrice,
		UpperBand:    result.UpperBand,
		LowerBand:    result.LowerBand,
	}

	// 如果有成交量判斷，加入成交量資訊
	if sub.EnableVolumeCheck && result.CurrentVolume > 0 {
		payload.Body += fmt.Sprintf(" | 成交量 %.2f (%.1fx)", result.CurrentVolume, result.VolumeRatio)
	}

	if err := w.telegramService.SendAlert(sub.TelegramChatID, payload); err != nil {
		log.Error().
			Err(err).
			Str("symbol", result.Symbol).
			Str("subscriptionId", sub.SubscriptionID).
			Msg("Error sending Telegram notification")
	} else {
		log.Info().
			Str("symbol", result.Symbol).
			Str("userId", sub.UserID).
			Str("direction", direction).
			Float64("price", result.CurrentPrice).
			Msg("Telegram alert sent")
	}
}

// GetIndicatorResult 獲取指標結果（供 API 使用）
func (w *IndicatorMonitor) GetIndicatorResult(symbol string) (*models.IndicatorResult, error) {
	// 嘗試從快取獲取
	cached, err := w.repo.GetIndicatorResult(symbol)
	if err == nil && cached != nil {
		return cached, nil
	}

	// 計算新的結果
	config, _ := w.repo.GetIndicatorConfig()
	if config == nil {
		defaultConfig := models.DefaultIndicatorConfig()
		config = &defaultConfig
	}

	return w.calculateIndicators(symbol, config)
}

