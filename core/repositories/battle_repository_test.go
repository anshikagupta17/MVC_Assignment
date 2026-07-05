package repositories

import (
	"context"
	"errors"
	"testing"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
)

func TestZeroDefense(t *testing.T) {
	result := CalculateResult(100, 100, 0, 0)

	if result.Stars != 3 {
		t.Errorf("Expected 3 stars, got %d", result.Stars)
	}

	if result.DestructionPercentage != 100 {
		t.Errorf("Expected 100%% destruction, got %d", result.DestructionPercentage)
	}

	if !result.TownhallDestroyed {
		t.Errorf("Expected townhall to be destroyed when defense is zero")
	}
}

func TestLowDestruction(t *testing.T) {
	result := CalculateResult(10, 10, 1000, 1000)

	if result.Stars != 0 {
		t.Errorf("Expected 0 stars, got %d", result.Stars)
	}

	if result.TrophyChange != -10 {
		t.Errorf("Expected -10 trophy change, got %d", result.TrophyChange)
	}
}

func TestFullDestruction(t *testing.T) {
	result := CalculateResult(500, 500, 500, 500)

	if result.DestructionPercentage != 100 {
		t.Errorf("Expected 100%% destruction, got %d", result.DestructionPercentage)
	}

	if result.Stars != 3 {
		t.Errorf("Expected 3 stars, got %d", result.Stars)
	}
}

//

func TestValidateArmy_NoTroops(t *testing.T) {
	db := &MockDBExecutor{}

	err := ValidateArmy(context.Background(), db, 1, []models.AttackTroop{})
	if err == nil {
		t.Fatalf("Expected error for empty army")
	}
	if err.Error() != "No troops selected" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestValidateArmy_InvalidQuantity(t *testing.T) {
	db := &MockDBExecutor{}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 0},
	}

	err := ValidateArmy(context.Background(), db, 1, troops)
	if err == nil {
		t.Fatalf("Expected error for zero quantity")
	}
	if err.Error() != "Invalid troop quantity" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestValidateArmy_TroopNotOwned(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("no rows")},
		},
	}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 5},
	}

	err := ValidateArmy(context.Background(), db, 1, troops)
	if err == nil {
		t.Fatalf("Expected error for troop not owned")
	}
	if err.Error() != "Troop not owned" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestValidateArmy_NotEnoughTroops(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{3}},
		},
	}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 10},
	}

	err := ValidateArmy(context.Background(), db, 1, troops)
	if err == nil {
		t.Fatalf("Expected error for not enough troops")
	}
	if err.Error() != "Not enough troops" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestValidateArmy_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{10}},
		},
	}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 5},
	}

	err := ValidateArmy(context.Background(), db, 1, troops)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

//

func TestCalculateAttackPower_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("connection failed")},
		},
	}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 2},
	}

	_, _, err := CalculateAttackPower(context.Background(), db, 1, troops)
	if err == nil {
		t.Errorf("Expected error from DB, got nil")
	}
}

func TestCalculateAttackPower_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{10, 5}},
		},
	}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 2},
	}

	power, health, err := CalculateAttackPower(context.Background(), db, 1, troops)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if power != 20 {
		t.Errorf("Expected power 20, got %d", power)
	}
	if health != 10 {
		t.Errorf("Expected health 10, got %d", health)
	}
}

//

func TestCalculateDefensePower_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryErr: errors.New("connection failed"),
	}

	_, _, err := CalculateDefensePower(context.Background(), db, 1)
	if err == nil {
		t.Errorf("Expected error from DB, got nil")
	}
}

func TestCalculateDefensePower_NoDefenses(t *testing.T) {
	db := &MockDBExecutor{
		queryResults: []*MockRows{
			{rows: [][]any{}},
		},
	}

	damage, health, err := CalculateDefensePower(context.Background(), db, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if damage != 0 {
		t.Errorf("Expected 0 damage, got %d", damage)
	}
	if health != 0 {
		t.Errorf("Expected 0 health, got %d", health)
	}
}

func TestCalculateDefensePower_WithDefenses(t *testing.T) {
	db := &MockDBExecutor{
		queryResults: []*MockRows{
			{rows: [][]any{
				{50, 100},
				{30, 80},
			}},
		},
	}

	damage, health, err := CalculateDefensePower(context.Background(), db, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if damage != 80 {
		t.Errorf("Expected damage 80, got %d", damage)
	}
	if health != 180 {
		t.Errorf("Expected health 180, got %d", health)
	}
}

//

func TestLoot_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("connection failed")},
		},
	}

	_, _, err := Loot(context.Background(), db, 1, 100)
	if err == nil {
		t.Errorf("Expected error from DB, got nil")
	}
}

func TestLoot_LowDestruction(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1000, 1000}},
		},
	}

	gold, elixir, err := Loot(context.Background(), db, 1, 20)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if gold != 0 {
		t.Errorf("Expected 0 gold for low destruction, got %d", gold)
	}
	if elixir != 0 {
		t.Errorf("Expected 0 elixir for low destruction, got %d", elixir)
	}
}

func TestLoot_FullDestruction(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1000, 1000}},
		},
	}

	gold, elixir, err := Loot(context.Background(), db, 1, 100)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if gold != 400 {
		t.Errorf("Expected 400 gold, got %d", gold)
	}
	if elixir != 400 {
		t.Errorf("Expected 400 elixir, got %d", elixir)
	}
}

//

func TestConsumeTroops_TroopNotFound(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("no rows")},
		},
	}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 3},
	}

	err := ConsumeTroops(context.Background(), db, 1, troops)
	if err == nil {
		t.Fatalf("Expected error for troop not found")
	}
	if err.Error() != "Troop not found" {
		t.Errorf("Wrong error: %s", err.Error())
	}
}

func TestConsumeTroops_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{10}},
		},
	}

	troops := []models.AttackTroop{
		{TroopID: 1, Quantity: 3},
	}

	err := ConsumeTroops(context.Background(), db, 1, troops)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

//

func TestApplyBattleResult_DBError(t *testing.T) {
	db := &MockDBExecutor{
		execErr: errors.New("connection failed"),
	}

	result := models.BattleResult{
		LootGold:   100,
		LootElixir: 50,
	}

	err := ApplyBattleResult(context.Background(), db, 1, 2, result)
	if err == nil {
		t.Errorf("Expected error from DB, got nil")
	}
}

func TestApplyBattleResult_Success(t *testing.T) {
	db := &MockDBExecutor{}

	result := models.BattleResult{
		LootGold:              100,
		LootElixir:            50,
		TrophyChange:          10,
		Stars:                 2,
		DestructionPercentage: 80,
	}

	err := ApplyBattleResult(context.Background(), db, 1, 2, result)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
