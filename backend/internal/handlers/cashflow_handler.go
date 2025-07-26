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
// @Security BearerAuth
// @Param months query int false "Number of months to project" default(36)
// @Param onlyChanges query bool false "Only return days with changes" default(false)
// @Success 200 {array} models.CashflowProjection
// @Router /cashflow-projection [get]
func (h *CashflowHandler) GetCashflowProjection(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id format in context"})
		return
	}

	monthsStr := c.DefaultQuery("months", "36")
	months, err := strconv.Atoi(monthsStr)
	if err != nil || months <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid months parameter"})
		return
	}

	if months > 120 {
		months = 120 // Limit to 120 months (10 years)
	}

	onlyChangesStr := c.DefaultQuery("onlyChanges", "false")
	onlyChanges := onlyChangesStr == "true"

	projections, err := h.cashflowService.GetCashflowProjection(userUUID, months, onlyChanges)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projections)
}
