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

// MarketType 市場類型
type MarketType string

const (
	MarketTypeSpot    MarketType = "spot"
	MarketTypeFutures MarketType = "futures"
)

// KlineData 表示單根 K 線數據
type KlineData struct {
	OpenTime  int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	CloseTime int64
}

type PriceService struct {
	repo       *repository.RedisRepository
	apiURL     string
	futuresURL string
	marketType MarketType
	symbols    []string
}

func NewPriceService(repo *repository.RedisRepository, apiURL string) *PriceService {
	return &PriceService{
		repo:       repo,
		apiURL:     apiURL,
		futuresURL: "https://fapi.binance.com/fapi/v1",
		marketType: MarketTypeFutures, // 預設合約（U本位永續）
		symbols: []string{
			"BTC", "ETH", "BNB", "SOL", "XRP",
			"DOGE", "ADA", "AVAX", "1000SHIB", "BCH",
			"DOT", "LINK", "TON", "UNI", "LTC",
			"NEAR", "ATOM", "AAVE", "RIVER",
		},
	}
}

// SetMarketType 設置市場類型（現貨或合約）
func (s *PriceService) SetMarketType(marketType MarketType) {
	s.marketType = marketType
}

// GetSymbols 返回監控的幣種清單
func (s *PriceService) GetSymbols() []string {
	return s.symbols
}

// SetSymbols 設置監控的幣種清單
func (s *PriceService) SetSymbols(symbols []string) {
	s.symbols = symbols
}

// getBaseURL 根據市場類型返回對應的 API URL
func (s *PriceService) getBaseURL() string {
	if s.marketType == MarketTypeFutures {
		return s.futuresURL
	}
	return s.apiURL
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
	baseURL := s.getBaseURL()
	url := fmt.Sprintf("%s/ticker/24hr?symbol=%s", baseURL, pair)

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
	baseURL := s.getBaseURL()
	url := fmt.Sprintf("%s/klines?symbol=%s&interval=%s&limit=1", baseURL, pair, interval)

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

// FetchKlines 獲取歷史 K 線數據
// symbol: 交易對 (例如 "BTC")
// interval: K 線週期 (例如 "4h", "1m")
// limit: 獲取的 K 線數量
func (s *PriceService) FetchKlines(symbol string, interval string, limit int) ([]KlineData, error) {
	pair := symbol + "USDT"
	baseURL := s.getBaseURL()
	url := fmt.Sprintf("%s/klines?symbol=%s&interval=%s&limit=%d", baseURL, pair, interval, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 檢查錯誤響應
	var errorResp map[string]interface{}
	if json.Unmarshal(body, &errorResp) == nil {
		if code, ok := errorResp["code"]; ok {
			msg := errorResp["msg"]
			return nil, fmt.Errorf("binance error %v: %v", code, msg)
		}
	}

	var rawKlines [][]interface{}
	if err := json.Unmarshal(body, &rawKlines); err != nil {
		return nil, err
	}

	klines := make([]KlineData, 0, len(rawKlines))
	for _, k := range rawKlines {
		if len(k) < 7 {
			continue
		}
		openTime, _ := k[0].(float64)
		openStr, _ := k[1].(string)
		highStr, _ := k[2].(string)
		lowStr, _ := k[3].(string)
		closeStr, _ := k[4].(string)
		volumeStr, _ := k[5].(string)
		closeTime, _ := k[6].(float64)

		open, _ := strconv.ParseFloat(openStr, 64)
		high, _ := strconv.ParseFloat(highStr, 64)
		low, _ := strconv.ParseFloat(lowStr, 64)
		closePrice, _ := strconv.ParseFloat(closeStr, 64)
		volume, _ := strconv.ParseFloat(volumeStr, 64)

		klines = append(klines, KlineData{
			OpenTime:  int64(openTime),
			Open:      open,
			High:      high,
			Low:       low,
			Close:     closePrice,
			Volume:    volume,
			CloseTime: int64(closeTime),
		})
	}

	return klines, nil
}

// GetClosePrices 從 K 線數據中提取收盤價
func GetClosePrices(klines []KlineData) []float64 {
	prices := make([]float64, len(klines))
	for i, k := range klines {
		prices[i] = k.Close
	}
	return prices
}

// GetVolumes 從 K 線數據中提取成交量
func GetVolumes(klines []KlineData) []float64 {
	volumes := make([]float64, len(klines))
	for i, k := range klines {
		volumes[i] = k.Volume
	}
	return volumes
}

// FetchCurrentPrice 獲取當前價格（從 Redis 快取或 API）
func (s *PriceService) FetchCurrentPrice(symbol string) (float64, error) {
	// 先嘗試從 Redis 獲取
	price, err := s.repo.GetPrice(symbol)
	if err == nil && price != nil {
		return price.Price, nil
	}

	// 從 API 獲取
	priceData, err := s.fetchPriceFromBinance(symbol)
	if err != nil {
		return 0, err
	}
	return priceData.Price, nil
}
