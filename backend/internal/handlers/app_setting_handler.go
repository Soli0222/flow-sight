package handlers

import (
	"flow-sight-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppSettingHandler struct {
	appSettingService *services.AppSettingService
}

func NewAppSettingHandler(appSettingService *services.AppSettingService) *AppSettingHandler {
	return &AppSettingHandler{
		appSettingService: appSettingService,
	}
}

// @Summary Get settings
// @Description Get all settings for a user
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.AppSetting
// @Router /settings [get]
func (h *AppSettingHandler) GetSettings(c *gin.Context) {
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

	settings, err := h.appSettingService.GetSettings(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

type UpdateSettingsRequest struct {
	Settings map[string]string `json:"settings"`
}

// @Summary Update settings
// @Description Update settings for a user
// @Tags settings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param settings body UpdateSettingsRequest true "Settings data"
// @Success 200 {object} map[string]string
// @Router /settings [put]
func (h *AppSettingHandler) UpdateSettings(c *gin.Context) {
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

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, value := range req.Settings {
		if err := h.appSettingService.UpdateSetting(userUUID, key, value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}
