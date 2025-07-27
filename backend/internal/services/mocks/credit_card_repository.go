package mocks

import (
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockCreditCardRepository は CreditCardRepositoryInterface のモック
type MockCreditCardRepository struct {
	mock.Mock
}

func (m *MockCreditCardRepository) GetAll(userID uuid.UUID) ([]models.CreditCard, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) GetByID(id uuid.UUID) (*models.CreditCard, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) Create(creditCard *models.CreditCard) error {
	args := m.Called(creditCard)
	return args.Error(0)
}

func (m *MockCreditCardRepository) Update(creditCard *models.CreditCard) error {
	args := m.Called(creditCard)
	return args.Error(0)
}

func (m *MockCreditCardRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
