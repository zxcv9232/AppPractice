package worker

import (
	"log"
	"time"

	"cryptowatch/internal/repository"
)

type AlertMonitor struct {
	repo *repository.RedisRepository
}

func NewAlertMonitor(repo *repository.RedisRepository) *AlertMonitor {
	return &AlertMonitor{repo: repo}
}

func (w *AlertMonitor) Start() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Println("Alert Monitor Worker started")

	for range ticker.C {
		w.checkAlerts()
	}
}

func (w *AlertMonitor) checkAlerts() {
	alerts, err := w.repo.GetAllAlerts()
	if err != nil {
		log.Printf("Error fetching alerts: %v", err)
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
			log.Printf("Alert triggered for %s: current price %.2f, target %.2f",
				alert.Symbol, price.Price, alert.TargetPrice)
			w.repo.DeleteAlert(alert.AlertID)
		}
	}
}
