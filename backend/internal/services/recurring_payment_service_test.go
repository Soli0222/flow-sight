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

func TestRecurringPaymentService_GetRecurringPayments(t *testing.T) {
	mockRepo := &mocks.MockRecurringPaymentRepository{}
	service := NewRecurringPaymentService(mockRepo)
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(*mocks.MockRecurringPaymentRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful retrieval",
			userID: userID,
			setupMock: func(m *mocks.MockRecurringPaymentRepository) {
				payments := []models.RecurringPayment{
					*helpers.CreateTestRecurringPayment(userID, bankAccountID),
					*helpers.CreateTestRecurringPayment(userID, bankAccountID),
				}
				m.On("GetAll", userID).Return(payments, nil)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "empty result",
			userID: userID,
			setupMock: func(m *mocks.MockRecurringPaymentRepository) {
				m.On("GetAll", userID).Return([]models.RecurringPayment{}, nil)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "repository error",
			userID: userID,
			setupMock: func(m *mocks.MockRecurringPaymentRepository) {
				m.On("GetAll", userID).Return([]models.RecurringPayment{}, assert.AnError)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock(mockRepo)

			result, err := service.GetRecurringPayments(tt.userID)

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

func TestRecurringPaymentService_CreateRecurringPayment(t *testing.T) {
	mockRepo := &mocks.MockRecurringPaymentRepository{}
	service := NewRecurringPaymentService(mockRepo)
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		payment       *models.RecurringPayment
		setupMock     func(*mocks.MockRecurringPaymentRepository, *models.RecurringPayment)
		expectedError bool
	}{
		{
			name:    "successful creation",
			payment: helpers.CreateTestRecurringPayment(userID, bankAccountID),
			setupMock: func(m *mocks.MockRecurringPaymentRepository, rp *models.RecurringPayment) {
				m.On("Create", mock.MatchedBy(func(payment *models.RecurringPayment) bool {
					return payment.UserID == userID && payment.Name == rp.Name
				})).Return(nil)
			},
			expectedError: false,
		},
		{
			name:    "repository error",
			payment: helpers.CreateTestRecurringPayment(userID, bankAccountID),
			setupMock: func(m *mocks.MockRecurringPaymentRepository, rp *models.RecurringPayment) {
				m.On("Create", mock.AnythingOfType("*models.RecurringPayment")).Return(assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			originalID := tt.payment.ID
			tt.setupMock(mockRepo, tt.payment)

			err := service.CreateRecurringPayment(tt.payment)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// IDが新しく設定されていることを確認
				assert.NotEqual(t, originalID, tt.payment.ID)
				assert.NotEqual(t, uuid.Nil, tt.payment.ID)
				// 作成日時が設定されていることを確認
				assert.False(t, tt.payment.CreatedAt.IsZero())
				assert.False(t, tt.payment.UpdatedAt.IsZero())
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
