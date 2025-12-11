package worker

import (
	"context"
	"time"

	"cryptowatch/internal/repository"
	"cryptowatch/internal/service"

	"github.com/rs/zerolog/log"
)

type VolumeMonitor struct {
	repo         *repository.RedisRepository
	priceService *service.PriceService
}

func NewVolumeMonitor(repo *repository.RedisRepository, priceService *service.PriceService) *VolumeMonitor {
	return &VolumeMonitor{
		repo:         repo,
		priceService: priceService,
	}
}

func (w *VolumeMonitor) Start(ctx context.Context) error {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Info().Msg("Volume Monitor Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Volume Monitor Worker stopped")
			return ctx.Err()
		case <-ticker.C:
		w.checkVolumeAlerts()
		}
	}
}

func (w *VolumeMonitor) checkVolumeAlerts() {
	alerts, err := w.repo.GetAllAlerts()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching alerts")
		return
	}

	for _, alert := range alerts {
		if alert.AlertType != "volume" {
			continue
		}

		interval := w.getIntervalString(alert.TimeWindow)

		currentVolume, err := w.priceService.FetchKlineVolume(alert.Symbol, interval)
		if err != nil {
			log.Error().
				Err(err).
				Str("symbol", alert.Symbol).
				Msg("Error fetching kline volume")
			continue
		}

		if currentVolume >= alert.TargetVolume {
			log.Info().
				Str("symbol", alert.Symbol).
				Str("interval", interval).
				Float64("current_volume", currentVolume).
				Float64("target_volume", alert.TargetVolume).
				Msg("Volume alert triggered")
			w.repo.DeleteAlert(alert.AlertID)
		}
	}
}

func (w *VolumeMonitor) getIntervalString(minutes int) string {
	switch minutes {
	case 1:
		return "1m"
	case 3:
		return "3m"
	case 5:
		return "5m"
	case 15:
		return "15m"
	case 30:
		return "30m"
	case 60:
		return "1h"
	default:
		return "1m"
	}
}
