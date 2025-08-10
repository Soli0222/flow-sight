package mocks

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockBankAccountService は BankAccountServiceInterface のモック
type MockBankAccountService struct {
	mock.Mock
}

func (m *MockBankAccountService) GetBankAccounts() ([]models.BankAccount, error) {
	args := m.Called()
	return args.Get(0).([]models.BankAccount), args.Error(1)
}

func (m *MockBankAccountService) GetBankAccount(id uuid.UUID) (*models.BankAccount, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountService) CreateBankAccount(account *models.BankAccount) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockBankAccountService) UpdateBankAccount(account *models.BankAccount) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockBankAccountService) DeleteBankAccount(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
