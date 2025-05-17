package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"services-api/internal/config"
	"services-api/internal/db"
	"services-api/internal/server"

	_ "services-api/docs"
)

// @title Services API
// @version 1.0
// @description API for managing services
// @BasePath /api/v1
// @host localhost:8080
// @schemes http

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	database, err := db.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Configure connection pool
	if err := db.ConfigureConnectionPool(database, 10, 100, time.Hour); err != nil {
		log.Fatal("Failed to configure connection pool:", err)
	}

	// Run migrations
	if err := db.Migrate(database); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize API server (routes are set up in the constructor)
	srv := server.NewServer(database, cfg)
	
	// Setup HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	// Create HTTP server with the router
	httpServer := &http.Server{
		Addr:    addr,
		Handler: srv.Router(),
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown the server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	// Close database connections
	if err := db.Close(database); err != nil {
		log.Fatal("Error closing database connection:", err)
	}

	log.Println("Server exited gracefully")
}