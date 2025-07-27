package services

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type CreditCardService struct {
	creditCardRepo CreditCardRepositoryInterface
}

func NewCreditCardService(creditCardRepo CreditCardRepositoryInterface) *CreditCardService {
	return &CreditCardService{
		creditCardRepo: creditCardRepo,
	}
}

func (s *CreditCardService) GetCreditCards(userID uuid.UUID) ([]models.CreditCard, error) {
	return s.creditCardRepo.GetAll(userID)
}

func (s *CreditCardService) GetCreditCard(id uuid.UUID) (*models.CreditCard, error) {
	return s.creditCardRepo.GetByID(id)
}

func (s *CreditCardService) CreateCreditCard(creditCard *models.CreditCard) error {
	creditCard.ID = uuid.New()
	creditCard.CreatedAt = time.Now()
	creditCard.UpdatedAt = time.Now()

	return s.creditCardRepo.Create(creditCard)
}

func (s *CreditCardService) UpdateCreditCard(creditCard *models.CreditCard) error {
	creditCard.UpdatedAt = time.Now()
	return s.creditCardRepo.Update(creditCard)
}

func (s *CreditCardService) DeleteCreditCard(id uuid.UUID) error {
	return s.creditCardRepo.Delete(id)
}
