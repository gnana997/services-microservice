package db

import (
	"testing"
	"time"
)

func TestInitialize_EmptyURL(t *testing.T) {
	_, err := Initialize("")
	if err == nil {
		t.Error("expected error for empty database URL")
	}
}

func TestMigrate_NilDB(t *testing.T) {
	err := Migrate(nil)
	if err == nil {
		t.Error("expected error for nil db")
	}
}

func TestConfigureConnectionPool_NilDB(t *testing.T) {
	err := ConfigureConnectionPool(nil, 1, 1, time.Second)
	if err == nil {
		t.Error("expected error for nil db")
	}
}

func TestClose_NilDB(t *testing.T) {
	err := Close(nil)
	if err == nil {
		t.Error("expected error for nil db")
	}
} 