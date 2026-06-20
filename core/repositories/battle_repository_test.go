package repositories

import "testing"

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
