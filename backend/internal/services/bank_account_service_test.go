package services

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/services/mocks"
	"flow-sight-backend/test/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBankAccountService_GetBankAccounts(t *testing.T) {
	mockRepo := &mocks.MockBankAccountRepository{}
	service := NewBankAccountService(mockRepo)
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(*mocks.MockBankAccountRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful retrieval",
			userID: userID,
			setupMock: func(m *mocks.MockBankAccountRepository) {
				accounts := []models.BankAccount{
					*helpers.CreateTestBankAccount(userID),
					*helpers.CreateTestBankAccount(userID),
				}
				m.On("GetAll", userID).Return(accounts, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "empty result",
			userID: userID,
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("GetAll", userID).Return([]models.BankAccount{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "repository error",
			userID: userID,
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("GetAll", userID).Return([]models.BankAccount{}, assert.AnError)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock for each test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.setupMock(mockRepo)

			accounts, err := service.GetBankAccounts(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, accounts, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBankAccountService_GetBankAccount(t *testing.T) {
	mockRepo := &mocks.MockBankAccountRepository{}
	service := NewBankAccountService(mockRepo)
	accountID := uuid.New()
	userID := uuid.New()

	tests := []struct {
		name          string
		accountID     uuid.UUID
		setupMock     func(*mocks.MockBankAccountRepository)
		expectedFound bool
		expectedError bool
	}{
		{
			name:      "account found",
			accountID: accountID,
			setupMock: func(m *mocks.MockBankAccountRepository) {
				account := helpers.CreateTestBankAccount(userID)
				account.ID = accountID
				m.On("GetByID", accountID).Return(account, nil)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:      "account not found",
			accountID: accountID,
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("GetByID", accountID).Return(nil, assert.AnError)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock for each test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.setupMock(mockRepo)

			account, err := service.GetBankAccount(tt.accountID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, account)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.accountID, account.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBankAccountService_CreateBankAccount(t *testing.T) {
	mockRepo := &mocks.MockBankAccountRepository{}
	service := NewBankAccountService(mockRepo)
	userID := uuid.New()

	tests := []struct {
		name          string
		account       *models.BankAccount
		setupMock     func(*mocks.MockBankAccountRepository)
		expectedError bool
	}{
		{
			name:    "successful creation",
			account: helpers.CreateTestBankAccount(userID),
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("Create", mock.AnythingOfType("*models.BankAccount")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "repository error",
			account: helpers.CreateTestBankAccount(userID),
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("Create", mock.AnythingOfType("*models.BankAccount")).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock for each test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.setupMock(mockRepo)

			err := service.CreateBankAccount(tt.account)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify that UUID was generated
				assert.NotEqual(t, uuid.Nil, tt.account.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBankAccountService_UpdateBankAccount(t *testing.T) {
	mockRepo := &mocks.MockBankAccountRepository{}
	service := NewBankAccountService(mockRepo)
	userID := uuid.New()

	tests := []struct {
		name          string
		account       *models.BankAccount
		setupMock     func(*mocks.MockBankAccountRepository)
		expectedError bool
	}{
		{
			name:    "successful update",
			account: helpers.CreateTestBankAccount(userID),
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("Update", mock.AnythingOfType("*models.BankAccount")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "repository error",
			account: helpers.CreateTestBankAccount(userID),
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("Update", mock.AnythingOfType("*models.BankAccount")).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock for each test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.setupMock(mockRepo)

			err := service.UpdateBankAccount(tt.account)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBankAccountService_DeleteBankAccount(t *testing.T) {
	mockRepo := &mocks.MockBankAccountRepository{}
	service := NewBankAccountService(mockRepo)
	accountID := uuid.New()

	tests := []struct {
		name          string
		accountID     uuid.UUID
		setupMock     func(*mocks.MockBankAccountRepository)
		expectedError bool
	}{
		{
			name:      "successful deletion",
			accountID: accountID,
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("Delete", accountID).Return(nil)
			},
			expectedError: false,
		},
		{
			name:      "repository error",
			accountID: accountID,
			setupMock: func(m *mocks.MockBankAccountRepository) {
				m.On("Delete", accountID).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock for each test
			mockRepo.ExpectedCalls = nil
			mockRepo.Calls = nil

			tt.setupMock(mockRepo)

			err := service.DeleteBankAccount(tt.accountID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
