package models

import "time"

// IndicatorSubscription 用戶對特定幣種的指標監控訂閱
type IndicatorSubscription struct {
	SubscriptionID string    `json:"subscriptionId"`
	UserID         string    `json:"userId"`
	Symbol         string    `json:"symbol"`  // BTC, ETH, etc.
	Enabled        bool      `json:"enabled"` // 主開關

	// Telegram 通知設定
	TelegramChatID string `json:"telegramChatId"` // Telegram Chat ID

	// 通知設定
	NotifyIntervalMin int `json:"notifyIntervalMin"` // 通知間隔（分鐘），預設 60

	// 成交量判斷設定
	EnableVolumeCheck bool    `json:"enableVolumeCheck"` // 成交量開關
	VolumeCheckMode   string  `json:"volumeCheckMode"`   // "fixed" 或 "multiplier"
	VolumeFixedValue  float64 `json:"volumeFixedValue"`  // 固定值模式：閾值
	VolumeMultiplier  float64 `json:"volumeMultiplier"`  // 倍數模式：N 倍
	VolumeAvgPeriod   int     `json:"volumeAvgPeriod"`   // 均量計算週期（幾根 K 線）

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateSubscriptionRequest 創建訂閱請求
type CreateSubscriptionRequest struct {
	UserID            string  `json:"userId" binding:"required"`
	Symbol            string  `json:"symbol" binding:"required"`
	TelegramChatID    string  `json:"telegramChatId" binding:"required"` // Telegram Chat ID
	NotifyIntervalMin int     `json:"notifyIntervalMin"`                 // 預設 60
	EnableVolumeCheck bool    `json:"enableVolumeCheck"`
	VolumeCheckMode   string  `json:"volumeCheckMode"`  // "fixed" 或 "multiplier"
	VolumeFixedValue  float64 `json:"volumeFixedValue"` // 固定值模式用
	VolumeMultiplier  float64 `json:"volumeMultiplier"` // 倍數模式用
	VolumeAvgPeriod   int     `json:"volumeAvgPeriod"`  // 預設 20
}

// UpdateSubscriptionRequest 更新訂閱請求
type UpdateSubscriptionRequest struct {
	Enabled           *bool    `json:"enabled"`
	TelegramChatID    *string  `json:"telegramChatId"`
	NotifyIntervalMin *int     `json:"notifyIntervalMin"`
	EnableVolumeCheck *bool    `json:"enableVolumeCheck"`
	VolumeCheckMode   *string  `json:"volumeCheckMode"`
	VolumeFixedValue  *float64 `json:"volumeFixedValue"`
	VolumeMultiplier  *float64 `json:"volumeMultiplier"`
	VolumeAvgPeriod   *int     `json:"volumeAvgPeriod"`
}

// ApplyDefaults 套用預設值
func (r *CreateSubscriptionRequest) ApplyDefaults() {
	if r.NotifyIntervalMin <= 0 {
		r.NotifyIntervalMin = 60 // 預設 1 小時
	}
	if r.VolumeCheckMode == "" {
		r.VolumeCheckMode = "multiplier"
	}
	if r.VolumeMultiplier <= 0 {
		r.VolumeMultiplier = 2.0
	}
	if r.VolumeAvgPeriod <= 0 {
		r.VolumeAvgPeriod = 20
	}
}

