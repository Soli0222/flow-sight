package mocks

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/stretchr/testify/mock"
)

// MockCashflowService は CashflowService のモック
type MockCashflowService struct {
	mock.Mock
}

func (m *MockCashflowService) GetCashflowProjection(months int, includeConfirmed bool) ([]models.CashflowProjection, error) {
	args := m.Called(months, includeConfirmed)
	return args.Get(0).([]models.CashflowProjection), args.Error(1)
}
