package services

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type CardMonthlyTotalService struct {
	cardMonthlyTotalRepo *repositories.CardMonthlyTotalRepository
}

func NewCardMonthlyTotalService(cardMonthlyTotalRepo *repositories.CardMonthlyTotalRepository) *CardMonthlyTotalService {
	return &CardMonthlyTotalService{
		cardMonthlyTotalRepo: cardMonthlyTotalRepo,
	}
}

func (s *CardMonthlyTotalService) GetCardMonthlyTotals(creditCardID uuid.UUID) ([]models.CardMonthlyTotal, error) {
	return s.cardMonthlyTotalRepo.GetByCreditCardID(creditCardID)
}

func (s *CardMonthlyTotalService) GetCardMonthlyTotal(id uuid.UUID) (*models.CardMonthlyTotal, error) {
	return s.cardMonthlyTotalRepo.GetByID(id)
}

func (s *CardMonthlyTotalService) CreateCardMonthlyTotal(total *models.CardMonthlyTotal) error {
	total.ID = uuid.New()
	total.CreatedAt = time.Now()
	total.UpdatedAt = time.Now()

	return s.cardMonthlyTotalRepo.Create(total)
}

func (s *CardMonthlyTotalService) UpdateCardMonthlyTotal(total *models.CardMonthlyTotal) error {
	total.UpdatedAt = time.Now()
	return s.cardMonthlyTotalRepo.Update(total)
}

func (s *CardMonthlyTotalService) DeleteCardMonthlyTotal(id uuid.UUID) error {
	return s.cardMonthlyTotalRepo.Delete(id)
}
