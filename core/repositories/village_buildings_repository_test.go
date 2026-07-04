package repositories

import (
	"context"
	"errors"
	"testing"
)

func TestInitialBuildings_DBError(t *testing.T) {
	db := &MockDBExecutor{
		execErr: errors.New("Connection failed"),
	}
	err := InitialBuildings(context.Background(), db, 1)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestInitialBuildings_Success(t *testing.T) {
	db := &MockDBExecutor{}
	err := InitialBuildings(context.Background(), db, 1)
	if err != nil {
		t.Errorf("Expected nil, got %s", err.Error())
	}
}
