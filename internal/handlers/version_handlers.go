package handlers

import (
	"net/http"
	"strconv"

	"services-api/internal/business"
	"services-api/internal/models"

	"github.com/gin-gonic/gin"
)

type VersionHandler struct {
	versionBusiness business.VersionBusiness
}

func NewVersionHandler(versionBusiness business.VersionBusiness) *VersionHandler {
	return &VersionHandler{
		versionBusiness: versionBusiness,
	}
}

// CreateVersion godoc
// @Summary Create a new version
// @Description Create a new version for a service
// @Tags versions
// @Accept json
// @Produce json
// @Param version body models.Version true "Version details"
// @Success 201 {object} models.Version "Created version"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to create version"
// @Router /versions [post]
func (h *VersionHandler) CreateVersion(c *gin.Context) {
	var version models.Version
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_request_body",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	createdVersion, err := h.versionBusiness.CreateVersion(c.Request.Context(), version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_server_error",
			Message: "Failed to create version",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdVersion)
}

// GetVersion godoc
// @Summary Get a version by ID
// @Description Get a version by ID
// @Tags versions
// @Accept json
// @Produce json
// @Param id path integer true "Version ID"
// @Success 200 {object} models.Version "Version details"
// @Failure 400 {object} ErrorResponse "Invalid version ID"
// @Failure 500 {object} ErrorResponse "Failed to get version"
// @Router /versions/{id} [get]
func (h *VersionHandler) GetVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_version_id",
			Message: "Invalid version ID",
			Details: err.Error(),
		})
		return
	}

	version, err := h.versionBusiness.GetVersion(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_server_error",
			Message: "Failed to get version",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, version)
}

// UpdateVersion godoc
// @Summary Update a version
// @Description Update a version by ID
// @Tags versions
// @Accept json
// @Produce json
// @Param version body models.Version true "Version details"
// @Success 200 {object} models.Version "Updated version"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to update version"
// @Router /versions/{id} [put]
func (h *VersionHandler) UpdateVersion(c *gin.Context) {
	var version models.Version
	if err := c.ShouldBindJSON(&version); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_request_body",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	updatedVersion, err := h.versionBusiness.UpdateVersion(c.Request.Context(), version)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_server_error",
			Message: "Failed to update version",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updatedVersion)
}

// DeleteVersion godoc
// @Summary Delete a version
// @Description Delete a version by ID
// @Tags versions
// @Accept json
// @Produce json
// @Param id path integer true "Version ID"
// @Success 204 "Version deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid version ID"
// @Failure 500 {object} ErrorResponse "Failed to delete version"
// @Router /versions/{id} [delete]
func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_version_id",
			Message: "Invalid version ID",
			Details: err.Error(),
		})
		return
	}

	err = h.versionBusiness.DeleteVersion(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_server_error",
			Message: "Failed to delete version",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}