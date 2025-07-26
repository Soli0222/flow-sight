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
// @Param user_id query string true "User ID"
// @Success 200 {array} models.AppSetting
// @Router /settings [get]
func (h *AppSettingHandler) GetSettings(c *gin.Context) {
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

	settings, err := h.appSettingService.GetSettings(userID)
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
// @Param user_id query string true "User ID"
// @Param settings body UpdateSettingsRequest true "Settings data"
// @Success 200 {object} map[string]string
// @Router /settings [put]
func (h *AppSettingHandler) UpdateSettings(c *gin.Context) {
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

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, value := range req.Settings {
		if err := h.appSettingService.UpdateSetting(userID, key, value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}
