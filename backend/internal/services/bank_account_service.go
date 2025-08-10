package services

import (
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

type BankAccountService struct {
	bankAccountRepo BankAccountRepositoryInterface
}

func NewBankAccountService(bankAccountRepo BankAccountRepositoryInterface) *BankAccountService {
	return &BankAccountService{
		bankAccountRepo: bankAccountRepo,
	}
}

func (s *BankAccountService) GetBankAccounts() ([]models.BankAccount, error) {
	return s.bankAccountRepo.GetAll()
}

func (s *BankAccountService) GetBankAccount(id uuid.UUID) (*models.BankAccount, error) {
	return s.bankAccountRepo.GetByID(id)
}

func (s *BankAccountService) CreateBankAccount(account *models.BankAccount) error {
	account.ID = uuid.New()
	account.CreatedAt = time.Now()
	account.UpdatedAt = time.Now()

	return s.bankAccountRepo.Create(account)
}

func (s *BankAccountService) UpdateBankAccount(account *models.BankAccount) error {
	account.UpdatedAt = time.Now()
	return s.bankAccountRepo.Update(account)
}

func (s *BankAccountService) DeleteBankAccount(id uuid.UUID) error {
	return s.bankAccountRepo.Delete(id)
}
