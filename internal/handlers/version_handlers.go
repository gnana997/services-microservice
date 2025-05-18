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
// @Param sid path integer true "Service ID"
// @Param version body models.VersionRequest true "Version details"
// @Success 201 {object} models.Version "Created version"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 500 {object} ErrorResponse "Failed to create version"
// @Router /services/{sid}/versions [post]
func (h *VersionHandler) CreateVersion(c *gin.Context) {
	serviceId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_service_id",
			Message: "Invalid service ID",
			Details: err.Error(),
		})
		return
	}

	var req models.VersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_request_body",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	version := models.Version{
		ServiceID:   uint(serviceId),
		Version:     req.Version,
		Description: req.Description,
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
// @Param sid path integer true "Service ID"
// @Param vid path integer true "Version ID"
// @Success 200 {object} models.Version "Version details"
// @Failure 400 {object} ErrorResponse "Invalid version ID"
// @Failure 404 {object} ErrorResponse "Version not found"
// @Failure 500 {object} ErrorResponse "Failed to get version"
// @Router /services/{sid}/versions/{vid} [get]
func (h *VersionHandler) GetVersion(c *gin.Context) {
	serviceId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_service_id",
			Message: "Invalid service ID",
			Details: err.Error(),
		})
		return
	}

	versionId, err := strconv.ParseUint(c.Param("vid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_version_id",
			Message: "Invalid version ID",
			Details: err.Error(),
		})
		return
	}

	version, err := h.versionBusiness.GetVersion(c.Request.Context(), uint(serviceId), uint(versionId))
	if err != nil {
		if err == business.ErrVersionNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "version_not_found",
				Message: "Version not found",
				Details: "The requested version does not exist",
			})
			return
		}
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
// @Param sid path integer true "Service ID"
// @Param vid path integer true "Version ID"
// @Param version body models.VersionRequest true "Version details"
// @Success 200 {object} models.Version "Updated version"
// @Failure 400 {object} ErrorResponse "Invalid request body"
// @Failure 404 {object} ErrorResponse "Version not found"
// @Failure 500 {object} ErrorResponse "Failed to update version"
// @Router /services/{sid}/versions/{vid} [put]
func (h *VersionHandler) UpdateVersion(c *gin.Context) {
	serviceId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_service_id",
			Message: "Invalid service ID",
			Details: err.Error(),
		})
		return
	}

	versionId, err := strconv.ParseUint(c.Param("vid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_version_id",
			Message: "Invalid version ID",
			Details: err.Error(),
		})
		return
	}

	var req models.VersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_request_body",
			Message: "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	version := models.Version{
		ID:          uint(versionId),
		ServiceID:   uint(serviceId),
		Version:     req.Version,
		Description: req.Description,
	}

	updatedVersion, err := h.versionBusiness.UpdateVersion(c.Request.Context(), version)
	if err != nil {
		if err == business.ErrVersionNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "version_not_found",
				Message: "Version not found",
				Details: "The requested version does not exist",
			})
			return
		}
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
// @Param sid path integer true "Service ID"
// @Param vid path integer true "Version ID"
// @Success 204 "Version deleted successfully"
// @Failure 400 {object} ErrorResponse "Invalid version ID"
// @Failure 404 {object} ErrorResponse "Version not found"
// @Failure 500 {object} ErrorResponse "Failed to delete version"
// @Router /services/{sid}/versions/{vid} [delete]
func (h *VersionHandler) DeleteVersion(c *gin.Context) {
	serviceId, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_service_id",
			Message: "Invalid service ID",
			Details: err.Error(),
		})
		return
	}

	versionId, err := strconv.ParseUint(c.Param("vid"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_version_id",
			Message: "Invalid version ID",
			Details: err.Error(),
		})
		return
	}

	err = h.versionBusiness.DeleteVersion(c.Request.Context(), uint(versionId), uint(serviceId))
	if err != nil {
		if err == business.ErrVersionNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "version_not_found",
				Message: "Version not found",
				Details: "The requested version does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_server_error",
			Message: "Failed to delete version",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}