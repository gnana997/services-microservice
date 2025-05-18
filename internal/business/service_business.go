package business

import (
	"context"
	"errors"

	"services-api/internal/models"
	"services-api/internal/repository"
)

var (
	// ErrServiceNotFound is returned when a requested service doesn't exist
	ErrServiceNotFound = errors.New("service not found")
)

// BusinessService interface defines service business logic operations
type BusinessService interface {
	// ListServices retrieves a paginated list of services based on filter criteria.
	// Returns a response containing services and pagination details.
	ListServices(ctx context.Context, filter models.ServiceFilter) (*models.ServiceResponse, error)
	
	// GetService retrieves a single service by its ID.
	// Returns the service with its associated versions or ErrServiceNotFound if not found.
	GetService(ctx context.Context, id uint) (*models.Service, error)

	// CreateService creates a new service
	// Returns the created service or an error if the service creation fails.
	CreateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error)

	// UpdateService updates a service
	// Returns the updated service or an error if the service update fails or if the service is not found.
	UpdateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error)

	// DeleteService deletes a service and all the versions of the service
	// Returns an error if the service deletion fails or if the service is not found.
	DeleteService(ctx context.Context, id uint) error
}

type serviceBusinessImpl struct {
	repo repository.ServiceRepository
}

// NewServiceBusiness creates a new business logic implementation
// with the provided repository.
func NewServiceBusiness(repo repository.ServiceRepository) BusinessService {
	return &serviceBusinessImpl{
		repo: repo,
	}
}

// ListServices returns a paginated list of services
func (s *serviceBusinessImpl) ListServices(ctx context.Context, filter models.ServiceFilter) (*models.ServiceResponse, error) {
	services, total, err := s.repo.ListServices(ctx, filter)
	if err != nil {
		return nil, err
	}

	totalPages := (total + filter.Limit - 1) / filter.Limit

	return &models.ServiceResponse{
		Services: services,
		Pagination: models.Pagination{
			CurrentPage:  filter.Page,
			TotalPages:   totalPages,
			TotalItems:   total,
			ItemsPerPage: filter.Limit,
		},
	}, nil
}

// GetService returns a single service by ID
func (s *serviceBusinessImpl) GetService(ctx context.Context, id uint) (*models.Service, error) {
	service, err := s.repo.GetService(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}
	return service, nil
}

// CreateService creates a new service
func (s *serviceBusinessImpl) CreateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
	return s.repo.CreateService(ctx, service)
}

// UpdateService updates a service
func (s *serviceBusinessImpl) UpdateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
	return s.repo.UpdateService(ctx, service)
}

// DeleteService deletes a service and all the versions of the service
func (s *serviceBusinessImpl) DeleteService(ctx context.Context, id uint) error {
	return s.repo.DeleteService(ctx, id)
}