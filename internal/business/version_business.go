package business

import (
	"context"
	"errors"
	"services-api/internal/models"
	"services-api/internal/repository"
)

var (
	// ErrVersionNotFound is returned when a requested version doesn't exist
	ErrVersionNotFound = errors.New("version not found")
)

type VersionBusiness interface {
	// CreateVersion creates a new version
	// Returns the created version or an error if the version creation fails.
	CreateVersion(ctx context.Context, version models.Version) (*models.Version, error)

	// GetVersion retrieves a version by its ID
	// Returns the version or ErrVersionNotFound if it doesn't exist.
	GetVersion(ctx context.Context, id uint) (*models.Version, error)

	// UpdateVersion updates a version
	// Returns the updated version or an error if the version update fails or if the version is not found.
	UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error)

	// DeleteVersion deletes a version
	// Returns an error if the version deletion fails or if the version is not found.
	DeleteVersion(ctx context.Context, id uint) error
}

type versionBusinessImpl struct {
	repo repository.VersionRepository
}

// NewVersionBusiness creates a new business logic implementation
// with the provided repository.
func NewVersionBusiness(repo repository.VersionRepository) VersionBusiness {
	return &versionBusinessImpl{
		repo: repo,
	}
}

// CreateVersion creates a new version
// Returns the created version or an error if the version creation fails.
func (b *versionBusinessImpl) CreateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	return b.repo.CreateVersion(ctx, version)
}

// GetVersion retrieves a version by its ID
// Returns the version or ErrVersionNotFound if it doesn't exist.
func (b *versionBusinessImpl) GetVersion(ctx context.Context, id uint) (*models.Version, error) {
	version, err := b.repo.GetVersion(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrVersionNotFound
		}
		return nil, err
	}
	return version, nil
}

// UpdateVersion updates a version
// Returns the updated version or an error if the version update fails or if the version is not found.
func (b *versionBusinessImpl) UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	return b.repo.UpdateVersion(ctx, version)
}

// DeleteVersion deletes a version
// Returns an error if the version deletion fails or if the version is not found.
func (b *versionBusinessImpl) DeleteVersion(ctx context.Context, id uint) error {
	return b.repo.DeleteVersion(ctx, id)
}