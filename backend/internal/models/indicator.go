package models

import "time"

// IndicatorResult 指標計算結果
type IndicatorResult struct {
	Symbol string `json:"symbol"`

	// LRC 結果
	UpperBand  float64 `json:"upperBand"`
	LowerBand  float64 `json:"lowerBand"`
	CenterLine float64 `json:"centerLine"`
	Slope      float64 `json:"slope"`
	Deviation  float64 `json:"deviation"`

	// 當前價格
	CurrentPrice float64 `json:"currentPrice"`

	// 1 分 K 成交量
	CurrentVolume float64 `json:"currentVolume"` // 當前 1 分 K
	AvgVolume     float64 `json:"avgVolume"`     // 近 N 根平均
	VolumeRatio   float64 `json:"volumeRatio"`   // 當前/平均 比值

	// 狀態
	IsAboveUpper bool `json:"isAboveUpper"`
	IsBelowLower bool `json:"isBelowLower"`

	// 計算時間
	CalculatedAt time.Time `json:"calculatedAt"`
}

// IndicatorConfig 系統級的指標參數配置
type IndicatorConfig struct {
	// 監控的幣種清單
	Symbols []string `json:"symbols"`

	// 交易所設定（現貨/合約）
	MarketType string `json:"marketType"` // "spot" 或 "futures"

	// LRC 固定參數
	LRCLength        int     `json:"lrcLength"`        // 預設 42
	LRCDevMultiplier float64 `json:"lrcDevMultiplier"` // 預設 2.0
	LRCInterval      string  `json:"lrcInterval"`      // 預設 "4h"

	// 1 分 K 成交量預設參數
	DefaultVolumeAvgPeriod int `json:"defaultVolumeAvgPeriod"` // 預設 20
}

// DefaultIndicatorConfig 返回預設配置
func DefaultIndicatorConfig() IndicatorConfig {
	return IndicatorConfig{
		Symbols: []string{
			"BTC",      // 比特幣
			"ETH",      // 以太坊
			"BNB",      // 幣安幣
			"SOL",      // Solana
			"XRP",      // 瑞波幣
			"DOGE",     // 狗狗幣
			"ADA",      // Cardano
			"AVAX",     // Avalanche
			"1000SHIB", // Shiba Inu (1000倍)
			"BCH",      // Bitcoin Cash
			"DOT",      // Polkadot
			"LINK",     // Chainlink
			"TON",      // Toncoin
			"UNI",      // Uniswap
			"LTC",      // Litecoin
			"NEAR",     // NEAR Protocol
			"ATOM",     // Cosmos
			"AAVE",     // Aave
			"RIVER",    // River
		},
		MarketType:             "futures", // U本位永續合約
		LRCLength:              42,
		LRCDevMultiplier:       2.0,
		LRCInterval:            "4h",
		DefaultVolumeAvgPeriod: 20,
	}
}

// AlertPayload 推播通知內容
type AlertPayload struct {
	Title        string            `json:"title"`
	Body         string            `json:"body"`
	Symbol       string            `json:"symbol"`
	Type         string            `json:"type"` // "above_upper", "below_lower"
	CurrentPrice float64           `json:"currentPrice"`
	UpperBand    float64           `json:"upperBand"`
	LowerBand    float64           `json:"lowerBand"`
	Data         map[string]string `json:"data,omitempty"`
}

