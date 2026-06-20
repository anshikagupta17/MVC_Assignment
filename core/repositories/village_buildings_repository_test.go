package repositories

import "testing"

func TestTrueCases(t *testing.T) {
	defenseIDs := []int64{2, 3, 4, 5}

	for _, id := range defenseIDs {
		if !isDefense(id) {
			t.Errorf("Expected building_id %d to be a defense building", id)
		}
	}
}

func TestIsDefense_FalseCases(t *testing.T) {
	nonDefenseIDs := []int64{1, 6, 7, 8, 9, 10}

	for _, id := range nonDefenseIDs {
		if isDefense(id) {
			t.Errorf("Expected building_id %d to NOT be a defense building", id)
		}
	}
}
