package repositories

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestTrainTroops_TroopLocked(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1, 500}},
			{values: []any{20}},
			{values: []any{3, 50, 1}},
		},
	}
	err := TrainTroops(context.Background(), db, 1, 3, 1)
	if err == nil {
		t.Fatalf("Expected an error but got nothing")
	}
	if err.Error() != "Troop locked" {
		t.Errorf("Wrong error message: got %s", err.Error())
	}
}

func TestTrainTroops_NoHousing(t *testing.T) {
	mockData := []*MockRow{
		{values: []any{1, 500}},
		{values: []any{20}},
		{values: []any{1, 50, 1}},
		{values: []any{18}},
	}
	db := &MockDBExecutor{queryRows: mockData}
	err := TrainTroops(context.Background(), db, 1, 1, 5)
	if err == nil {
		t.Fatalf("Should have failed for housing space")
	}
	if err.Error() != "Not enough housing space" {
		t.Errorf("Got bad error: %v", err)
	}
}

func TestTrainTroops_NoElixir(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1, 100}},
			{values: []any{20}},
			{values: []any{1, 50, 1}},
			{values: []any{0}},
		},
	}
	err := TrainTroops(context.Background(), db, 1, 1, 5)
	if err == nil {
		t.Fatalf("Expected elixir error")
	}
	if err.Error() != "Not enough Elixir" {
		t.Errorf("Expected Not enough Elixir, got %s", err.Error())
	}
}

func TestTrainTroops_Works(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1, 500}},
			{values: []any{20}},
			{values: []any{1, 50, 1}},
			{values: []any{0}},
		},
	}
	err := TrainTroops(context.Background(), db, 1, 1, 2)
	if err != nil {
		t.Errorf("It should work but got error: %s", err.Error())
	}
}

func TestTrainTroops_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("sql connection lost")},
		},
	}
	err := TrainTroops(context.Background(), db, 1, 1, 1)
	if err == nil {
		t.Errorf("Expected error from database")
	}
}

func TestUpgradeTroops_AlreadyUpgrading(t *testing.T) {
	timer := time.Now().Add(60 * time.Second)
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1, 10, &timer}},
		},
	}
	err := UpgradeTroops(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Should fail because already upgrading")
	}
	if err.Error() != "Troops already upgrading" {
		t.Errorf("Got: %s", err.Error())
	}
}

func TestUpgradeTroops_MaxLevel(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{2, 10, (*time.Time)(nil)}},
			{values: []any{2, 500}},
		},
	}
	err := UpgradeTroops(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected max level error")
	}
	if err.Error() != "Troop at max upgrade level" {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestUpgradeTroops_NoElixir(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1, 10, (*time.Time)(nil)}},
			{values: []any{2, 100}},
			{values: []any{int64(500)}},
		},
	}
	err := UpgradeTroops(context.Background(), db, 1, 1)
	if err == nil {
		t.Fatalf("Expected error")
	}
	if err.Error() != "Not enough elixir" {
		t.Errorf("Expected not enough elixir, got %s", err.Error())
	}
}

func TestUpgradeTroops_Success(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{values: []any{1, 10, (*time.Time)(nil)}},
			{values: []any{2, 1000}},
			{values: []any{int64(500)}},
			{values: []any{1080}},
		},
	}
	err := UpgradeTroops(context.Background(), db, 1, 1)
	if err != nil {
		t.Errorf("Upgrade failed: %s", err.Error())
	}
}

func TestUpgradeTroops_DBError(t *testing.T) {
	db := &MockDBExecutor{
		queryRows: []*MockRow{
			{err: errors.New("broken pipe")},
		},
	}
	err := UpgradeTroops(context.Background(), db, 1, 1)
	if err == nil {
		t.Errorf("Should have failed")
	}
}
