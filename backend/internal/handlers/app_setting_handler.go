package handlers

import (
	"net/http"

	"github.com/Soli0222/flow-sight/backend/internal/services"

	"github.com/gin-gonic/gin"
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
// @Description Get all settings
// @Tags settings
// @Accept json
// @Produce json
// @Success 200 {array} models.AppSetting
// @Router /settings [get]
func (h *AppSettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.appSettingService.GetSettings()
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
// @Description Update settings
// @Tags settings
// @Accept json
// @Produce json
// @Param settings body UpdateSettingsRequest true "Settings data"
// @Success 200 {object} map[string]string
// @Router /settings [put]
func (h *AppSettingHandler) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for key, value := range req.Settings {
		if err := h.appSettingService.UpdateSetting(key, value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}
