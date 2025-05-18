package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"services-api/internal/business"
	"services-api/internal/config"
	"services-api/internal/handlers"
	"services-api/internal/middleware"
	"services-api/internal/repository"
)

// HealthStatus represents the health check response
type HealthStatus struct {
	Status     string            `json:"status"`     // Overall service status
	Version    string            `json:"version"`    // API version
	Components map[string]string `json:"components"` // Status of individual components
}

// Server represents the API server
type Server struct {
	router  *gin.Engine
	db      *gorm.DB
	config  *config.Config
	version string
}

// NewServer creates a new API server
func NewServer(db *gorm.DB, cfg *config.Config) *Server {
	server := &Server{
		router:  gin.Default(),
		db:      db,
		config:  cfg,
		version: "1.0.0", // Set your API version here
	}
	
	// Set up routes immediately on creation
	server.setupRoutes()
	
	return server
}

// Router returns the gin.Engine instance
func (s *Server) Router() *gin.Engine {
	return s.router
}

// Run starts the server
func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}

// setupRoutes configures all the routes for the API server
func (s *Server) setupRoutes() {
	// Initialize repositories
	serviceRepo := repository.NewServiceRepository(s.db)
	versionRepo := repository.NewVersionRepository(s.db)
	
	// Initialize services
	serviceBusiness := business.NewServiceBusiness(serviceRepo)
	versionBusiness := business.NewVersionBusiness(versionRepo)
	
	// Initialize handlers
	serviceHandler := handlers.NewServiceHandler(serviceBusiness)
	versionHandler := handlers.NewVersionHandler(versionBusiness)
	
	// Initialize Swagger UI
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.NewHandler()))

	// Setup middleware
	s.router.Use(middleware.RequestLogger())
	s.router.Use(middleware.ErrorHandler())
	
	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// setup auth middleware for all the protected routes
		// v1.Use(middleware.AuthMiddleware())

		// Services endpoints
		services := v1.Group("/services")
		{
			services.GET("", serviceHandler.ListServices)
			services.GET("/:id", serviceHandler.GetService)
			services.POST("", serviceHandler.CreateService)
			services.PATCH("/:id", serviceHandler.UpdateService)
			services.DELETE("/:id", serviceHandler.DeleteService)
		}

		// Versions endpoints
		versions := v1.Group("/versions")
		{
			versions.POST("", versionHandler.CreateVersion)
			versions.GET("/:id", versionHandler.GetVersion)
			versions.PUT("/:id", versionHandler.UpdateVersion)
			versions.DELETE("/:id", versionHandler.DeleteVersion)
		}
		
	}

	// Enhanced health check
	s.router.GET("/health", s.healthCheck)
}

// healthCheck handles the /health endpoint
func (s *Server) healthCheck(c *gin.Context) {
	health := HealthStatus{
		Status:    "ok",
		Version:   s.version,
		Components: map[string]string{
			"database": "ok",
		},
	}
	
	// Check database connectivity
	sqlDB, err := s.db.DB()
	if err != nil {
		health.Status = "degraded"
		health.Components["database"] = "error: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}
	
	if err := sqlDB.Ping(); err != nil {
		health.Status = "degraded"
		health.Components["database"] = "error: " + err.Error()
		c.JSON(http.StatusServiceUnavailable, health)
		return
	}
	
	c.JSON(http.StatusOK, health)
}