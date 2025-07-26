package handlers

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IncomeHandler struct {
	incomeService *services.IncomeService
}

func NewIncomeHandler(incomeService *services.IncomeService) *IncomeHandler {
	return &IncomeHandler{
		incomeService: incomeService,
	}
}

// Income Source handlers

// @Summary Get all income sources
// @Description Get all income sources for a user
// @Tags income
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} models.IncomeSource
// @Router /income-sources [get]
func (h *IncomeHandler) GetIncomeSources(c *gin.Context) {
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

	sources, err := h.incomeService.GetIncomeSources(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sources)
}

// @Summary Get income source by ID
// @Description Get a specific income source by ID
// @Tags income
// @Accept json
// @Produce json
// @Param id path string true "Income Source ID"
// @Success 200 {object} models.IncomeSource
// @Router /income-sources/{id} [get]
func (h *IncomeHandler) GetIncomeSource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid income source id format"})
		return
	}

	source, err := h.incomeService.GetIncomeSource(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "income source not found"})
		return
	}

	c.JSON(http.StatusOK, source)
}

// @Summary Create income source
// @Description Create a new income source
// @Tags income
// @Accept json
// @Produce json
// @Param source body models.IncomeSource true "Income Source data"
// @Success 201 {object} models.IncomeSource
// @Router /income-sources [post]
func (h *IncomeHandler) CreateIncomeSource(c *gin.Context) {
	var source models.IncomeSource
	if err := c.ShouldBindJSON(&source); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.incomeService.CreateIncomeSource(&source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, source)
}

// @Summary Update income source
// @Description Update an existing income source
// @Tags income
// @Accept json
// @Produce json
// @Param id path string true "Income Source ID"
// @Param source body models.IncomeSource true "Income Source data"
// @Success 200 {object} models.IncomeSource
// @Router /income-sources/{id} [put]
func (h *IncomeHandler) UpdateIncomeSource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid income source id format"})
		return
	}

	var source models.IncomeSource
	if err := c.ShouldBindJSON(&source); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	source.ID = id
	if err := h.incomeService.UpdateIncomeSource(&source); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, source)
}

// @Summary Delete income source
// @Description Delete an income source
// @Tags income
// @Accept json
// @Produce json
// @Param id path string true "Income Source ID"
// @Success 204
// @Router /income-sources/{id} [delete]
func (h *IncomeHandler) DeleteIncomeSource(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid income source id format"})
		return
	}

	if err := h.incomeService.DeleteIncomeSource(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// Monthly Income Record handlers

// @Summary Get monthly income records
// @Description Get monthly income records for a specific income source
// @Tags income
// @Accept json
// @Produce json
// @Param income_source_id query string true "Income Source ID"
// @Success 200 {array} models.MonthlyIncomeRecord
// @Router /monthly-income-records [get]
func (h *IncomeHandler) GetMonthlyIncomeRecords(c *gin.Context) {
	incomeSourceIDStr := c.Query("income_source_id")
	if incomeSourceIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "income_source_id is required"})
		return
	}

	incomeSourceID, err := uuid.Parse(incomeSourceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid income_source_id format"})
		return
	}

	records, err := h.incomeService.GetMonthlyIncomeRecords(incomeSourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

// @Summary Get monthly income record by ID
// @Description Get a specific monthly income record by ID
// @Tags income
// @Accept json
// @Produce json
// @Param id path string true "Monthly Income Record ID"
// @Success 200 {object} models.MonthlyIncomeRecord
// @Router /monthly-income-records/{id} [get]
func (h *IncomeHandler) GetMonthlyIncomeRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid monthly income record id format"})
		return
	}

	record, err := h.incomeService.GetMonthlyIncomeRecord(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "monthly income record not found"})
		return
	}

	c.JSON(http.StatusOK, record)
}

// @Summary Create monthly income record
// @Description Create a new monthly income record
// @Tags income
// @Accept json
// @Produce json
// @Param record body models.MonthlyIncomeRecord true "Monthly Income Record data"
// @Success 201 {object} models.MonthlyIncomeRecord
// @Router /monthly-income-records [post]
func (h *IncomeHandler) CreateMonthlyIncomeRecord(c *gin.Context) {
	var record models.MonthlyIncomeRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.incomeService.CreateMonthlyIncomeRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, record)
}

// @Summary Update monthly income record
// @Description Update an existing monthly income record
// @Tags income
// @Accept json
// @Produce json
// @Param id path string true "Monthly Income Record ID"
// @Param record body models.MonthlyIncomeRecord true "Monthly Income Record data"
// @Success 200 {object} models.MonthlyIncomeRecord
// @Router /monthly-income-records/{id} [put]
func (h *IncomeHandler) UpdateMonthlyIncomeRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid monthly income record id format"})
		return
	}

	var record models.MonthlyIncomeRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.ID = id
	if err := h.incomeService.UpdateMonthlyIncomeRecord(&record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, record)
}

// @Summary Delete monthly income record
// @Description Delete a monthly income record
// @Tags income
// @Accept json
// @Produce json
// @Param id path string true "Monthly Income Record ID"
// @Success 204
// @Router /monthly-income-records/{id} [delete]
func (h *IncomeHandler) DeleteMonthlyIncomeRecord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid monthly income record id format"})
		return
	}

	if err := h.incomeService.DeleteMonthlyIncomeRecord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
