package worker

import (
	"log"
	"time"

	"cryptowatch/internal/repository"
)

type VolumeMonitor struct {
	repo *repository.RedisRepository
}

func NewVolumeMonitor(repo *repository.RedisRepository) *VolumeMonitor {
	return &VolumeMonitor{repo: repo}
}

func (w *VolumeMonitor) Start() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("Volume Monitor Worker started")

	for range ticker.C {
		w.trackVolumes()
		w.checkVolumeAlerts()
	}
}

func (w *VolumeMonitor) trackVolumes() {
	symbols := []string{"BTC", "ETH", "BNB", "SOL", "XRP"}
	
	for _, symbol := range symbols {
		price, err := w.repo.GetPrice(symbol)
		if err != nil {
			continue
		}

		if err := w.repo.SaveVolumeSnapshot(symbol, price.Volume); err != nil {
			log.Printf("Error saving volume snapshot for %s: %v", symbol, err)
		}
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

		timeWindow := alert.TimeWindow
		if timeWindow == 0 {
			timeWindow = 1
		}

		cumulativeVolume, err := w.repo.GetCumulativeVolume(alert.Symbol, timeWindow)
		if err != nil {
			continue
		}

		if cumulativeVolume >= alert.TargetVolume {
			log.Printf("Volume alert triggered for %s: cumulative volume %.2f in %d minutes (target: %.2f)",
				alert.Symbol, cumulativeVolume, timeWindow, alert.TargetVolume)
			w.repo.DeleteAlert(alert.AlertID)
		}
	}
}

