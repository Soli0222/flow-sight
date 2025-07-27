package mocks

import (
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockMonthlyIncomeRepository は MonthlyIncomeRepositoryInterface のモック
type MockMonthlyIncomeRepository struct {
	mock.Mock
}

func (m *MockMonthlyIncomeRepository) GetByIncomeSourceID(incomeSourceID uuid.UUID) ([]models.MonthlyIncomeRecord, error) {
	args := m.Called(incomeSourceID)
	return args.Get(0).([]models.MonthlyIncomeRecord), args.Error(1)
}

func (m *MockMonthlyIncomeRepository) GetByYearMonth(userID uuid.UUID, yearMonth string) ([]models.MonthlyIncomeRecord, error) {
	args := m.Called(userID, yearMonth)
	return args.Get(0).([]models.MonthlyIncomeRecord), args.Error(1)
}

// Add convenience method used by services
func (m *MockMonthlyIncomeRepository) GetByUserIDAndYearMonth(userID uuid.UUID, yearMonth string) ([]models.MonthlyIncomeRecord, error) {
	args := m.Called(userID, yearMonth)
	return args.Get(0).([]models.MonthlyIncomeRecord), args.Error(1)
}
