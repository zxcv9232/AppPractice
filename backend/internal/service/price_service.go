package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"cryptowatch/internal/models"
	"cryptowatch/internal/repository"
)

type PriceService struct {
	repo      *repository.RedisRepository
	apiURL    string
	symbols   []string
}

func NewPriceService(repo *repository.RedisRepository, apiURL string) *PriceService {
	return &PriceService{
		repo:    repo,
		apiURL:  apiURL,
		symbols: []string{"BTC", "ETH", "BNB", "SOL", "XRP"},
	}
}

func (s *PriceService) FetchAndStore() error {
	for _, symbol := range s.symbols {
		price, err := s.fetchPriceFromBinance(symbol)
		if err != nil {
			continue
		}
		if err := s.repo.SetPrice(price); err != nil {
			continue
		}
	}
	return nil
}

func (s *PriceService) fetchPriceFromBinance(symbol string) (*models.Price, error) {
	pair := symbol + "USDT"
	url := fmt.Sprintf("%s/ticker/24hr?symbol=%s", s.apiURL, pair)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	price, _ := strconv.ParseFloat(data["lastPrice"].(string), 64)
	change, _ := strconv.ParseFloat(data["priceChangePercent"].(string), 64)
	volume, _ := strconv.ParseFloat(data["volume"].(string), 64)

	return &models.Price{
		Symbol:    symbol,
		Price:     price,
		Change24h: change,
		Volume:    volume,
		Timestamp: time.Now(),
	}, nil
}

func (s *PriceService) GetAllPrices() ([]*models.Price, error) {
	return s.repo.GetAllPrices(s.symbols)
}

