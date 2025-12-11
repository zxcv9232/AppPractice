package service

import (
	"time"

	"cryptowatch/internal/models"
	"cryptowatch/internal/repository"

	"github.com/google/uuid"
)

// SubscriptionService 訂閱服務
type SubscriptionService struct {
	repo *repository.RedisRepository
}

// NewSubscriptionService 創建訂閱服務
func NewSubscriptionService(repo *repository.RedisRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

// CreateSubscription 創建訂閱
func (s *SubscriptionService) CreateSubscription(req *models.CreateSubscriptionRequest) (*models.IndicatorSubscription, error) {
	// 套用預設值
	req.ApplyDefaults()

	sub := &models.IndicatorSubscription{
		SubscriptionID:    uuid.New().String(),
		UserID:            req.UserID,
		Symbol:            req.Symbol,
		Enabled:           true, // 創建時預設啟用
		TelegramChatID:    req.TelegramChatID,
		NotifyIntervalMin: req.NotifyIntervalMin,
		EnableVolumeCheck: req.EnableVolumeCheck,
		VolumeCheckMode:   req.VolumeCheckMode,
		VolumeFixedValue:  req.VolumeFixedValue,
		VolumeMultiplier:  req.VolumeMultiplier,
		VolumeAvgPeriod:   req.VolumeAvgPeriod,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.repo.SaveSubscription(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// GetSubscription 獲取訂閱
func (s *SubscriptionService) GetSubscription(subscriptionID string) (*models.IndicatorSubscription, error) {
	return s.repo.GetSubscription(subscriptionID)
}

// GetUserSubscriptions 獲取用戶的所有訂閱
func (s *SubscriptionService) GetUserSubscriptions(userID string) ([]*models.IndicatorSubscription, error) {
	return s.repo.GetUserSubscriptions(userID)
}

// UpdateSubscription 更新訂閱
func (s *SubscriptionService) UpdateSubscription(subscriptionID string, req *models.UpdateSubscriptionRequest) (*models.IndicatorSubscription, error) {
	sub, err := s.repo.GetSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// 更新欄位（只更新有提供的欄位）
	if req.Enabled != nil {
		sub.Enabled = *req.Enabled
	}
	if req.TelegramChatID != nil {
		sub.TelegramChatID = *req.TelegramChatID
	}
	if req.NotifyIntervalMin != nil {
		sub.NotifyIntervalMin = *req.NotifyIntervalMin
	}
	if req.EnableVolumeCheck != nil {
		sub.EnableVolumeCheck = *req.EnableVolumeCheck
	}
	if req.VolumeCheckMode != nil {
		sub.VolumeCheckMode = *req.VolumeCheckMode
	}
	if req.VolumeFixedValue != nil {
		sub.VolumeFixedValue = *req.VolumeFixedValue
	}
	if req.VolumeMultiplier != nil {
		sub.VolumeMultiplier = *req.VolumeMultiplier
	}
	if req.VolumeAvgPeriod != nil {
		sub.VolumeAvgPeriod = *req.VolumeAvgPeriod
	}

	sub.UpdatedAt = time.Now()

	if err := s.repo.UpdateSubscription(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// DeleteSubscription 刪除訂閱
func (s *SubscriptionService) DeleteSubscription(subscriptionID string) error {
	return s.repo.DeleteSubscription(subscriptionID)
}

// ToggleSubscription 切換訂閱開關
func (s *SubscriptionService) ToggleSubscription(subscriptionID string) (*models.IndicatorSubscription, error) {
	sub, err := s.repo.GetSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	sub.Enabled = !sub.Enabled
	sub.UpdatedAt = time.Now()

	if err := s.repo.UpdateSubscription(sub); err != nil {
		return nil, err
	}

	return sub, nil
}

