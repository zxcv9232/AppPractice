package repository

import (
	"context"
	"encoding/json"
	"fmt"
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

func (r *RedisRepository) SaveVolumeSnapshot(symbol string, volume float64) error {
	snapshot := &models.VolumeSnapshot{
		Symbol:    symbol,
		Volume:    volume,
		Timestamp: time.Now(),
	}
	
	data, err := json.Marshal(snapshot)
	if err != nil {
		return err
	}
	
	key := "volume:" + symbol
	return r.client.ZAdd(r.ctx, key, redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: data,
	}).Err()
}

func (r *RedisRepository) GetCumulativeVolume(symbol string, minutes int) (float64, error) {
	key := "volume:" + symbol
	now := time.Now()
	startTime := now.Add(-time.Duration(minutes) * time.Minute)
	
	results, err := r.client.ZRangeByScore(r.ctx, key, &redis.ZRangeBy{
		Min: fmt.Sprintf("%d", startTime.Unix()),
		Max: fmt.Sprintf("%d", now.Unix()),
	}).Result()
	
	if err != nil {
		return 0, err
	}
	
	if len(results) == 0 {
		return 0, nil
	}
	
	var firstSnapshot, lastSnapshot models.VolumeSnapshot
	
	if err := json.Unmarshal([]byte(results[0]), &firstSnapshot); err != nil {
		return 0, err
	}
	
	if err := json.Unmarshal([]byte(results[len(results)-1]), &lastSnapshot); err != nil {
		return 0, err
	}
	
	return lastSnapshot.Volume - firstSnapshot.Volume, nil
}

func (r *RedisRepository) CleanOldVolumeSnapshots(symbol string, hoursToKeep int) error {
	key := "volume:" + symbol
	cutoffTime := time.Now().Add(-time.Duration(hoursToKeep) * time.Hour)
	
	return r.client.ZRemRangeByScore(r.ctx, key, "0", fmt.Sprintf("%d", cutoffTime.Unix())).Err()
}

