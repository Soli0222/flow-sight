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

func TestCreditCardService_GetCreditCards(t *testing.T) {
	mockRepo := &mocks.MockCreditCardRepository{}
	service := NewCreditCardService(mockRepo)
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(*mocks.MockCreditCardRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful retrieval",
			userID: userID,
			setupMock: func(m *mocks.MockCreditCardRepository) {
				creditCards := []models.CreditCard{
					*helpers.CreateTestCreditCard(userID, bankAccountID),
					*helpers.CreateTestCreditCard(userID, bankAccountID),
				}
				m.On("GetAll", userID).Return(creditCards, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "empty result",
			userID: userID,
			setupMock: func(m *mocks.MockCreditCardRepository) {
				m.On("GetAll", userID).Return([]models.CreditCard{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "repository error",
			userID: userID,
			setupMock: func(m *mocks.MockCreditCardRepository) {
				m.On("GetAll", userID).Return([]models.CreditCard{}, assert.AnError)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock(mockRepo)

			result, err := service.GetCreditCards(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreditCardService_GetCreditCard(t *testing.T) {
	mockRepo := &mocks.MockCreditCardRepository{}
	service := NewCreditCardService(mockRepo)
	creditCardID := uuid.New()
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		creditCardID  uuid.UUID
		setupMock     func(*mocks.MockCreditCardRepository)
		expectedNil   bool
		expectedError bool
	}{
		{
			name:         "successful retrieval",
			creditCardID: creditCardID,
			setupMock: func(m *mocks.MockCreditCardRepository) {
				creditCard := helpers.CreateTestCreditCard(userID, bankAccountID)
				creditCard.ID = creditCardID
				m.On("GetByID", creditCardID).Return(creditCard, nil)
			},
			expectedNil:   false,
			expectedError: false,
		},
		{
			name:         "credit card not found",
			creditCardID: creditCardID,
			setupMock: func(m *mocks.MockCreditCardRepository) {
				m.On("GetByID", creditCardID).Return(nil, assert.AnError)
			},
			expectedNil:   true,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock(mockRepo)

			result, err := service.GetCreditCard(tt.creditCardID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.creditCardID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreditCardService_CreateCreditCard(t *testing.T) {
	mockRepo := &mocks.MockCreditCardRepository{}
	service := NewCreditCardService(mockRepo)
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		creditCard    *models.CreditCard
		setupMock     func(*mocks.MockCreditCardRepository, *models.CreditCard)
		expectedError bool
	}{
		{
			name:       "successful creation",
			creditCard: helpers.CreateTestCreditCard(userID, bankAccountID),
			setupMock: func(m *mocks.MockCreditCardRepository, cc *models.CreditCard) {
				m.On("Create", mock.MatchedBy(func(card *models.CreditCard) bool {
					return card.UserID == userID && card.Name == cc.Name
				})).Return(nil)
			},
			expectedError: false,
		},
		{
			name:       "repository error",
			creditCard: helpers.CreateTestCreditCard(userID, bankAccountID),
			setupMock: func(m *mocks.MockCreditCardRepository, cc *models.CreditCard) {
				m.On("Create", mock.AnythingOfType("*models.CreditCard")).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			originalID := tt.creditCard.ID
			tt.setupMock(mockRepo, tt.creditCard)

			err := service.CreateCreditCard(tt.creditCard)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// IDが新しく設定されていることを確認
				assert.NotEqual(t, originalID, tt.creditCard.ID)
				assert.NotEqual(t, uuid.Nil, tt.creditCard.ID)
				// 作成日時が設定されていることを確認
				assert.False(t, tt.creditCard.CreatedAt.IsZero())
				assert.False(t, tt.creditCard.UpdatedAt.IsZero())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreditCardService_UpdateCreditCard(t *testing.T) {
	mockRepo := &mocks.MockCreditCardRepository{}
	service := NewCreditCardService(mockRepo)
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		creditCard    *models.CreditCard
		setupMock     func(*mocks.MockCreditCardRepository)
		expectedError bool
	}{
		{
			name:       "successful update",
			creditCard: helpers.CreateTestCreditCard(userID, bankAccountID),
			setupMock: func(m *mocks.MockCreditCardRepository) {
				m.On("Update", mock.AnythingOfType("*models.CreditCard")).Return(nil)
			},
			expectedError: false,
		},
		{
			name:       "repository error",
			creditCard: helpers.CreateTestCreditCard(userID, bankAccountID),
			setupMock: func(m *mocks.MockCreditCardRepository) {
				m.On("Update", mock.AnythingOfType("*models.CreditCard")).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			originalUpdatedAt := tt.creditCard.UpdatedAt
			tt.setupMock(mockRepo)

			err := service.UpdateCreditCard(tt.creditCard)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 更新日時が更新されていることを確認
				assert.True(t, tt.creditCard.UpdatedAt.After(originalUpdatedAt))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreditCardService_DeleteCreditCard(t *testing.T) {
	mockRepo := &mocks.MockCreditCardRepository{}
	service := NewCreditCardService(mockRepo)
	creditCardID := uuid.New()

	tests := []struct {
		name          string
		creditCardID  uuid.UUID
		setupMock     func(*mocks.MockCreditCardRepository)
		expectedError bool
	}{
		{
			name:         "successful deletion",
			creditCardID: creditCardID,
			setupMock: func(m *mocks.MockCreditCardRepository) {
				m.On("Delete", creditCardID).Return(nil)
			},
			expectedError: false,
		},
		{
			name:         "repository error",
			creditCardID: creditCardID,
			setupMock: func(m *mocks.MockCreditCardRepository) {
				m.On("Delete", creditCardID).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock(mockRepo)

			err := service.DeleteCreditCard(tt.creditCardID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
