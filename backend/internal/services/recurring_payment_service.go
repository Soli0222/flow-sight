package services

import (
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

type RecurringPaymentService struct {
	recurringPaymentRepo RecurringPaymentRepositoryInterface
}

func NewRecurringPaymentService(recurringPaymentRepo RecurringPaymentRepositoryInterface) *RecurringPaymentService {
	return &RecurringPaymentService{
		recurringPaymentRepo: recurringPaymentRepo,
	}
}

func (s *RecurringPaymentService) GetRecurringPayments() ([]models.RecurringPayment, error) {
	return s.recurringPaymentRepo.GetAll()
}

func (s *RecurringPaymentService) GetRecurringPayment(id uuid.UUID) (*models.RecurringPayment, error) {
	return s.recurringPaymentRepo.GetByID(id)
}

func (s *RecurringPaymentService) CreateRecurringPayment(payment *models.RecurringPayment) error {
	payment.ID = uuid.New()
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()

	return s.recurringPaymentRepo.Create(payment)
}

func (s *RecurringPaymentService) UpdateRecurringPayment(payment *models.RecurringPayment) error {
	payment.UpdatedAt = time.Now()
	return s.recurringPaymentRepo.Update(payment)
}

func (s *RecurringPaymentService) DeleteRecurringPayment(id uuid.UUID) error {
	return s.recurringPaymentRepo.Delete(id)
}
