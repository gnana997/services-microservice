package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"services-api/internal/business"
	"services-api/internal/models"
)

// ErrorResponse represents a standardized error response format
type ErrorResponse struct {
	Code    string `json:"code"`              // Machine-readable error code
	Message string `json:"message"`           // Human-readable error message
	Details any    `json:"details,omitempty"` // Optional additional details
}

// ServiceHandler handles service-related HTTP requests
type ServiceHandler struct {
	service business.BusinessService
}

// NewServiceHandler creates a new service handler with the required service business logic.
func NewServiceHandler(service business.BusinessService) *ServiceHandler {
	return &ServiceHandler{
		service: service,
	}
}

// ListServices godoc
// @Summary List services
// @Description Get a list of services with optional filtering, sorting, and pagination
// @Tags services
// @Param name query string false "Filter by service name (case-insensitive, partial match)"
// @Param description query string false "Filter by service description (case-insensitive, partial match)"
// @Param sort query string false "Sort field (name, created_at)" default(name)
// @Param order query string false "Sort order (asc, desc)" default(asc)
// @Param page query integer false "Page number" minimum(1) default(1)
// @Param limit query integer false "Items per page" minimum(1) maximum(100) default(10)
// @Success 200 {object} models.ServiceResponse "List of services"
// @Failure 500 {object} map[string]string "Error message"
// @Router /services [get]
func (h *ServiceHandler) ListServices(c *gin.Context) {
	filter := models.ServiceFilter{
		Name:        c.Query("name"),
		Description: c.Query("description"),
		Sort:        c.DefaultQuery("sort", "name"),
		Order:       c.DefaultQuery("order", "asc"),
		Page:        parseIntOrDefault(c.Query("page"), 1),
		Limit:       parseIntOrDefault(c.Query("limit"), 10),
	}

	// Validate pagination parameters
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 || filter.Limit > 100 {
		filter.Limit = 10
	}

	result, err := h.service.ListServices(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_error",
			Message: "An error occurred while retrieving services",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetService godoc
// @Summary Get a service
// @Description Get details of a service by ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path integer true "Service ID" minimum(1)
// @Success 200 {object} models.Service "Service details with versions"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Error message"
// @Router /services/{id} [get]
func (h *ServiceHandler) GetService(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_service_id",
			Message: "The service ID must be a positive integer",
			Details: fmt.Sprintf("Provided ID: %s", c.Param("id")),
		})
		return
	}

	svc, err := h.service.GetService(c.Request.Context(), uint(id))
	if err != nil {
		if err == business.ErrServiceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "service_not_found",
				Message: fmt.Sprintf("Service with ID %d could not be found", id),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_error",
			Message: "An error occurred while retrieving the service",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, svc)
}

// GetServiceVersions godoc
// @Summary Get a service
// @Description Get details of a service by ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path integer true "Service ID" minimum(1)
// @Success 200 {object} models.Service "Service details with versions"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Service not found"
// @Failure 500 {object} map[string]string "Error message"
// @Router /services/{id}/versions [get]
func (h *ServiceHandler) GetServiceVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_service_id",
			Message: "The service ID must be a positive integer",
			Details: fmt.Sprintf("Provided ID: %s", c.Param("id")),
		})
		return
	}

	versions, err := h.service.GetServiceVersions(c.Request.Context(), uint(id))
	if err != nil {
		if err == business.ErrServiceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "service_not_found",
				Message: fmt.Sprintf("Service with ID %d could not be found", id),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_error",
			Message: "An error occurred while retrieving service versions",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"versions": versions})
}

// GetServiceVersion godoc
// @Summary Get a service version
// @Description Get details of a service version by ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path integer true "Service ID" minimum(1)
// @Param versionId path integer true "Version ID" minimum(1)
// @Success 200 {object} models.Version "Service version details"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "Service version not found"
// @Failure 500 {object} map[string]string "Error message"
// @Router /services/{id}/versions/{versionId} [get]
func (h *ServiceHandler) GetServiceVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_service_id",
			Message: "The service ID must be a positive integer",
			Details: fmt.Sprintf("Provided ID: %s", c.Param("id")),
		})
		return
	}

	versionId, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Code:    "invalid_version_id",
			Message: "The version ID must be a positive integer",
			Details: fmt.Sprintf("Provided ID: %s", c.Param("versionId")),
		})
		return
	}

	version, err := h.service.GetServiceVersion(c.Request.Context(), uint(id), uint(versionId))
	if err != nil {
		if err == business.ErrServiceNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Code:    "service_not_found",
				Message: fmt.Sprintf("Service with ID %d could not be found", id),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Code:    "internal_error",
			Message: "An error occurred while retrieving the service version",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, version)
}

// parseIntOrDefault parses a string into an integer. 
// If parsing fails, it returns the provided default value.
func parseIntOrDefault(s string, defaultValue int) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultValue
}