package handlers

import (
	"net/http"
	"strconv"

	"github.com/Soli0222/flow-sight/backend/internal/services"

	"github.com/gin-gonic/gin"
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
// @Description Get cashflow projection
// @Tags cashflow
// @Accept json
// @Produce json
// @Param months query int false "Number of months to project" default(36)
// @Param onlyChanges query bool false "Only return days with changes" default(false)
// @Success 200 {array} models.CashflowProjection
// @Router /cashflow-projection [get]
func (h *CashflowHandler) GetCashflowProjection(c *gin.Context) {
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

	projections, err := h.cashflowService.GetCashflowProjection(months, onlyChanges)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projections)
}
