package worker

import (
	"log"
	"time"

	"cryptowatch/internal/service"
)

type PriceFetcher struct {
	service  *service.PriceService
	interval int
}

func NewPriceFetcher(service *service.PriceService, interval int) *PriceFetcher {
	return &PriceFetcher{
		service:  service,
		interval: interval,
	}
}

func (w *PriceFetcher) Start() {
	ticker := time.NewTicker(time.Duration(w.interval) * time.Second)
	defer ticker.Stop()

	log.Println("Price Fetcher Worker started")

	for range ticker.C {
		if err := w.service.FetchAndStore(); err != nil {
			log.Printf("Error fetching prices: %v", err)
		} else {
			log.Println("Prices updated successfully")
		}
	}
}

