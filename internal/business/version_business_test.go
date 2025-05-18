package business

import (
	"context"
	"testing"

	"services-api/internal/models"
)

type mockVersionRepository struct {
	CreateVersionFn func(ctx context.Context, version models.Version) (*models.Version, error)
	GetVersionFn    func(ctx context.Context, versionId uint, serviceId uint) (*models.Version, error)
	UpdateVersionFn func(ctx context.Context, version models.Version) (*models.Version, error)
	DeleteVersionFn func(ctx context.Context, versionId uint, serviceId uint) error
}

func (m *mockVersionRepository) CreateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	return &models.Version{ID: 1, Version: "1.0.0"}, nil
}
func (m *mockVersionRepository) GetVersion(ctx context.Context, versionId uint, serviceId uint) (*models.Version, error) {
	return &models.Version{ID: 1, Version: "1.0.0"}, nil
}
func (m *mockVersionRepository) UpdateVersion(ctx context.Context, version models.Version) (*models.Version, error) {
	return &models.Version{ID: 1, Version: "1.0.0"}, nil
}
func (m *mockVersionRepository) DeleteVersion(ctx context.Context, versionId uint, serviceId uint) error {
	return nil
}

func TestCreateVersion(t *testing.T) {
	repo := &mockVersionRepository{}
	business := NewVersionBusiness(repo)

	version := models.Version{
		Version: "1.0.0",
	}

	createdVersion, err := business.CreateVersion(context.Background(), version)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdVersion.ID != 1 {
		t.Fatalf("expected version ID to be 1, got %d", createdVersion.ID)
	}

	if createdVersion.Version != "1.0.0" {
		t.Fatalf("expected version to be 1.0.0, got %s", createdVersion.Version)
	}
}

func TestGetVersion(t *testing.T) {
	repo := &mockVersionRepository{}
	business := NewVersionBusiness(repo)

	version, err := business.GetVersion(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if version.ID != 1 {
		t.Fatalf("expected version ID to be 1, got %d", version.ID)
	}

	if version.Version != "1.0.0" {
		t.Fatalf("expected version to be 1.0.0, got %s", version.Version)
	}
}

func TestUpdateVersion(t *testing.T) {
	repo := &mockVersionRepository{}
	business := NewVersionBusiness(repo)

	version := models.Version{
		ID: 1,
		Version: "1.0.0",
	}

	updatedVersion, err := business.UpdateVersion(context.Background(), version)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedVersion.ID != 1 {
		t.Fatalf("expected version ID to be 1, got %d", updatedVersion.ID)
	}

	if updatedVersion.Version != "1.0.0" {
		t.Fatalf("expected version to be 1.0.0, got %s", updatedVersion.Version)
	}
}

func TestDeleteVersion(t *testing.T) {
	repo := &mockVersionRepository{}
	business := NewVersionBusiness(repo)

	err := business.DeleteVersion(context.Background(), 1, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
