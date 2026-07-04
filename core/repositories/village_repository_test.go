package repositories

import (
	"context"
	"errors"
	"testing"
)

func TestCreateVillage_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("Connection failed")},
		},
	}

	_, err := CreateVillage(context.Background(), db, 1)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestCreateVillage_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{int64(1)}},
		},
	}
	id, err := CreateVillage(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected nil, got %s", err.Error())
	}

	if id != 1 {
		t.Fatalf("Expected id as 1, got %d", id)
	}
}
