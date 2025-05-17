package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"services-api/internal/business"
	"services-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockBusinessService struct {
	ListServicesFn        func(ctx context.Context, filter models.ServiceFilter) (*models.ServiceResponse, error)
	GetServiceFn          func(ctx context.Context, id uint) (*models.Service, error)
	GetServiceVersionFn   func(ctx context.Context, serviceID uint, versionID uint) (*models.Version, error)
	CreateServiceFn       func(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error)
	UpdateServiceFn       func(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error)
	DeleteServiceFn       func(ctx context.Context, id uint) error
}

func (m *mockBusinessService) ListServices(ctx context.Context, filter models.ServiceFilter) (*models.ServiceResponse, error) {
	return m.ListServicesFn(ctx, filter)
}
func (m *mockBusinessService) GetService(ctx context.Context, id uint) (*models.Service, error) {
	return m.GetServiceFn(ctx, id)
}
func (m *mockBusinessService) GetServiceVersion(ctx context.Context, serviceID uint, versionID uint) (*models.Version, error) {
	return m.GetServiceVersionFn(ctx, serviceID, versionID)
}
func (m *mockBusinessService) CreateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
	return m.CreateServiceFn(ctx, service)
}
func (m *mockBusinessService) UpdateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
	return m.UpdateServiceFn(ctx, service)
}
func (m *mockBusinessService) DeleteService(ctx context.Context, id uint) error {
	return m.DeleteServiceFn(ctx, id)
}


func TestListServicesHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockBusinessService{
		ListServicesFn: func(ctx context.Context, filter models.ServiceFilter) (*models.ServiceResponse, error) {
			return &models.ServiceResponse{Services: []models.ServiceModel{{ID: 1, Name: "Test Service"}}}, nil
		},
	}
	h := NewServiceHandler(mockSvc)
	r := gin.New()
	r.GET("/services", h.ListServices)

	req, _ := http.NewRequest("GET", "/services", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Service")
}

func TestGetServiceHandler_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockBusinessService{
		GetServiceFn: func(ctx context.Context, id uint) (*models.Service, error) {
			return nil, business.ErrServiceNotFound
		},
	}
	h := NewServiceHandler(mockSvc)
	r := gin.New()
	r.GET("/services/:id", h.GetService)

	req, _ := http.NewRequest("GET", "/services/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetServiceHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockBusinessService{
		GetServiceFn: func(ctx context.Context, id uint) (*models.Service, error) {
			return &models.Service{ID: 1, Name: "Test Service", Description: "Test Description", CreatedAt: time.Now(), UpdatedAt: time.Now(), VersionCount: 1}, nil
		},
	}
	h := NewServiceHandler(mockSvc)
	r := gin.New()
	r.GET("/services/:id", h.GetService)

	req, _ := http.NewRequest("GET", "/services/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Test Service")
}

func TestGetServiceVersionHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockBusinessService{
		GetServiceVersionFn: func(ctx context.Context, serviceID uint, versionID uint) (*models.Version, error) {
			return &models.Version{ID: versionID, Version: "1.0.0"}, nil
		},
	}
	h := NewServiceHandler(mockSvc)
	r := gin.New()
	r.GET("/services/:id/versions/:versionId", h.GetServiceVersion)

	req, _ := http.NewRequest("GET", "/services/1/versions/2", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "1.0.0")
} 

func TestCreateServiceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockBusinessService{
		CreateServiceFn: func(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
			return &models.ServiceModel{ID: 1, Name: "Test Service"}, nil
		},
	}
	h := NewServiceHandler(mockSvc)
	r := gin.New()
	r.POST("/services", h.CreateService)

	req, _ := http.NewRequest("POST", "/services", bytes.NewBufferString(`{"name": "Test Service"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Test Service")
}

func TestUpdateServiceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockBusinessService{
		UpdateServiceFn: func(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
			return &models.ServiceModel{ID: 1, Name: "Updated Service"}, nil
		},
	}
	h := NewServiceHandler(mockSvc)
	r := gin.New()
	r.PATCH("/services/:id", h.UpdateService)

	req, _ := http.NewRequest("PATCH", "/services/1", bytes.NewBufferString(`{"name": "Updated Service"}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Service")
}

func TestDeleteServiceHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := &mockBusinessService{
		DeleteServiceFn: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	h := NewServiceHandler(mockSvc)
	r := gin.New()
	r.DELETE("/services/:id", h.DeleteService)

	req, _ := http.NewRequest("DELETE", "/services/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
