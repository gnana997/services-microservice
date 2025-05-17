package models

import (
	"time"
)

// Service represents a service in the organization
type Service struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name        string    `json:"name" gorm:"unique;not null;index" example:"User Service"`
	Description string    `json:"description" example:"Manages user authentication and profiles"`
	CreatedAt   time.Time `json:"created_at" example:"2025-05-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-05-01T00:00:00Z"`
	Versions    []Version `json:"versions" gorm:"foreignKey:ServiceID"`
	VersionCount int       `json:"version_count" gorm:"-" example:"1"`
}

type ServiceModel struct {
	ID uint `json:"id" example:"1"`
	Name string `json:"name" example:"User Service"`
	Description string `json:"description" example:"Manages user authentication and profiles"`
	CreatedAt time.Time `json:"created_at" example:"2025-05-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-05-01T00:00:00Z"`
	VersionCount int `json:"version_count" example:"1"`
}

// Version represents a version of a service
type Version struct {
	ID          uint      `json:"id" gorm:"primaryKey" example:"1"`
	ServiceID   uint      `json:"service_id" gorm:"not null;index" example:"1"`
	Version     string    `json:"version" gorm:"not null" example:"1.0.0"`
	Description string    `json:"description" example:"Initial release"`
	IsActive    bool      `json:"is_active" gorm:"default:true" example:"true"`
	CreatedAt   time.Time `json:"created_at" example:"2025-05-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-05-01T00:00:00Z"`
}

// ServiceFilter contains filter parameters for service queries
type ServiceFilter struct {
	Name        string `json:"name" example:"auth"`
	Description string `json:"description" example:"authentication"`
	Sort        string `json:"sort" example:"name"`
	Order       string `json:"order" example:"asc"`
	Page        int    `json:"page" example:"1"`
	Limit       int    `json:"limit" example:"10"`
}

// ServiceResponse represents the response for service list
type ServiceResponse struct {
	Services   []ServiceModel   `json:"services"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination contains pagination information
type Pagination struct {
	CurrentPage  int `json:"current_page" example:"1"`
	TotalPages   int `json:"total_pages" example:"5"`
	TotalItems   int `json:"total_items" example:"50"`
	ItemsPerPage int `json:"items_per_page" example:"10"`
}