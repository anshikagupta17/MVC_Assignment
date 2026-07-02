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

func ValidateArmy(ctx context.Context, db DBExecutor, village_id int64, deployedTroops []models.AttackTroop) error {

	if len(deployedTroops) == 0 {
		return errors.New("No troops selected")
	}

	for _, troop := range deployedTroops {

		if troop.Quantity <= 0 {
			return errors.New("Invalid troop quantity")
		}

		var owned_amount int

		err := db.QueryRow(
			ctx,
			`SELECT quantity
			FROM troops_village
			WHERE village_id = $1
			AND troops_id = $2
			FOR UPDATE`, village_id, troop.TroopID).Scan(&owned_amount)

		if err != nil {
			return errors.New("Troop not owned")
		}

		if troop.Quantity > owned_amount {
			return errors.New("Not enough troops")
		}
	}

	return nil
}

func CalculateAttackPower(ctx context.Context, db DBExecutor, village_id int64, deployed_troops []models.AttackTroop) (int, int, error) {
	total_power := 0
	total_health := 0
	for _, troop := range deployed_troops {
		var damage int
		var health int

		err := db.QueryRow(ctx,
			`SELECT tlm.damage, tlm.max_health
			FROM troops_level_metadata tlm
			JOIN troops_village tv 
			ON tlm.type_id = tv.troops_id AND tlm.level=tv.level
			WHERE tv.village_id = $1
			AND tv.troops_id = $2`, village_id, troop.TroopID).Scan(&damage, &health)

		if err != nil {
			return 0, 0, err
		}

		total_power += damage * troop.Quantity
		total_health += health * troop.Quantity
	}

	return total_power, total_health, nil
}

func CalculateDefensePower(ctx context.Context, db DBExecutor, village_id int64) (int, int, error) {
	total_damage := 0
	total_health := 0
	rows, err := db.Query(ctx,
		`SELECT
			dm.damage,
			dm.max_health
		FROM defense_metadata dm 
		JOIN buildings_village bv
			ON bv.building_id = dm.type_id
			AND bv.level = dm.level
		WHERE bv.village_id = $1`, village_id)

	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {

		var damage int
		var health int

		err := rows.Scan(
			&damage,
			&health,
		)

		if err != nil {
			return 0, 0, err
		}

		total_damage += damage
		total_health += health
	}

	if err := rows.Err(); err != nil {
		return 0, 0, err
	}

	return total_damage, total_health, nil
}

func CalculateResult(AttackerDamage int, AttackerHealth int, DefenderDamage int, DefenderHealth int) models.BattleResult {

	AttackPower := AttackerDamage + AttackerHealth
	DefensePower := DefenderDamage + DefenderHealth

	if DefensePower == 0 {
		return models.BattleResult{
			Stars:                 3,
			DestructionPercentage: 100,
			TownhallDestroyed:     true,
			TrophyChange:          30,
		}
	}

	destruction := (AttackPower * 100 / DefensePower)
	destruction = min(destruction, 100)

	result := models.BattleResult{
		DestructionPercentage: destruction,
	}

	switch {

	case destruction < 50:
		result.Stars = 0
		result.TownhallDestroyed = false
		result.TrophyChange = -10

	case destruction < 75:
		result.Stars = 1
		result.TownhallDestroyed = false
		result.TrophyChange = 10

	case destruction < 100:
		result.Stars = 2
		result.TownhallDestroyed = true
		result.TrophyChange = 20

	default:
		result.Stars = 3
		result.TownhallDestroyed = true
		result.TrophyChange = 30
	}

	return result
}

