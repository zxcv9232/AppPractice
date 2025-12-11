package worker

import (
	"context"
	"time"

	"cryptowatch/internal/repository"

	"github.com/rs/zerolog/log"
)

type AlertMonitor struct {
	repo *repository.RedisRepository
}

func NewAlertMonitor(repo *repository.RedisRepository) *AlertMonitor {
	return &AlertMonitor{repo: repo}
}

func (w *AlertMonitor) Start(ctx context.Context) error {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Info().Msg("Alert Monitor Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Alert Monitor Worker stopped")
			return ctx.Err()
		case <-ticker.C:
		w.checkAlerts()
		}
	}
}

func (w *AlertMonitor) checkAlerts() {
	alerts, err := w.repo.GetAllAlerts()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching alerts")
		return
	}

	for _, alert := range alerts {
		if alert.AlertType != "price" && alert.AlertType != "" {
			continue
		}

		price, err := w.repo.GetPrice(alert.Symbol)
		if err != nil {
			continue
		}

		shouldTrigger := false
		if alert.Direction == "above" && price.Price >= alert.TargetPrice {
			shouldTrigger = true
		} else if alert.Direction == "below" && price.Price <= alert.TargetPrice {
			shouldTrigger = true
		}

		if shouldTrigger {
			log.Info().
				Str("symbol", alert.Symbol).
				Float64("current_price", price.Price).
				Float64("target_price", alert.TargetPrice).
				Msg("Alert triggered")
			w.repo.DeleteAlert(alert.AlertID)
		}
	}
}
