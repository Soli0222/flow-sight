package mocks

import (
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockIncomeSourceRepository は IncomeSourceRepositoryInterface のモック
type MockIncomeSourceRepository struct {
	mock.Mock
}

func (m *MockIncomeSourceRepository) GetAll(userID uuid.UUID) ([]models.IncomeSource, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.IncomeSource), args.Error(1)
}

func (m *MockIncomeSourceRepository) GetByID(id uuid.UUID) (*models.IncomeSource, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.IncomeSource), args.Error(1)
}

func (m *MockIncomeSourceRepository) GetActiveByUserID(userID uuid.UUID) ([]models.IncomeSource, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.IncomeSource), args.Error(1)
}

func (m *MockIncomeSourceRepository) Create(incomeSource *models.IncomeSource) error {
	args := m.Called(incomeSource)
	return args.Error(0)
}

func (m *MockIncomeSourceRepository) Update(incomeSource *models.IncomeSource) error {
	args := m.Called(incomeSource)
	return args.Error(0)
}

func (m *MockIncomeSourceRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
