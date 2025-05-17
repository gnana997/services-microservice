package config

import (
	"os"
	"testing"
)

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	tests := []struct {
		key          string
		defaultValue string
		expected     string
	}{
		{"TEST_KEY", "default", "test_value"},
		{"NON_EXISTENT_KEY", "default", "default"},
	}

	for _, tt := range tests {
		result := getEnvOrDefault(tt.key, tt.defaultValue)
		if result != tt.expected {
			t.Errorf("getEnvOrDefault(%q, %q) = %q; want %q", tt.key, tt.defaultValue, result, tt.expected)
		}
	}
}

func TestLoad(t *testing.T) {
	os.Setenv("DATABASE_URL", "postgres://test")
	os.Setenv("JWT_SECRET", "mysecret")
	os.Setenv("ENVIRONMENT", "production")
	defer os.Unsetenv("DATABASE_URL")
	defer os.Unsetenv("JWT_SECRET")
	defer os.Unsetenv("ENVIRONMENT")

	cfg := Load()
	if cfg.DatabaseURL != "postgres://test" {
		t.Errorf("expected DATABASE_URL to be 'postgres://test', got %q", cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "mysecret" {
		t.Errorf("expected JWT_SECRET to be 'mysecret', got %q", cfg.JWTSecret)
	}
	if cfg.Environment != "production" {
		t.Errorf("expected ENVIRONMENT to be 'production', got %q", cfg.Environment)
	}
} 