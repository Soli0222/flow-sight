package handlers

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RecurringPaymentHandler struct {
	recurringPaymentService *services.RecurringPaymentService
}

func NewRecurringPaymentHandler(recurringPaymentService *services.RecurringPaymentService) *RecurringPaymentHandler {
	return &RecurringPaymentHandler{
		recurringPaymentService: recurringPaymentService,
	}
}

// @Summary Get all recurring payments
// @Description Get all recurring payments for a user
// @Tags recurring-payments
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} models.RecurringPayment
// @Router /recurring-payments [get]
func (h *RecurringPaymentHandler) GetRecurringPayments(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	payments, err := h.recurringPaymentService.GetRecurringPayments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payments)
}

// @Summary Get recurring payment by ID
// @Description Get a specific recurring payment by ID
// @Tags recurring-payments
// @Accept json
// @Produce json
// @Param id path string true "Recurring Payment ID"
// @Success 200 {object} models.RecurringPayment
// @Router /recurring-payments/{id} [get]
func (h *RecurringPaymentHandler) GetRecurringPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recurring payment id format"})
		return
	}

	payment, err := h.recurringPaymentService.GetRecurringPayment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "recurring payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Create recurring payment
// @Description Create a new recurring payment
// @Tags recurring-payments
// @Accept json
// @Produce json
// @Param payment body models.RecurringPayment true "Recurring Payment data"
// @Success 201 {object} models.RecurringPayment
// @Router /recurring-payments [post]
func (h *RecurringPaymentHandler) CreateRecurringPayment(c *gin.Context) {
	var payment models.RecurringPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.recurringPaymentService.CreateRecurringPayment(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// @Summary Update recurring payment
// @Description Update an existing recurring payment
// @Tags recurring-payments
// @Accept json
// @Produce json
// @Param id path string true "Recurring Payment ID"
// @Param payment body models.RecurringPayment true "Recurring Payment data"
// @Success 200 {object} models.RecurringPayment
// @Router /recurring-payments/{id} [put]
func (h *RecurringPaymentHandler) UpdateRecurringPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recurring payment id format"})
		return
	}

	var payment models.RecurringPayment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment.ID = id
	if err := h.recurringPaymentService.UpdateRecurringPayment(&payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// @Summary Delete recurring payment
// @Description Delete a recurring payment
// @Tags recurring-payments
// @Accept json
// @Produce json
// @Param id path string true "Recurring Payment ID"
// @Success 204
// @Router /recurring-payments/{id} [delete]
func (h *RecurringPaymentHandler) DeleteRecurringPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid recurring payment id format"})
		return
	}

	if err := h.recurringPaymentService.DeleteRecurringPayment(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
