package worker

import (
	"log"
	"time"

	"cryptowatch/internal/repository"
	"cryptowatch/internal/service"
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

func (w *VolumeMonitor) Start() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	log.Println("Volume Monitor Worker started")

	for range ticker.C {
		w.checkVolumeAlerts()
	}
}

func (w *VolumeMonitor) checkVolumeAlerts() {
	alerts, err := w.repo.GetAllAlerts()
	if err != nil {
		log.Printf("Error fetching alerts: %v", err)
		return
	}

	for _, alert := range alerts {
		if alert.AlertType != "volume" {
			continue
		}

		interval := w.getIntervalString(alert.TimeWindow)

		currentVolume, err := w.priceService.FetchKlineVolume(alert.Symbol, interval)
		if err != nil {
			log.Printf("Error fetching kline volume for %s: %v", alert.Symbol, err)
			continue
		}

		if currentVolume >= alert.TargetVolume {
			log.Printf("Volume alert triggered for %s: current %s volume %.2f (target: %.2f)",
				alert.Symbol, interval, currentVolume, alert.TargetVolume)
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
