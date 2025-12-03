package service

import (
	"time"

	"cryptowatch/internal/models"
	"cryptowatch/internal/repository"
	"github.com/google/uuid"
)

type AlertService struct {
	repo *repository.RedisRepository
}

func NewAlertService(repo *repository.RedisRepository) *AlertService {
	return &AlertService{repo: repo}
}

func (s *AlertService) CreateAlert(req *models.CreateAlertRequest) (*models.Alert, error) {
	alert := &models.Alert{
		AlertID:      uuid.New().String(),
		UserID:       req.UserID,
		Symbol:       req.Symbol,
		AlertType:    req.AlertType,
		TargetPrice:  req.TargetPrice,
		Direction:    req.Direction,
		TargetVolume: req.TargetVolume,
		TimeWindow:   req.TimeWindow,
		CreatedAt:    time.Now(),
	}
	if err := s.repo.SaveAlert(alert); err != nil {
		return nil, err
	}
	return alert, nil
}

func (s *AlertService) GetUserAlerts(userID string) ([]*models.Alert, error) {
	return s.repo.GetUserAlerts(userID)
}

func (s *AlertService) DeleteAlert(alertID string) error {
	return s.repo.DeleteAlert(alertID)
}

