package models

import "time"

type Price struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Change24h float64 `json:"change24h"`
	Volume    float64 `json:"volume"`
	Timestamp time.Time `json:"timestamp"`
}

