package handlers

import (
	"net/http"

	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreditCardHandler struct {
	creditCardService CreditCardServiceInterface
}

func NewCreditCardHandler(creditCardService CreditCardServiceInterface) *CreditCardHandler {
	return &CreditCardHandler{
		creditCardService: creditCardService,
	}
}

// @Summary Get all credit cards
// @Description Get all credit cards
// @Tags credit-cards
// @Accept json
// @Produce json
// @Success 200 {array} models.CreditCard
// @Router /credit-cards [get]
func (h *CreditCardHandler) GetCreditCards(c *gin.Context) {
	creditCards, err := h.creditCardService.GetCreditCards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, creditCards)
}

// @Summary Get credit card by ID
// @Description Get a specific credit card by ID
// @Tags credit-cards
// @Accept json
// @Produce json
// @Param id path string true "Credit Card ID"
// @Success 200 {object} models.CreditCard
// @Router /credit-cards/{id} [get]
func (h *CreditCardHandler) GetCreditCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit card id format"})
		return
	}

	creditCard, err := h.creditCardService.GetCreditCard(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "credit card not found"})
		return
	}

	c.JSON(http.StatusOK, creditCard)
}

// @Summary Create credit card
// @Description Create a new credit card
// @Tags credit-cards
// @Accept json
// @Produce json
// @Param creditCard body models.CreditCard true "Credit Card data"
// @Success 201 {object} models.CreditCard
// @Router /credit-cards [post]
func (h *CreditCardHandler) CreateCreditCard(c *gin.Context) {
	var creditCard models.CreditCard
	if err := c.ShouldBindJSON(&creditCard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.creditCardService.CreateCreditCard(&creditCard); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, creditCard)
}

// @Summary Update credit card
// @Description Update an existing credit card
// @Tags credit-cards
// @Accept json
// @Produce json
// @Param id path string true "Credit Card ID"
// @Param creditCard body models.CreditCard true "Credit Card data"
// @Success 200 {object} models.CreditCard
// @Router /credit-cards/{id} [put]
func (h *CreditCardHandler) UpdateCreditCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit card id format"})
		return
	}

	var creditCard models.CreditCard
	if err := c.ShouldBindJSON(&creditCard); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creditCard.ID = id
	if err := h.creditCardService.UpdateCreditCard(&creditCard); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, creditCard)
}

// @Summary Delete credit card
// @Description Delete a credit card
// @Tags credit-cards
// @Accept json
// @Produce json
// @Param id path string true "Credit Card ID"
// @Success 204
// @Router /credit-cards/{id} [delete]
func (h *CreditCardHandler) DeleteCreditCard(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit card id format"})
		return
	}

	if err := h.creditCardService.DeleteCreditCard(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
