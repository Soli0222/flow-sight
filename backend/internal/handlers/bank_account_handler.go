package handlers

import (
	"context"
	"net/http"

	"github.com/Soli0222/flow-sight/backend/internal/middleware"
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BankAccountHandler struct {
	bankAccountService BankAccountServiceInterface
}

func NewBankAccountHandler(bankAccountService BankAccountServiceInterface) *BankAccountHandler {
	return &BankAccountHandler{
		bankAccountService: bankAccountService,
	}
}

// @Summary Get all bank accounts
// @Description Get all bank accounts
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Success 200 {array} models.BankAccount
// @Router /bank-accounts [get]
func (h *BankAccountHandler) GetBankAccounts(c *gin.Context) {
	accounts, err := h.bankAccountService.GetBankAccounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

// @Summary Get bank account by ID
// @Description Get a specific bank account by ID
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Success 200 {object} models.BankAccount
// @Router /bank-accounts/{id} [get]
func (h *BankAccountHandler) GetBankAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank account id format"})
		return
	}

	account, err := h.bankAccountService.GetBankAccount(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bank account not found"})
		return
	}

	c.JSON(http.StatusOK, account)
}

// @Summary Create bank account
// @Description Create a new bank account
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param account body models.BankAccount true "Bank Account data"
// @Success 201 {object} models.BankAccount
// @Router /bank-accounts [post]
func (h *BankAccountHandler) CreateBankAccount(c *gin.Context) {
	logger := middleware.GetLogger(c)
	ctx := context.Background()

	var account models.BankAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		logger.WarnContext(ctx, "Invalid request body for bank account creation",
			"error", err.Error(),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bankAccountService.CreateBankAccount(&account); err != nil {
		logger.ErrorContext(ctx, "Failed to create bank account",
			"account_name", account.Name,
			"error", err.Error(),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.BusinessOperation(ctx, "bank_account_created", "single-user", map[string]interface{}{
		"account_id":   account.ID.String(),
		"account_name": account.Name,
		"balance":      account.Balance,
	})

	c.JSON(http.StatusCreated, account)
}

// @Summary Update bank account
// @Description Update an existing bank account
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Param account body models.BankAccount true "Bank Account data"
// @Success 200 {object} models.BankAccount
// @Router /bank-accounts/{id} [put]
func (h *BankAccountHandler) UpdateBankAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank account id format"})
		return
	}

	var account models.BankAccount
	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	account.ID = id
	if err := h.bankAccountService.UpdateBankAccount(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, account)
}

// @Summary Delete bank account
// @Description Delete a bank account
// @Tags bank-accounts
// @Accept json
// @Produce json
// @Param id path string true "Bank Account ID"
// @Success 204
// @Router /bank-accounts/{id} [delete]
func (h *BankAccountHandler) DeleteBankAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid bank account id format"})
		return
	}

	if err := h.bankAccountService.DeleteBankAccount(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
