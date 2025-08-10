package handlers

import (
	"net/http"

	"github.com/Soli0222/flow-sight/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService *services.DashboardService
}

func NewDashboardHandler(dashboardService *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// @Summary Get dashboard summary
// @Description Get dashboard summary data including balance, monthly income/expense, asset count
// @Tags dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.DashboardSummary
// @Router /dashboard/summary [get]
func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	summary, err := h.dashboardService.GetDashboardSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
