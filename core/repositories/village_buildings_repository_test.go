package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
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

//

func TestCanPlace_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryErr: errors.New("connection failed"),
	}
	_, err := CanPlace(context.Background(), db, 1, 1, 1, 5, 5, 1)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestCanPlace_OutOfBounds(t *testing.T) {
	db := &MockDBExecutor{}

	canPlace, err := CanPlace(context.Background(), db, 1, 3, 3, 48, 5, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if canPlace {
		t.Errorf("Expected false for out of bounds placement, got true")
	}
}

func TestCanPlace_NoExistingBuildings(t *testing.T) {
	db := &MockDBExecutor{
		queryResults: []*MockRows{
			{rows: [][]any{}},
		},
	}

	canPlace, err := CanPlace(context.Background(), db, 1, 1, 1, 5, 5, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !canPlace {
		t.Errorf("Expected true with no existing buildings, got false")
	}
}

func TestCanPlace_Overlap(t *testing.T) {
	db := &MockDBExecutor{
		queryResults: []*MockRows{
			{rows: [][]any{
				{int64(1), 5, 5, 2, 2},
			}},
		},
	}

	canPlace, err := CanPlace(context.Background(), db, 1, 1, 1, 5, 5, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if canPlace {
		t.Errorf("Expected false for overlapping placement, got true")
	}
}

func TestCanPlace_NoOverlap(t *testing.T) {
	db := &MockDBExecutor{
		queryResults: []*MockRows{
			{rows: [][]any{
				{int64(1), 0, 0, 2, 2},
			}},
		},
	}

	canPlace, err := CanPlace(context.Background(), db, 1, 2, 2, 5, 5, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !canPlace {
		t.Errorf("Expected true for non-overlapping placement, got false")
	}
}

func TestCanPlace_SkipsOwnBuilding(t *testing.T) {
	db := &MockDBExecutor{
		queryResults: []*MockRows{
			{rows: [][]any{
				{int64(1), 5, 5, 2, 2},
			}},
		},
	}

	canPlace, err := CanPlace(context.Background(), db, 1, 2, 2, 5, 5, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !canPlace {
		t.Errorf("Expected true when skipping own building, got false")
	}
}

//

func completedUpgradesRow() *MockRow {
	return &MockRow{values: []any{0}}
}

func TestBuildingUpgrade_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("connection failed")},
		},
	}

	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err == nil {
		t.Errorf("Expected error from DB, got nil")
	}
}

func TestBuildingUpgrade_AlreadyUpgrading(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			completedUpgradesRow(),
			{values: []any{int64(2), 1, pgtype.Timestamp{Valid: true}}},
		},
	}
	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected error for already upgrading building")
	}
	if err.Error() != "Building already upgrading" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestBuildingUpgrade_AnotherUpgrading(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			completedUpgradesRow(),
			{values: []any{int64(2), 1, pgtype.Timestamp{Valid: false}}},
			{values: []any{1}},
		},
	}

	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected error for another building upgrading")
	}
	if err.Error() != "Another building is already upgrading" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestBuildingUpgrade_MaxTownhallLevel(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			completedUpgradesRow(),
			{values: []any{int64(1), 4, pgtype.Timestamp{Valid: false}}},
			{values: []any{0}},
			{values: []any{1, 500, 500}},
		},
	}

	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected error for max townhall level")
	}
	if err.Error() != "Townhall has reached its maximum level of 4" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestBuildingUpgrade_ExceedsTownhallLevel(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			completedUpgradesRow(),
			{values: []any{int64(2), 2, pgtype.Timestamp{Valid: false}}},
			{values: []any{0}},
			{values: []any{2, 500, 500}},
		},
	}

	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected error for exceeding townhall level")
	}
	if err.Error() != "building cannot exceed townhall level" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestBuildingUpgrade_NotEnoughGold(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			completedUpgradesRow(),
			{values: []any{int64(2), 1, pgtype.Timestamp{Valid: false}}},
			{values: []any{0}},
			{values: []any{4, 100, 500}},
			{values: []any{500, "Gold", 60}},
		},
	}

	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected error for not enough gold")
	}
	if err.Error() != "Not enough gold" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}
