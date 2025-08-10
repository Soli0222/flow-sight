package mocks

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockBankAccountRepository は BankAccountRepositoryInterface のモック
type MockBankAccountRepository struct {
	mock.Mock
}

func (m *MockBankAccountRepository) GetAll() ([]models.BankAccount, error) {
	args := m.Called()
	return args.Get(0).([]models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) GetByID(id uuid.UUID) (*models.BankAccount, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) Create(account *models.BankAccount) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockBankAccountRepository) Update(account *models.BankAccount) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockBankAccountRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
