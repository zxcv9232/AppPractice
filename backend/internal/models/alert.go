package models

import "time"

type Alert struct {
	AlertID       string    `json:"alertId"`
	UserID        string    `json:"userId"`
	Symbol        string    `json:"symbol"`
	AlertType     string    `json:"alertType"`
	TargetPrice   float64   `json:"targetPrice,omitempty"`
	Direction     string    `json:"direction,omitempty"`
	TargetVolume  float64   `json:"targetVolume,omitempty"`
	TimeWindow    int       `json:"timeWindow,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
}

type CreateAlertRequest struct {
	UserID       string  `json:"userId" binding:"required"`
	Symbol       string  `json:"symbol" binding:"required"`
	AlertType    string  `json:"alertType" binding:"required"`
	TargetPrice  float64 `json:"targetPrice,omitempty"`
	Direction    string  `json:"direction,omitempty"`
	TargetVolume float64 `json:"targetVolume,omitempty"`
	TimeWindow   int     `json:"timeWindow,omitempty"`
}
