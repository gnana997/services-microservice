package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"services-api/internal/models"
)

var (
	// ErrNotFound is returned when a requested record doesn't exist in the database
	ErrNotFound = errors.New("record not found")
)

// ServiceRepository interface defines data access methods for service entities
type ServiceRepository interface {
	// ListServices retrieves a paginated list of services based on filter criteria.
	// It returns the matched services, the total count of matches, and any error encountered.
	ListServices(ctx context.Context, filter models.ServiceFilter) ([]models.ServiceModel, int, error)
	
	// GetService retrieves a single service by its ID.
	// It returns the service with its associated versions or an error if not found.
	GetService(ctx context.Context, id uint) (*models.Service, error)
	
	// CreateService creates a new service
	// It returns the created service or an error if the service creation fails.
	CreateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error)

	// UpdateService updates a service
	// It returns the updated service or an error if the service update fails or if the service is not found.
	UpdateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error)

	// DeleteService deletes a service
	// It returns an error if the service deletion fails or if the service is not found.
	DeleteService(ctx context.Context, id uint) error
}

// serviceRepositoryImpl implements ServiceRepository
type serviceRepositoryImpl struct {
	db *gorm.DB
}

// NewServiceRepository creates a new service repository with the provided database connection.
func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepositoryImpl{
		db: db,
	}
}

// ListServices returns paginated services with filtering and sorting
func (r *serviceRepositoryImpl) ListServices(ctx context.Context, filter models.ServiceFilter) ([]models.ServiceModel, int, error) {
	var services []models.Service
	var servicesModel []models.ServiceModel = make([]models.ServiceModel, 0)
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Service{})

	// Apply filters
	if filter.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Description != "" {
		query = query.Where("description ILIKE ?", "%"+filter.Description+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortColumn := "name"
	if filter.Sort != "" {
		switch filter.Sort {
		case "name", "created_at", "updated_at":
			sortColumn = filter.Sort
		}
	}

	sortOrder := "ASC"
	if filter.Order == "desc" {
		sortOrder = "DESC"
	}
	query = query.Order(sortColumn + " " + sortOrder)

	// Apply pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Limit(filter.Limit).Offset(offset)

	// Execute query
	if err := query.Find(&services).Error; err != nil {
		return nil, 0, err
	}

	// If we have services, get the version counts efficiently in a single query
	if len(services) > 0 {
		// Extract service IDs
		serviceIDs := make([]uint, len(services))
		for i, service := range services {
			serviceIDs[i] = service.ID
		}

		// Create a struct to hold the count results
		type VersionCount struct {
			ServiceID uint
			Count     int
		}
		var counts []VersionCount

		// Execute a single query to get counts for all services
		err := r.db.WithContext(ctx).
			Model(&models.Version{}).
			Select("service_id, COUNT(*) as count").
			Where("service_id IN ?", serviceIDs).
			Group("service_id").
			Scan(&counts).Error

		if err != nil {
			return nil, 0, err
		}

		// Create a map for quick lookup of counts by service ID
		countMap := make(map[uint]int)
		for _, count := range counts {
			countMap[count.ServiceID] = count.Count
		}

		// Update services with their respective version counts
		for i := range services {
			var serviceModel models.ServiceModel
			serviceModel.ID = services[i].ID
			serviceModel.Name = services[i].Name
			serviceModel.Description = services[i].Description
			serviceModel.CreatedAt = services[i].CreatedAt
			serviceModel.UpdatedAt = services[i].UpdatedAt
			if count, exists := countMap[services[i].ID]; exists {
				serviceModel.VersionCount = count
			} else {
				serviceModel.VersionCount = 0
			}
			servicesModel = append(servicesModel, serviceModel)
		}
	}

	return servicesModel, int(total), nil
}

// GetService returns a single service by ID
func (r *serviceRepositoryImpl) GetService(ctx context.Context, id uint) (*models.Service, error) {
	var service models.Service
	
	err := r.db.WithContext(ctx).
		Preload("Versions").
		First(&service, id).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &service, nil
}

// GetServiceVersion returns a single version for a service
func (r *serviceRepositoryImpl) GetServiceVersion(ctx context.Context, serviceID uint, versionID uint) (*models.Version, error) {
	var version models.Version
	
	err := r.db.WithContext(ctx).
		Where("service_id = ? AND id = ?", serviceID, versionID).
		First(&version).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &version, nil
}

// CreateService creates a new service
// It returns the created service or an error if the service creation fails.
func (r *serviceRepositoryImpl) CreateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
	if err := r.db.WithContext(ctx).Model(&models.Service{}).Create(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

// UpdateService updates a service
func (r *serviceRepositoryImpl) UpdateService(ctx context.Context, service models.ServiceModel) (*models.ServiceModel, error) {
	if err := r.db.WithContext(ctx).Model(&models.Service{}).Where("id = ?", service.ID).Updates(&service).Error; err != nil {
		return nil, err
	}
	return &service, nil
}

// DeleteService deletes a service by ID and all the versions of the service
func (r *serviceRepositoryImpl) DeleteService(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Service{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	if err := r.db.WithContext(ctx).Where("service_id = ?", id).Delete(&models.Version{}).Error; err != nil {
		return err
	}
	return nil
}
