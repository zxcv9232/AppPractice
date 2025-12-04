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

func (s *PriceService) FetchKlineVolume(symbol string, interval string) (float64, error) {
	pair := symbol + "USDT"
	url := fmt.Sprintf("%s/klines?symbol=%s&interval=%s&limit=1", s.apiURL, pair, interval)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var errorResp map[string]interface{}
	if json.Unmarshal(body, &errorResp) == nil {
		if code, ok := errorResp["code"]; ok {
			msg := errorResp["msg"]
			return 0, fmt.Errorf("binance error %v: %v", code, msg)
		}
	}

	var klines [][]interface{}
	if err := json.Unmarshal(body, &klines); err != nil {
		return 0, err
	}

	if len(klines) == 0 {
		return 0, fmt.Errorf("no kline data for %s", symbol)
	}

	volumeStr, ok := klines[0][5].(string)
	if !ok {
		return 0, fmt.Errorf("invalid volume format for %s", symbol)
	}

	return strconv.ParseFloat(volumeStr, 64)
}

