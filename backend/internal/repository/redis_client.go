package repository

import (
	"context"
	"encoding/json"
	"time"

	"cryptowatch/internal/models"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisRepository(addr string) *RedisRepository {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisRepository) SetPrice(price *models.Price) error {
	data, err := json.Marshal(price)
	if err != nil {
		return err
	}
	key := "prices:" + price.Symbol
	return r.client.Set(r.ctx, key, data, 20*time.Second).Err()
}

func (r *RedisRepository) GetPrice(symbol string) (*models.Price, error) {
	key := "prices:" + symbol
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var price models.Price
	if err := json.Unmarshal([]byte(data), &price); err != nil {
		return nil, err
	}
	return &price, nil
}

func (r *RedisRepository) GetAllPrices(symbols []string) ([]*models.Price, error) {
	prices := make([]*models.Price, 0, len(symbols))
	for _, symbol := range symbols {
		price, err := r.GetPrice(symbol)
		if err == nil {
			prices = append(prices, price)
		}
	}
	return prices, nil
}

func (r *RedisRepository) SaveAlert(alert *models.Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return err
	}
	key := "alert:" + alert.AlertID
	return r.client.Set(r.ctx, key, data, 0).Err()
}

func (r *RedisRepository) GetUserAlerts(userID string) ([]*models.Alert, error) {
	pattern := "alert:*"
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	alerts := make([]*models.Alert, 0)
	for _, key := range keys {
		data, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			continue
		}
		var alert models.Alert
		if err := json.Unmarshal([]byte(data), &alert); err != nil {
			continue
		}
		if alert.UserID == userID {
			alerts = append(alerts, &alert)
		}
	}
	return alerts, nil
}

func (r *RedisRepository) DeleteAlert(alertID string) error {
	key := "alert:" + alertID
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisRepository) GetAllAlerts() ([]*models.Alert, error) {
	pattern := "alert:*"
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	alerts := make([]*models.Alert, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			continue
		}
		var alert models.Alert
		if err := json.Unmarshal([]byte(data), &alert); err != nil {
			continue
		}
		alerts = append(alerts, &alert)
	}
	return alerts, nil
}

// ==================== 指標訂閱相關方法 ====================

// SaveSubscription 儲存訂閱
func (r *RedisRepository) SaveSubscription(sub *models.IndicatorSubscription) error {
	data, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	// 儲存訂閱資料
	key := "indicator_sub:" + sub.SubscriptionID
	if err := r.client.Set(r.ctx, key, data, 0).Err(); err != nil {
		return err
	}

	// 加入用戶的訂閱集合
	userKey := "indicator_subs:user:" + sub.UserID
	if err := r.client.SAdd(r.ctx, userKey, sub.SubscriptionID).Err(); err != nil {
		return err
	}

	// 加入幣種的訂閱集合
	symbolKey := "indicator_subs:symbol:" + sub.Symbol
	return r.client.SAdd(r.ctx, symbolKey, sub.SubscriptionID).Err()
}

