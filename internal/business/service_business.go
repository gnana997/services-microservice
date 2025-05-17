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
	
	// GetServiceVersions retrieves all versions for a specific service.
	// Returns ErrServiceNotFound if the service doesn't exist.
	GetServiceVersions(ctx context.Context, serviceID uint) ([]models.Version, error)

	// GetServiceVersion retrieves a single version for a service
	GetServiceVersion(ctx context.Context, serviceID uint, versionID uint) (*models.Version, error)
}

// businessServiceImpl implements BusinessService
type businessServiceImpl struct {
	repo repository.ServiceRepository
}

// NewBusinessService creates a new business logic implementation
// with the provided repository.
func NewBusinessService(repo repository.ServiceRepository) BusinessService {
	return &businessServiceImpl{
		repo: repo,
	}
}

// ListServices returns a paginated list of services
func (s *businessServiceImpl) ListServices(ctx context.Context, filter models.ServiceFilter) (*models.ServiceResponse, error) {
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
func (s *businessServiceImpl) GetService(ctx context.Context, id uint) (*models.Service, error) {
	service, err := s.repo.GetService(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}
	return service, nil
}

// GetServiceVersions returns all versions for a service
func (s *businessServiceImpl) GetServiceVersions(ctx context.Context, serviceID uint) ([]models.Version, error) {
	// First check if service exists
	_, err := s.repo.GetService(ctx, serviceID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrServiceNotFound
		}
		return nil, err
	}

	return s.repo.GetServiceVersions(ctx, serviceID)
}

// GetServiceVersion returns a single version for a service
func (s *businessServiceImpl) GetServiceVersion(ctx context.Context, serviceID uint, versionID uint) (*models.Version, error) {
	return s.repo.GetServiceVersion(ctx, serviceID, versionID)
}

