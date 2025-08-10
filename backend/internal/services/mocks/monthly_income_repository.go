package mocks

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"

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

func (m *MockMonthlyIncomeRepository) GetByYearMonth(yearMonth string) ([]models.MonthlyIncomeRecord, error) {
	args := m.Called(yearMonth)
	return args.Get(0).([]models.MonthlyIncomeRecord), args.Error(1)
}

// Deprecated helper retained for backward compatibility in tests (no-op userID)
func (m *MockMonthlyIncomeRepository) GetByUserIDAndYearMonth(_ uuid.UUID, yearMonth string) ([]models.MonthlyIncomeRecord, error) {
	args := m.Called(yearMonth)
	return args.Get(0).([]models.MonthlyIncomeRecord), args.Error(1)
}