// GetSubscription 獲取單個訂閱
func (r *RedisRepository) GetSubscription(subscriptionID string) (*models.IndicatorSubscription, error) {
	key := "indicator_sub:" + subscriptionID
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var sub models.IndicatorSubscription
	if err := json.Unmarshal([]byte(data), &sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

// GetUserSubscriptions 獲取用戶的所有訂閱
func (r *RedisRepository) GetUserSubscriptions(userID string) ([]*models.IndicatorSubscription, error) {
	userKey := "indicator_subs:user:" + userID
	subIDs, err := r.client.SMembers(r.ctx, userKey).Result()
	if err != nil {
		return nil, err
	}

	subs := make([]*models.IndicatorSubscription, 0, len(subIDs))
	for _, subID := range subIDs {
		sub, err := r.GetSubscription(subID)
		if err != nil {
			continue
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// GetSubscriptionsBySymbol 獲取特定幣種的所有訂閱
func (r *RedisRepository) GetSubscriptionsBySymbol(symbol string) ([]*models.IndicatorSubscription, error) {
	symbolKey := "indicator_subs:symbol:" + symbol
	subIDs, err := r.client.SMembers(r.ctx, symbolKey).Result()
	if err != nil {
		return nil, err
	}

	subs := make([]*models.IndicatorSubscription, 0, len(subIDs))
	for _, subID := range subIDs {
		sub, err := r.GetSubscription(subID)
		if err != nil {
			continue
		}
		subs = append(subs, sub)
	}
	return subs, nil
}

// UpdateSubscription 更新訂閱
func (r *RedisRepository) UpdateSubscription(sub *models.IndicatorSubscription) error {
	data, err := json.Marshal(sub)
	if err != nil {
		return err
	}
	key := "indicator_sub:" + sub.SubscriptionID
	return r.client.Set(r.ctx, key, data, 0).Err()
}

// DeleteSubscription 刪除訂閱
func (r *RedisRepository) DeleteSubscription(subscriptionID string) error {
	// 先獲取訂閱資料以取得 UserID 和 Symbol
	sub, err := r.GetSubscription(subscriptionID)
	if err != nil {
		return err
	}

	// 從用戶集合中移除
	userKey := "indicator_subs:user:" + sub.UserID
	r.client.SRem(r.ctx, userKey, subscriptionID)

	// 從幣種集合中移除
	symbolKey := "indicator_subs:symbol:" + sub.Symbol
	r.client.SRem(r.ctx, symbolKey, subscriptionID)

	// 刪除訂閱資料
	key := "indicator_sub:" + subscriptionID
	return r.client.Del(r.ctx, key).Err()
}

// GetAllSubscriptions 獲取所有訂閱
func (r *RedisRepository) GetAllSubscriptions() ([]*models.IndicatorSubscription, error) {
	pattern := "indicator_sub:*"
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}

	subs := make([]*models.IndicatorSubscription, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(r.ctx, key).Result()
		if err != nil {
			continue
		}
		var sub models.IndicatorSubscription
		if err := json.Unmarshal([]byte(data), &sub); err != nil {
			continue
		}
		subs = append(subs, &sub)
	}
	return subs, nil
}

// ==================== 通知冷卻相關方法 ====================

// GetLastNotifyTime 獲取最後通知時間
func (r *RedisRepository) GetLastNotifyTime(key string) (time.Time, error) {
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, data)
}

// SetLastNotifyTime 設置最後通知時間
func (r *RedisRepository) SetLastNotifyTime(key string) error {
	return r.client.Set(r.ctx, key, time.Now().Format(time.RFC3339), 24*time.Hour).Err()
}

// ==================== 指標結果快取相關方法 ====================

// SetIndicatorResult 快取指標計算結果
func (r *RedisRepository) SetIndicatorResult(result *models.IndicatorResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}
	key := "indicator_result:" + result.Symbol
	return r.client.Set(r.ctx, key, data, 30*time.Second).Err() // 快取 30 秒
}

// GetIndicatorResult 獲取快取的指標結果
func (r *RedisRepository) GetIndicatorResult(symbol string) (*models.IndicatorResult, error) {
	key := "indicator_result:" + symbol
	data, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var result models.IndicatorResult
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// ==================== 系統配置相關方法 ====================

// SetIndicatorConfig 儲存指標配置
func (r *RedisRepository) SetIndicatorConfig(config *models.IndicatorConfig) error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, "indicator_config", data, 0).Err()
}

// GetIndicatorConfig 獲取指標配置
func (r *RedisRepository) GetIndicatorConfig() (*models.IndicatorConfig, error) {
	data, err := r.client.Get(r.ctx, "indicator_config").Result()
	if err != nil {
		// 如果不存在，返回預設配置
		if err == redis.Nil {
			defaultConfig := models.DefaultIndicatorConfig()
			return &defaultConfig, nil
		}
		return nil, err
	}
	var config models.IndicatorConfig
	if err := json.Unmarshal([]byte(data), &config); err != nil {
		return nil, err
	}
	return &config, nil
}
