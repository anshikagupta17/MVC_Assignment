package repositories

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestCollectResources_NormalCase(t *testing.T) {
	time_since_last := time.Now().Add(-100 * time.Second)

	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{500, 500}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{
				{
					int64(1),
					1,
					int64(6),
					pgtype.Timestamp{Time: time_since_last, Valid: true},
					nil,
					int64(6),
					float64(0.05),
				},
			}},
			{rows: [][]any{
				{int64(8), 1500},
				{int64(9), 1500},
			}},
		},
	}

	result, err := CollectResources(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Gold != 5 {
		t.Errorf("Expected 5 gold collected, got %d", result.Gold)
	}
	if result.Elixir != 0 {
		t.Errorf("Expected 0 elixir collected, got %d", result.Elixir)
	}
}

func TestCollectResources_StorageFull(t *testing.T) {
	time_since_last := time.Now().Add(-100 * time.Second)

	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1500, 1500}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{
				{
					int64(1),
					1,
					int64(6),
					pgtype.Timestamp{Time: time_since_last, Valid: true},
					nil,
					int64(6),
					float64(0.05),
				},
			}},
			{rows: [][]any{
				{int64(8), 1500},
				{int64(9), 1500},
			}},
		},
	}

	result, err := CollectResources(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Gold != 0 {
		t.Errorf("Expected 0 gold when storage is full, got %d", result.Gold)
	}
}

func TestCollectResources_MineUpgrading(t *testing.T) {
	time_since_last := time.Now().Add(-100 * time.Second)
	UpgradeFinishesLater := time.Now().Add(60 * time.Second)

	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{500, 500}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{
				{
					int64(1),
					1,
					int64(6),
					pgtype.Timestamp{Time: time_since_last, Valid: true},
					UpgradeFinishesLater,
					int64(6),
					float64(0.05),
				},
			}},
			{rows: [][]any{
				{int64(8), 1500},
				{int64(9), 1500},
			}},
		},
	}

	result, err := CollectResources(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Gold != 0 {
		t.Errorf("Expected 0 gold from upgrading mine, got %d", result.Gold)
	}
}

func TestCollectResources_NoTimeElapsed(t *testing.T) {
	now := time.Now()

	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{500, 500}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{
				{
					int64(1),
					1,
					int64(6),
					pgtype.Timestamp{Time: now, Valid: true},
					nil,
					int64(6),
					float64(0.05),
				},
			}},
			{rows: [][]any{
				{int64(8), 1500},
				{int64(9), 1500},
			}},
		},
	}

	result, err := CollectResources(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Gold != 0 {
		t.Errorf("Expected 0 gold when no time elapsed, got %d", result.Gold)
	}
}

func TestCollectResources_NoMines(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{500, 500}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{}},
			{rows: [][]any{
				{int64(8), 1500},
				{int64(9), 1500},
			}},
		},
	}

	result, err := CollectResources(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Gold != 0 {
		t.Errorf("Expected 0 gold with no mines, got %d", result.Gold)
	}
	if result.Elixir != 0 {
		t.Errorf("Expected 0 elixir with no mines, got %d", result.Elixir)
	}
}

func TestCollectResources_BothMineTypes(t *testing.T) {
	time_since_last := time.Now().Add(-100 * time.Second)

	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{500, 500}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{
				{
					int64(1),
					1,
					int64(6),
					pgtype.Timestamp{Time: time_since_last, Valid: true},
					nil,
					int64(6),
					float64(0.05),
				},
				{
					int64(2),
					1,
					int64(7),
					pgtype.Timestamp{Time: time_since_last, Valid: true},
					nil,
					int64(7),
					float64(0.05),
				},
			}},
			{rows: [][]any{
				{int64(8), 1500},
				{int64(9), 1500},
			}},
		},
	}

	result, err := CollectResources(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Gold != 5 {
		t.Errorf("Expected 5 gold collected, got %d", result.Gold)
	}
	if result.Elixir != 5 {
		t.Errorf("Expected 5 elixir collected, got %d", result.Elixir)
	}
}

func TestCollectResources_PartialStorageRoom(t *testing.T) {
	time_since_last := time.Now().Add(-400 * time.Second)

	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1490, 500}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{
				{
					int64(1),
					1,
					int64(6),
					pgtype.Timestamp{Time: time_since_last, Valid: true},
					nil,
					int64(6),
					float64(0.05),
				},
			}},
			{rows: [][]any{
				{int64(8), 1500},
				{int64(9), 1500},
			}},
		},
	}

	result, err := CollectResources(context.Background(), db, 1)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if result.Gold != 10 {
		t.Errorf("Expected 10 gold after partial cap, got %d", result.Gold)
	}
}

func TestCollectResources_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("connection failed")},
		},
	}

	_, err := CollectResources(context.Background(), db, 1)

	if err == nil {
		t.Errorf("Expected error when DB fails, got nil")
	}
}
