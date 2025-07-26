package services

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type AppSettingService struct {
	appSettingRepo *repositories.AppSettingRepository
}

func NewAppSettingService(appSettingRepo *repositories.AppSettingRepository) *AppSettingService {
	return &AppSettingService{
		appSettingRepo: appSettingRepo,
	}
}

func (s *AppSettingService) GetSettings(userID uuid.UUID) ([]models.AppSetting, error) {
	return s.appSettingRepo.GetByUserID(userID)
}

func (s *AppSettingService) UpdateSetting(userID uuid.UUID, key, value string) error {
	setting := &models.AppSetting{
		ID:        uuid.New(),
		UserID:    userID,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.appSettingRepo.Upsert(setting)
}
