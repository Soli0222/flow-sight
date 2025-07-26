package services

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type AssetService struct {
	assetRepo *repositories.AssetRepository
}

func NewAssetService(assetRepo *repositories.AssetRepository) *AssetService {
	return &AssetService{
		assetRepo: assetRepo,
	}
}

func (s *AssetService) GetAssets(userID uuid.UUID) ([]models.Asset, error) {
	return s.assetRepo.GetAll(userID)
}

func (s *AssetService) GetAsset(id uuid.UUID) (*models.Asset, error) {
	return s.assetRepo.GetByID(id)
}

func (s *AssetService) CreateAsset(asset *models.Asset) error {
	asset.ID = uuid.New()
	asset.CreatedAt = time.Now()
	asset.UpdatedAt = time.Now()

	return s.assetRepo.Create(asset)
}

func (s *AssetService) UpdateAsset(asset *models.Asset) error {
	asset.UpdatedAt = time.Now()
	return s.assetRepo.Update(asset)
}

func (s *AssetService) DeleteAsset(id uuid.UUID) error {
	return s.assetRepo.Delete(id)
}
