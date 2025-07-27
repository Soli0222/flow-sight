package services

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type IncomeService struct {
	incomeSourceRepo  *repositories.IncomeSourceRepository
	monthlyIncomeRepo *repositories.MonthlyIncomeRepository
}

func NewIncomeService(incomeSourceRepo *repositories.IncomeSourceRepository, monthlyIncomeRepo *repositories.MonthlyIncomeRepository) *IncomeService {
	return &IncomeService{
		incomeSourceRepo:  incomeSourceRepo,
		monthlyIncomeRepo: monthlyIncomeRepo,
	}
}

// Income Source methods
func (s *IncomeService) GetIncomeSources(userID uuid.UUID) ([]models.IncomeSource, error) {
	return s.incomeSourceRepo.GetAll(userID)
}

func (s *IncomeService) GetIncomeSource(id uuid.UUID) (*models.IncomeSource, error) {
	return s.incomeSourceRepo.GetByID(id)
}

func (s *IncomeService) CreateIncomeSource(source *models.IncomeSource) error {
	source.ID = uuid.New()
	source.CreatedAt = time.Now()
	source.UpdatedAt = time.Now()

	return s.incomeSourceRepo.Create(source)
}

func (s *IncomeService) UpdateIncomeSource(source *models.IncomeSource) error {
	source.UpdatedAt = time.Now()
	return s.incomeSourceRepo.Update(source)
}

func (s *IncomeService) DeleteIncomeSource(id uuid.UUID) error {
	return s.incomeSourceRepo.Delete(id)
}

// Monthly Income Record methods
func (s *IncomeService) GetMonthlyIncomeRecords(incomeSourceID uuid.UUID) ([]models.MonthlyIncomeRecord, error) {
	return s.monthlyIncomeRepo.GetByIncomeSourceID(incomeSourceID)
}

func (s *IncomeService) GetMonthlyIncomeRecord(id uuid.UUID) (*models.MonthlyIncomeRecord, error) {
	return s.monthlyIncomeRepo.GetByID(id)
}

func (s *IncomeService) CreateMonthlyIncomeRecord(record *models.MonthlyIncomeRecord) error {
	record.ID = uuid.New()
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()

	return s.monthlyIncomeRepo.Create(record)
}

func (s *IncomeService) UpdateMonthlyIncomeRecord(record *models.MonthlyIncomeRecord) error {
	record.UpdatedAt = time.Now()
	return s.monthlyIncomeRepo.Update(record)
}

func (s *IncomeService) DeleteMonthlyIncomeRecord(id uuid.UUID) error {
	return s.monthlyIncomeRepo.Delete(id)
}