func TestBuildingUpgrade_NotEnoughElixir(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			completedUpgradesRow(),
			{values: []any{int64(2), 1, pgtype.Timestamp{Valid: false}}},
			{values: []any{0}},
			{values: []any{4, 500, 100}},
			{values: []any{500, "Elixir", 60}},
		},
	}

	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected error for not enough elixir")
	}
	if err.Error() != "Not enough elixir" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestBuildingUpgrade_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			completedUpgradesRow(),
			{values: []any{int64(2), 1, pgtype.Timestamp{Valid: false}}},
			{values: []any{0}},
			{values: []any{4, 1000, 500}},
			{values: []any{200, "Gold", 60}},
		},
	}

	err := BuildingUpgrade(context.Background(), db, 1, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

//

func TestCompleteUpgrades_NoCompletedUpgrades(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{0}},
		},
	}

	err := CompleteUpgrades(context.Background(), db, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCompleteUpgrades_TownhallCompleted(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1}},
		},
	}

	err := CompleteUpgrades(context.Background(), db, 1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCompleteUpgrades_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("connection failed")},
		},
	}

	err := CompleteUpgrades(context.Background(), db, 1)
	if err == nil {
		t.Errorf("Expected error from DB, got nil")
	}
}

func TestCompleteUpgrades_ExecError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1}},
		},
		execErr: errors.New("update failed"),
	}

	err := CompleteUpgrades(context.Background(), db, 1)
	if err == nil {
		t.Errorf("Expected error when Exec fails, got nil")
	}
}

//

func TestIsDefense_True(t *testing.T) {
	defenseIDs := []int64{2, 3, 4, 5}

	for _, id := range defenseIDs {
		if !isDefense(id) {
			t.Errorf("Expected building %d to be a defense building", id)
		}
	}
}

func TestIsDefense_False(t *testing.T) {
	nonDefenseIDs := []int64{1, 6, 7, 8, 9, 10}

	for _, id := range nonDefenseIDs {
		if isDefense(id) {
			t.Errorf("Expected building %d to not be a defense building", id)
		}
	}
}

//

func TestAddBuilding_InvalidPosition(t *testing.T) {
	db := &MockDBExecutor{}

	err := AddBuilding(context.Background(), db, 1, 2, -1, 5)
	if err == nil {
		t.Fatalf("Expected error for invalid position")
	}
	if err.Error() != "invalid position" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestAddBuilding_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("connection failed")},
		},
	}

	err := AddBuilding(context.Background(), db, 1, 2, 5, 5)
	if err == nil {
		t.Errorf("Expected error from DB, got nil")
	}
}

func TestAddBuilding_DefenseLocked(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1}},
			{values: []any{200, "Gold", 2, 2}},
			{values: []any{3}},
		},
	}

	err := AddBuilding(context.Background(), db, 1, 2, 5, 5)
	if err == nil {
		t.Fatalf("Expected error for locked building")
	}
	if err.Error() != "Building locked" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestAddBuilding_MaxQuantityReached(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{2}},
			{values: []any{150, "Elixir", 3, 3}},
			{values: []any{1}},
			{values: []any{1}},
		},
	}

	err := AddBuilding(context.Background(), db, 1, 6, 5, 5)
	if err == nil {
		t.Fatalf("Expected error for max quantity")
	}
	if err.Error() != "max quantity reached" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestAddBuilding_InvalidPlacement(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{2}},
			{values: []any{150, "Elixir", 3, 3}},
			{values: []any{2}},
			{values: []any{1}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{
				{int64(1), 5, 5, 3, 3},
			}},
		},
	}

	err := AddBuilding(context.Background(), db, 1, 6, 5, 5)
	if err == nil {
		t.Fatalf("Expected error for invalid placement")
	}
	if err.Error() != "invalid placement" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestAddBuilding_NotEnoughGold(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{2}},
			{values: []any{500, "Gold", 2, 2}},
			{values: []any{2}},
			{values: []any{0}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{}},
		},
	}

	err := AddBuilding(context.Background(), db, 1, 8, 5, 5)
	if err == nil {
		t.Fatalf("Expected error for not enough gold")
	}
	if err.Error() != "Not enough gold" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestAddBuilding_NotEnoughElixir(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{2}},
			{values: []any{500, "Elixir", 3, 3}},
			{values: []any{2}},
			{values: []any{0}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{}},
		},
	}

	err := AddBuilding(context.Background(), db, 1, 6, 5, 5)
	if err == nil {
		t.Fatalf("Expected error for not enough elixir")
	}
	if err.Error() != "Not enough elixir" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestAddBuilding_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{2}},
			{values: []any{200, "Gold", 3, 3}},
			{values: []any{2}},
			{values: []any{0}},
		},
		queryResults: []*MockRows{
			{rows: [][]any{}},
		},
		execRowsAffected: 1,
	}

	err := AddBuilding(context.Background(), db, 1, 8, 5, 5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
