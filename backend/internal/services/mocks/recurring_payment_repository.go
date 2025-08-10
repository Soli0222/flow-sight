package mocks

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockRecurringPaymentRepository は RecurringPaymentRepositoryInterface のモック
type MockRecurringPaymentRepository struct {
	mock.Mock
}

func (m *MockRecurringPaymentRepository) GetAll() ([]models.RecurringPayment, error) {
	args := m.Called()
	return args.Get(0).([]models.RecurringPayment), args.Error(1)
}

func (m *MockRecurringPaymentRepository) GetActive() ([]models.RecurringPayment, error) {
	args := m.Called()
	return args.Get(0).([]models.RecurringPayment), args.Error(1)
}

func (m *MockRecurringPaymentRepository) GetByID(id uuid.UUID) (*models.RecurringPayment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.RecurringPayment), args.Error(1)
}

func (m *MockRecurringPaymentRepository) Create(payment *models.RecurringPayment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockRecurringPaymentRepository) Update(payment *models.RecurringPayment) error {
	args := m.Called(payment)
	return args.Error(0)
}

func (m *MockRecurringPaymentRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
