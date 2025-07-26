package handlers

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssetHandler struct {
	assetService *services.AssetService
}

func NewAssetHandler(assetService *services.AssetService) *AssetHandler {
	return &AssetHandler{
		assetService: assetService,
	}
}

// @Summary Get all assets
// @Description Get all assets for a user
// @Tags assets
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} models.Asset
// @Router /assets [get]
func (h *AssetHandler) GetAssets(c *gin.Context) {
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

	assets, err := h.assetService.GetAssets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, assets)
}

// @Summary Get asset by ID
// @Description Get a specific asset by ID
// @Tags assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Success 200 {object} models.Asset
// @Router /assets/{id} [get]
func (h *AssetHandler) GetAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id format"})
		return
	}

	asset, err := h.assetService.GetAsset(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	c.JSON(http.StatusOK, asset)
}

// @Summary Create asset
// @Description Create a new asset
// @Tags assets
// @Accept json
// @Produce json
// @Param asset body models.Asset true "Asset data"
// @Success 201 {object} models.Asset
// @Router /assets [post]
func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.assetService.CreateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, asset)
}

// @Summary Update asset
// @Description Update an existing asset
// @Tags assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Param asset body models.Asset true "Asset data"
// @Success 200 {object} models.Asset
// @Router /assets/{id} [put]
func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id format"})
		return
	}

	var asset models.Asset
	if err := c.ShouldBindJSON(&asset); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset.ID = id
	if err := h.assetService.UpdateAsset(&asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, asset)
}

// @Summary Delete asset
// @Description Delete an asset
// @Tags assets
// @Accept json
// @Produce json
// @Param id path string true "Asset ID"
// @Success 204
// @Router /assets/{id} [delete]
func (h *AssetHandler) DeleteAsset(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id format"})
		return
	}

	if err := h.assetService.DeleteAsset(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
