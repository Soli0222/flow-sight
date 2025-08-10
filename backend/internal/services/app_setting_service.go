package services

import (
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/repositories"

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

func (s *AppSettingService) GetSettings() ([]models.AppSetting, error) {
	return s.appSettingRepo.GetAll()
}

func (s *AppSettingService) UpdateSetting(key, value string) error {
	setting := &models.AppSetting{
		ID:        uuid.New(),
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.appSettingRepo.Upsert(setting)
}
