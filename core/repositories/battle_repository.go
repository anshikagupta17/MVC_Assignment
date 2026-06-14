package repositories

import (
	"context"
	"errors"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
)

func (r *VillageRepository) FindOpponent(village_ID int64) (*models.MatchmakingResult, error) {
	ctx := context.Background()
	var townhallLevel int
	var trophies int
	err := r.DB.QueryRow(ctx,
		`SELECT townhall_level, trophies
		FROM villages
		WHERE id = $1`, village_ID).Scan(&townhallLevel, &trophies)

	if err != nil {
		return nil, err
	}

	var opponent models.MatchmakingResult

	err = r.DB.QueryRow(
		ctx,
		`SELECT
			id,
			townhall_level,
			trophies,
			gold,
			elixir
		FROM villages
		WHERE id != $1
		AND ABS(townhall_level - $2) <= 1
		AND ABS(trophies - $3) <= 100
		ORDER BY RANDOM()
		LIMIT 1`, village_ID, townhallLevel, trophies,
	).Scan(
		&opponent.VillageID,
		&opponent.TownhallLevel,
		&opponent.Trophies,
		&opponent.Gold,
		&opponent.Elixir,
	)

	if err != nil {
		return nil, err
	}

	return &opponent, nil
}

func (r *VillageRepository) ValidateArmy(village_id int64, deployedTroops []models.AttackTroop) error {
	ctx := context.Background()

	if len(deployedTroops) == 0 {
		return errors.New("No troops selected")
	}

	for _, troop := range deployedTroops {

		if troop.Quantity <= 0 {
			return errors.New("Invalid troop quantity")
		}

		var owned_amount int

		err := r.DB.QueryRow(
			ctx,
			`SELECT quantity
			FROM troops_village
			WHERE village_id = $1
			AND troops_id = $2`, village_id, troop.TroopID).Scan(&owned_amount)

		if err != nil {
			return errors.New("Troop not owned")
		}

		if troop.Quantity > owned_amount {
			return errors.New("Not enough troops")
		}
	}

	return nil
}

func (r *VillageRepository) CalculateAttackPower(village_id int64, deployedTroops []models.AttackTroop) (int, error)

func (r *VillageRepository) CalculateDefensePower(village_id int64) (int, error)

func (r *VillageRepository) CalculateDestruction(defender_village_id int64, attack_power int) (int, bool, error)

func CalculateBattleResult(destruction int, townhall_destroyed bool, defender_gold, defender_elixir int) (stars, loot_gold, loot_elixir, trophy_change int)
