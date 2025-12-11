package worker

import (
	"context"
	"time"

	"cryptowatch/internal/service"

	"github.com/rs/zerolog/log"
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

func (w *PriceFetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Duration(w.interval) * time.Second)
	defer ticker.Stop()

	log.Info().Msg("Price Fetcher Worker started")

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Price Fetcher Worker stopped")
			return ctx.Err()
		case <-ticker.C:
		if err := w.service.FetchAndStore(); err != nil {
				log.Error().Err(err).Msg("Error fetching prices")
		} else {
				log.Info().Msg("Prices updated successfully")
			}
		}
	}
}

