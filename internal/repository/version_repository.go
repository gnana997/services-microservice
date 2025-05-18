package repository

import (
	"context"
	"services-api/internal/models"

	"gorm.io/gorm"
)

// VersionRepository interface defines the operations for the version repository
type VersionRepository interface {
	// CreateVersion creates a new version
	// Returns the created version or an error if the version creation fails.
	CreateVersion(ctx context.Context, version models.Version) (*models.Version, error)

	// GetVersion retrieves a version by its ID
	// Returns the version or ErrNotFound if it doesn't exist.
	GetVersion(ctx context.Context, id uint, serviceId uint) (*models.Version, error)

	// UpdateVersion updates a version
	// Returns the updated version or an error if the version update fails or if the version is not found.
	UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error)

	// DeleteVersion deletes a version
	// Returns an error if the version deletion fails or if the version is not found.
	DeleteVersion(ctx context.Context, id uint, serviceId uint) error
}

type versionRepositoryImpl struct {
	db *gorm.DB
}

// NewVersionRepository creates a new version repository
// with the provided database connection.
func NewVersionRepository(db *gorm.DB) VersionRepository {
	return &versionRepositoryImpl{db: db}
}

// CreateVersion creates a new version
// Returns the created version or an error if the version creation fails.
func (r *versionRepositoryImpl) CreateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	if err := r.db.WithContext(ctx).Create(&version).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// GetVersion retrieves a version by its ID
// Returns the version or ErrNotFound if it doesn't exist.
func (r *versionRepositoryImpl) GetVersion(ctx context.Context, id uint, serviceId uint	) (*models.Version, error) {
	var version models.Version
	if err := r.db.WithContext(ctx).Where("id = ? AND service_id = ?", id, serviceId).First(&version).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// UpdateVersion updates a version
// Returns the updated version or an error if the version update fails or if the version is not found.
func (r *versionRepositoryImpl) UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	result := r.db.WithContext(ctx).Model(&models.Version{}).Where("id = ? AND service_id = ?", version.ID, version.ServiceID).Updates(&version)
	if result.Error != nil {
		return nil, result.Error
	}
	
	// Check if any rows were affected (record exists)
	if result.RowsAffected == 0 {
		return nil, ErrNotFound
	}
	
	// Fetch the updated record to return
	var updatedVersion models.Version
	if err := r.db.WithContext(ctx).Where("id = ? AND service_id = ?", version.ID, version.ServiceID).First(&updatedVersion).Error; err != nil {
		return nil, err
	}
	
	return &updatedVersion, nil
}

// DeleteVersion deletes a version
// Returns an error if the version deletion fails or if the version is not found.
func (r *versionRepositoryImpl) DeleteVersion(ctx context.Context, id uint, serviceId uint) error {
	result := r.db.WithContext(ctx).Where("id = ? AND service_id = ?", id, serviceId).Delete(&models.Version{})
	if result.Error != nil {
		return result.Error
	}
	
	// Check if any rows were affected (record exists)
	if result.RowsAffected == 0 {
		return ErrNotFound
	}
	
	return nil
}