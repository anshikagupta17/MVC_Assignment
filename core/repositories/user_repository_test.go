package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestCreateUser_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("Duplicate key")},
		},
	}
	_, err := CreateUser(context.Background(), db, "abc", "hashed_pass")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestCreateUser_DuplicateUsername(t *testing.T) {
	pgErr := &pgconn.PgError{Code: "23505"}

	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: pgErr},
		},
	}

	_, err := CreateUser(context.Background(), db, "abc", "hashed_pass")
	if err == nil {
		t.Fatalf("Expected error for duplicate username, got nil")
	}
	if err.Error() != "Username already taken" {
		t.Errorf("Expected 'Username already taken', got %s", err.Error())
	}
}

func TestCreateUser_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{int64(7)}},
		},
	}

	id, err := CreateUser(context.Background(), db, "abc", "hashed_pass")

	if err != nil {
		t.Fatalf("Expected no error, got %s", err.Error())
	}

	if id != 7 {
		t.Errorf("Expected id to be 7, got %d", id)
	}
}
