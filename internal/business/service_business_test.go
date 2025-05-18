package business

import (
	"context"
	"errors"
	"services-api/internal/models"
	"services-api/internal/repository"
	"testing"
	"time"
)

type mockRepo struct {
	ListServicesFn      func(ctx context.Context, filter models.ServiceFilter) ([]models.ServiceModel, int, error)
	GetServiceFn        func(ctx context.Context, id uint) (*models.Service, error)
	CreateServiceFn     func(ctx context.Context, service models.Service) (*models.Service, error)
	UpdateServiceFn     func(ctx context.Context, service models.Service) (*models.Service, error)
	DeleteServiceFn     func(ctx context.Context, id uint) error
}

func (m *mockRepo) ListServices(ctx context.Context, filter models.ServiceFilter) ([]models.ServiceModel, int, error) {
	return m.ListServicesFn(ctx, filter)
}
func (m *mockRepo) GetService(ctx context.Context, id uint) (*models.Service, error) {
	return m.GetServiceFn(ctx, id)
}
func (m *mockRepo) CreateService(ctx context.Context, service models.Service) (*models.Service, error) {
	return m.CreateServiceFn(ctx, service)
}
func (m *mockRepo) UpdateService(ctx context.Context, service models.Service) (*models.Service, error) {
	return m.UpdateServiceFn(ctx, service)
}
func (m *mockRepo) DeleteService(ctx context.Context, id uint) error {
	return m.DeleteServiceFn(ctx, id)
}

func TestListServices(t *testing.T) {
	repo := &mockRepo{
		ListServicesFn: func(ctx context.Context, filter models.ServiceFilter) ([]models.ServiceModel, int, error) {
			return []models.ServiceModel{{ID: 1, Name: "Test Service"}}, 1, nil
		},
	}
	bs := NewServiceBusiness(repo)
	resp, err := bs.ListServices(context.Background(), models.ServiceFilter{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Services) != 1 || resp.Services[0].Name != "Test Service" {
		t.Errorf("unexpected response: %+v", resp)
	}
}

func TestGetService_NotFound(t *testing.T) {
	repo := &mockRepo{
		GetServiceFn: func(ctx context.Context, id uint) (*models.Service, error) {
			return nil, repository.ErrNotFound
		},
	}
	bs := NewServiceBusiness(repo)
	_, err := bs.GetService(context.Background(), 1)
	if !errors.Is(err, ErrServiceNotFound) {
		t.Errorf("expected ErrServiceNotFound, got %v", err)
	}
}

func TestGetService_Success(t *testing.T) {
	repo := &mockRepo{
		GetServiceFn: func(ctx context.Context, id uint) (*models.Service, error) {
			return &models.Service{ID: 1, Name: "Test Service", Description: "Test Description", CreatedAt: time.Now(), UpdatedAt: time.Now(), VersionCount: 1}, nil
		},
	}
	bs := NewServiceBusiness(repo)
	service, err := bs.GetService(context.Background(), 1)
	if err != nil || service.ID != 1 {
		t.Errorf("unexpected result: %v, %v", service, err)
	}
}

func TestCreateService(t *testing.T) {
	repo := &mockRepo{
		CreateServiceFn: func(ctx context.Context, service models.Service) (*models.Service, error) {
			return &models.Service{ID: 1, Name: "Test Service"}, nil
		},
	}
	bs := NewServiceBusiness(repo)
	service, err := bs.CreateService(context.Background(), models.Service{Name: "Test Service"})
	if err != nil || service.ID != 1 {
		t.Errorf("unexpected result: %v, %v", service, err)
	}
}

func TestUpdateService(t *testing.T) {
	repo := &mockRepo{
		UpdateServiceFn: func(ctx context.Context, service models.Service) (*models.Service, error) {
			return &models.Service{ID: 1, Name: "Updated Service"}, nil
		},
	}
	bs := NewServiceBusiness(repo)
	service, err := bs.UpdateService(context.Background(), models.Service{ID: 1, Name: "Updated Service"})
	if err != nil || service.ID != 1 {
		t.Errorf("unexpected result: %v, %v", service, err)
	}
}

func TestDeleteService(t *testing.T) {
	repo := &mockRepo{
		DeleteServiceFn: func(ctx context.Context, id uint) error {
			return nil
		},
	}
	bs := NewServiceBusiness(repo)
	err := bs.DeleteService(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}