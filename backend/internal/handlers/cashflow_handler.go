package handlers

import (
	"flow-sight-backend/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CashflowHandler struct {
	cashflowService *services.CashflowService
}

func NewCashflowHandler(cashflowService *services.CashflowService) *CashflowHandler {
	return &CashflowHandler{
		cashflowService: cashflowService,
	}
}

// @Summary Get cashflow projection
// @Description Get cashflow projection for a user
// @Tags cashflow
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Param months query int false "Number of months to project" default(6)
// @Success 200 {array} models.CashflowProjection
// @Router /cashflow-projection [get]
func (h *CashflowHandler) GetCashflowProjection(c *gin.Context) {
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

	monthsStr := c.DefaultQuery("months", "6")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid months parameter"})
		return
	}

	if months > 36 {
		months = 36 // Limit to 36 months as per specification
	}

	projections, err := h.cashflowService.GetCashflowProjection(userID, months)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projections)
}
