package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"services-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockVersionBusiness struct {
	CreateVersionFn func(ctx context.Context, version models.Version) (*models.Version, error)
	GetVersionFn    func(ctx context.Context, id uint, serviceId uint) (*models.Version, error)
	UpdateVersionFn func(ctx context.Context, version models.Version) (*models.Version, error)
	DeleteVersionFn func(ctx context.Context, id uint, serviceId uint) error
}

func (m *mockVersionBusiness) CreateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	return m.CreateVersionFn(ctx, version)
}
func (m *mockVersionBusiness) GetVersion(ctx context.Context, id uint, serviceId uint) (*models.Version, error) {
	return m.GetVersionFn(ctx, id, serviceId)
}
func (m *mockVersionBusiness) UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	return m.UpdateVersionFn(ctx, version)
}
func (m *mockVersionBusiness) DeleteVersion(ctx context.Context, id uint, serviceId uint) error {
	return m.DeleteVersionFn(ctx, id, serviceId)
}

func TestCreateVersion_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{
		CreateVersionFn: func(ctx context.Context, version models.Version) (*models.Version, error) {
			version.ID = 1
			return &version, nil
		},
	}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.POST("/services/:sid/versions", h.CreateVersion)

	payload := `{"service_id":1,"version":"1.0.0","description":"desc"}`
	req, _ := http.NewRequest("POST", "/services/1/versions", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "1.0.0")
}

func TestCreateVersion_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.POST("/services/:sid/versions", h.CreateVersion)

	req, _ := http.NewRequest("POST", "/services/1/versions", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetVersion_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{
		GetVersionFn: func(ctx context.Context, id uint, serviceId uint) (*models.Version, error) {
			return &models.Version{ID: id, Version: "1.0.0"}, nil
		},
	}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.GET("/services/:sid/versions/:vid", h.GetVersion)

	req, _ := http.NewRequest("GET", "/services/1/versions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "1.0.0")
}

func TestGetVersion_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.GET("/services/:sid/versions/:vid", h.GetVersion)

	req, _ := http.NewRequest("GET", "/services/bad-id/versions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetVersion_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{
		GetVersionFn: func(ctx context.Context, id uint, serviceId uint) (*models.Version, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.GET("/services/:sid/versions/:vid", h.GetVersion)

	req, _ := http.NewRequest("GET", "/services/1/versions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateVersion_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{
		UpdateVersionFn: func(ctx context.Context, version models.Version) (*models.Version, error) {
			version.Description = "updated"
			return &version, nil
		},
	}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.PUT("/services/:sid/versions/:vid", h.UpdateVersion)

	payload := `{"service_id":1,"version":"1.0.0","description":"desc"}`
	req, _ := http.NewRequest("PUT", "/services/1/versions/1", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "updated")
}

func TestUpdateVersion_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.PUT("/services/:sid/versions/:vid", h.UpdateVersion)

	req, _ := http.NewRequest("PUT", "/services/1/versions/1", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteVersion_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{
		DeleteVersionFn: func(ctx context.Context, id uint, serviceId uint) error {
			return nil
		},
	}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.DELETE("/services/:sid/versions/:vid", h.DeleteVersion)

	req, _ := http.NewRequest("DELETE", "/services/1/versions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteVersion_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.DELETE("/services/:sid/versions/:vid", h.DeleteVersion)

	req, _ := http.NewRequest("DELETE", "/services/bad-id/versions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteVersion_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockBiz := &mockVersionBusiness{
		DeleteVersionFn: func(ctx context.Context, id uint, serviceId uint) error {
			return errors.New("db error")
		},
	}
	h := NewVersionHandler(mockBiz)
	r := gin.New()
	r.DELETE("/services/:sid/versions/:vid", h.DeleteVersion)

	req, _ := http.NewRequest("DELETE", "/services/1/versions/1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