func Loot(ctx context.Context, db DBExecutor, DefenderVillageID int64, destruction_percent int) (int, int, error) {
	var defender_gold, defender_elixir int
	err := db.QueryRow(ctx,
		`SELECT gold, elixir
		FROM villages
		WHERE id=$1
		FOR UPDATE`, DefenderVillageID).Scan(&defender_gold, &defender_elixir)
	if err != nil {
		return 0, 0, err
	}
	var loot_elixir, loot_gold int
	switch {
	case destruction_percent < 30:
		loot_gold = 0
		loot_elixir = 0
	case destruction_percent < 50:
		loot_gold = (defender_gold * 10) / 100
		loot_elixir = (defender_elixir * 10) / 100
	case destruction_percent < 75:
		loot_gold = (defender_gold * 20) / 100
		loot_elixir = (defender_elixir * 20) / 100
	case destruction_percent < 100:
		loot_gold = (defender_gold * 30) / 100
		loot_elixir = (defender_elixir * 30) / 100
	default:
		loot_gold = (defender_gold * 40) / 100
		loot_elixir = (defender_elixir * 40) / 100
	}

	return loot_gold, loot_elixir, nil
}
func ConsumeTroops(ctx context.Context, db DBExecutor, village_id int64, deployed_troops []models.AttackTroop) error {
	for _, troops := range deployed_troops {
		var Quantity int
		err := db.QueryRow(ctx,
			`SELECT quantity
			FROM troops_village
			WHERE village_id=$1 AND troops_id=$2
			FOR UPDATE`, village_id, troops.TroopID).Scan(&Quantity)

		if err != nil {
			return errors.New("Troop not found")
		}
		if Quantity > 0 {
			Quantity -= troops.Quantity
			Quantity = max(0, Quantity)
		}
		_, err = db.Exec(ctx,
			`UPDATE troops_village
		SET quantity=$1
		Where village_id=$2 AND troops_id=$3`, Quantity, village_id, troops.TroopID)

		if err != nil {
			return err
		}
	}
	return nil

}

func ApplyBattleResult(ctx context.Context, db DBExecutor, AttackerVillageID int64, DefenderVillageID int64, result models.BattleResult) error {
	_, err := db.Exec(ctx,
		`UPDATE villages
		SET
			gold = gold + $1,
			elixir = elixir + $2,
			trophies = GREATEST(0,trophies + $3)
		WHERE id = $4`, result.LootGold, result.LootElixir, result.TrophyChange, AttackerVillageID)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx,
		`UPDATE villages
		SET
			gold = GREATEST(0,gold - $1),
			elixir = GREATEST(0,elixir - $2),
			trophies = GREATEST(0, trophies - $3)
		WHERE id = $4`, result.LootGold, result.LootElixir, result.TrophyChange, DefenderVillageID)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx,
		`INSERT INTO battles(
			attacker_id,
			defender_id,
			loot_gold,
			loot_elixir,
			start_time,
			stars,
			destruction_percent,
			trophy_change
		)
		VALUES(
			$1,$2,$3,$4,NOW(),$5,$6,$7
		)
		`,
		AttackerVillageID,
		DefenderVillageID,
		result.LootGold,
		result.LootElixir,
		result.Stars,
		result.DestructionPercentage,
		result.TrophyChange,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *VillageRepository) AttackVillage(AttackerVillageID int64, req models.AttackRequest) (models.BattleResult, error) {
	ctx := context.Background()
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return models.BattleResult{}, err
	}
	defer tx.Rollback(ctx)

	err = ValidateArmy(ctx, tx, AttackerVillageID, req.Troops)

	if err != nil {
		return models.BattleResult{}, err
	}
	AttackerDamage, AttackerHealth, err := CalculateAttackPower(ctx, tx, AttackerVillageID, req.Troops)
	if err != nil {
		return models.BattleResult{}, err
	}

	DefenderDamage, DefenderHealth, err := CalculateDefensePower(ctx, tx, req.DefenderID)
	if err != nil {
		return models.BattleResult{}, err
	}

	result := CalculateResult(AttackerDamage, AttackerHealth, DefenderDamage, DefenderHealth)

	LootGold, LootElixir, err := Loot(ctx, tx, req.DefenderID, result.DestructionPercentage)
	if err != nil {
		return models.BattleResult{}, err
	}

	result.LootGold = LootGold
	result.LootElixir = LootElixir

	err = ConsumeTroops(ctx, tx, AttackerVillageID, req.Troops)
	if err != nil {
		return models.BattleResult{}, err
	}

	err = ApplyBattleResult(ctx, tx, AttackerVillageID, req.DefenderID, result)
	if err != nil {
		return models.BattleResult{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return models.BattleResult{}, err
	}

	return result, nil
}

func (r *VillageRepository) GetOpponentBuildings(village_id int64) ([]models.VillageBuilding, error) {
	ctx := context.Background()

	rows, err := r.DB.Query(ctx,
		`SELECT bv.id, bv.building_id, bv.level, bv.x, bv.y, bm.size_x, bm.size_y
        FROM buildings_village bv
        JOIN buildings_metadata bm ON bm.id = bv.building_id
        WHERE bv.village_id = $1`, village_id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.VillageBuilding

	for rows.Next() {
		var b models.VillageBuilding

		err := rows.Scan(
			&b.ID,
			&b.BuildingId,
			&b.Level,
			&b.X,
			&b.Y,
			&b.SizeX,
			&b.SizeY,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
