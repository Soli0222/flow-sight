package handlers

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CardMonthlyTotalHandler struct {
	cardMonthlyTotalService *services.CardMonthlyTotalService
}

func NewCardMonthlyTotalHandler(cardMonthlyTotalService *services.CardMonthlyTotalService) *CardMonthlyTotalHandler {
	return &CardMonthlyTotalHandler{
		cardMonthlyTotalService: cardMonthlyTotalService,
	}
}

// @Summary Get card monthly totals
// @Description Get card monthly totals for a specific credit card
// @Tags card-monthly-totals
// @Accept json
// @Produce json
// @Param credit_card_id query string true "Credit Card ID"
// @Success 200 {array} models.CardMonthlyTotal
// @Router /card-monthly-totals [get]
func (h *CardMonthlyTotalHandler) GetCardMonthlyTotals(c *gin.Context) {
	creditCardIDStr := c.Query("credit_card_id")
	if creditCardIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "credit_card_id is required"})
		return
	}

	creditCardID, err := uuid.Parse(creditCardIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credit_card_id format"})
		return
	}

	totals, err := h.cardMonthlyTotalService.GetCardMonthlyTotals(creditCardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, totals)
}

// @Summary Get card monthly total by ID
// @Description Get a specific card monthly total by ID
// @Tags card-monthly-totals
// @Accept json
// @Produce json
// @Param id path string true "Card Monthly Total ID"
// @Success 200 {object} models.CardMonthlyTotal
// @Router /card-monthly-totals/{id} [get]
func (h *CardMonthlyTotalHandler) GetCardMonthlyTotal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card monthly total id format"})
		return
	}

	total, err := h.cardMonthlyTotalService.GetCardMonthlyTotal(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "card monthly total not found"})
		return
	}

	c.JSON(http.StatusOK, total)
}

// @Summary Create card monthly total
// @Description Create a new card monthly total
// @Tags card-monthly-totals
// @Accept json
// @Produce json
// @Param total body models.CardMonthlyTotal true "Card Monthly Total data"
// @Success 201 {object} models.CardMonthlyTotal
// @Router /card-monthly-totals [post]
func (h *CardMonthlyTotalHandler) CreateCardMonthlyTotal(c *gin.Context) {
	var total models.CardMonthlyTotal
	if err := c.ShouldBindJSON(&total); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.cardMonthlyTotalService.CreateCardMonthlyTotal(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, total)
}

// @Summary Update card monthly total
// @Description Update an existing card monthly total
// @Tags card-monthly-totals
// @Accept json
// @Produce json
// @Param id path string true "Card Monthly Total ID"
// @Param total body models.CardMonthlyTotal true "Card Monthly Total data"
// @Success 200 {object} models.CardMonthlyTotal
// @Router /card-monthly-totals/{id} [put]
func (h *CardMonthlyTotalHandler) UpdateCardMonthlyTotal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card monthly total id format"})
		return
	}

	var total models.CardMonthlyTotal
	if err := c.ShouldBindJSON(&total); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	total.ID = id
	if err := h.cardMonthlyTotalService.UpdateCardMonthlyTotal(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, total)
}

// @Summary Delete card monthly total
// @Description Delete a card monthly total
// @Tags card-monthly-totals
// @Accept json
// @Produce json
// @Param id path string true "Card Monthly Total ID"
// @Success 204
// @Router /card-monthly-totals/{id} [delete]
func (h *CardMonthlyTotalHandler) DeleteCardMonthlyTotal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid card monthly total id format"})
		return
	}

	if err := h.cardMonthlyTotalService.DeleteCardMonthlyTotal(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
