package repository

import (
	"context"
	"errors"
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
	GetVersion(ctx context.Context, id uint) (*models.Version, error)

	// UpdateVersion updates a version
	// Returns the updated version or an error if the version update fails or if the version is not found.
	UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error)

	// DeleteVersion deletes a version
	// Returns an error if the version deletion fails or if the version is not found.
	DeleteVersion(ctx context.Context, id uint) error
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
func (r *versionRepositoryImpl) GetVersion(ctx context.Context, id uint) (*models.Version, error) {
	var version models.Version
	if err := r.db.WithContext(ctx).First(&version, id).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// UpdateVersion updates a version
// Returns the updated version or an error if the version update fails or if the version is not found.
func (r *versionRepositoryImpl) UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	if err := r.db.WithContext(ctx).Model(&models.Version{}).Where("id = ?", version.ID).Updates(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &version, nil
}

// DeleteVersion deletes a version
// Returns an error if the version deletion fails or if the version is not found.
func (r *versionRepositoryImpl) DeleteVersion(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&models.Version{}, id).Error; err != nil {
		return err
	}
	return nil
}