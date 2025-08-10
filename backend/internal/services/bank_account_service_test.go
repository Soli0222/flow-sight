package services

import (
	"testing"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/services/mocks"
	"github.com/Soli0222/flow-sight/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBankAccountService_GetBankAccounts(t *testing.T) {
	mockRepo := new(mocks.MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	accounts := []models.BankAccount{*helpers.CreateTestBankAccount(), *helpers.CreateTestBankAccount()}
	mockRepo.On("GetAll").Return(accounts, nil)

	res, err := service.GetBankAccounts()
	assert.NoError(t, err)
	assert.Len(t, res, 2)
}

func TestBankAccountService_GetBankAccount(t *testing.T) {
	mockRepo := new(mocks.MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	id := uuid.New()
	acc := helpers.CreateTestBankAccount()
	mockRepo.On("GetByID", id).Return(acc, nil)

	res, err := service.GetBankAccount(id)
	assert.NoError(t, err)
	assert.Equal(t, acc.ID, res.ID)
}

func TestBankAccountService_CreateUpdateDelete(t *testing.T) {
	mockRepo := new(mocks.MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	acc := helpers.CreateTestBankAccount()

	// Create: service will mutate acc.ID and timestamps
	mockRepo.On("Create", mock.AnythingOfType("*models.BankAccount")).Return(nil).Once()

	// Update: expect the same struct passed back
	mockRepo.On("Update", mock.AnythingOfType("*models.BankAccount")).Return(nil).Once()

	// create
	err := service.CreateBankAccount(acc)
	assert.NoError(t, err)

	// update
	err = service.UpdateBankAccount(acc)
	assert.NoError(t, err)

	// delete: set expectation after create to use the finalized ID
	finalID := acc.ID
	mockRepo.On("Delete", finalID).Return(nil).Once()

	err = service.DeleteBankAccount(finalID)
	assert.NoError(t, err)
}
