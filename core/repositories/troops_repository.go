package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
)

func (r *VillageRepository) TrainTroops(village_id int64, troop_id int64, quantity int) error {
	ctx := context.Background()

	var townhall_level int
	var elixir int
	var housing_space int

	err := r.DB.QueryRow(ctx,
		`SELECT townhall_level, elixir, housing_space
		FROM villages
		WHERE id = $1`, village_id).Scan(&townhall_level, &elixir, &housing_space)

	if err != nil {
		return err
	}

	var unlock_level int
	var training_cost int
	var housing_one int

	err = r.DB.QueryRow(ctx,
		`SELECT unlock_level, training_cost, housing_space
		FROM troops_base_metadata
		WHERE id = $1`, troop_id).Scan(&unlock_level, &training_cost, &housing_one)

	if err != nil {
		return err
	}

	if townhall_level < unlock_level {
		return errors.New("Troop locked")
	}

	var used_space int

	err = r.DB.QueryRow(ctx,
		`SELECT COALESCE(SUM(tv.quantity * tb.housing_space),0)
		FROM troops_village tv
		JOIN troops_base_metadata tb
		ON tb.id = tv.troops_id
		WHERE tv.village_id = $1`, village_id).Scan(&used_space)

	if err != nil {
		return err
	}

	required := housing_one * quantity

	if used_space+required > housing_space {
		return errors.New("Not enough housing space")
	}
	totalCost := training_cost * quantity

	if elixir < totalCost {
		return errors.New("Not enough Elixir")
	}
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE villages
		SET elixir = elixir - $1
		WHERE id = $2`, totalCost, village_id)

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO troops_village
		(village_id, troops_id, level, quantity)
		VALUES ($1,$2,1,$3)
		ON CONFLICT (village_id, troops_id)
		DO UPDATE SET quantity = troops_village.quantity + EXCLUDED.quantity`, village_id, troop_id, quantity)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)

}

func (r *VillageRepository) GetVillageTroops(village_id int64) ([]models.VillageTroop, error) {
	err := r.CompleteTroopUpgrades(village_id)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	rows, err := r.DB.Query(ctx,
		`SELECT
			tv.troops_id,
			tb.name,
			tv.level,
			tv.quantity,
			tlm.damage,
			tlm.max_health
		FROM troops_village tv
		JOIN troops_base_metadata tb
			ON tb.id = tv.troops_id
		JOIN troops_level_metadata tlm
			ON tlm.type_id = tv.troops_id AND tlm.level = tv.level
		WHERE tv.village_id = $1`,
		village_id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var troops []models.VillageTroop

	for rows.Next() {
		var t models.VillageTroop

		err := rows.Scan(
			&t.TroopID,
			&t.Name,
			&t.Level,
			&t.Quantity,
			&t.Damage,
			&t.MaxHealth,
		)

		if err != nil {
			return nil, err
		}

		troops = append(troops, t)
	}

	return troops, nil
}

func (r *VillageRepository) UpgradeTroops(village_id int64, troops_id int) error {
	ctx := context.Background()
	var level, quantity int
	var upgrade_ends_at *time.Time
	err := r.DB.QueryRow(ctx,
		`SELECT level, quantity, upgrade_ends_at
	FROM troops_village
	WHERE id=$1 AND troops_id=$2`,
		village_id, troops_id).Scan(&level, &quantity, &upgrade_ends_at)
	if err != nil {
		return err
	}

	if upgrade_ends_at != nil {
		return errors.New("Troops already upgrading")
	}
	var townhall_level, elixir int
	err = r.DB.QueryRow(ctx,
		`SELECT townhall_level, elixir
	FROM villages
	WHERE village_id=$1`, village_id).Scan(&townhall_level, &elixir)

	if level >= townhall_level {
		return errors.New("Troop at max upgrade level")
	}

	var cost int64
	err = r.DB.QueryRow(ctx,
		`SELECT upgrade_cost from troops_level_metadata
	WHERE level= level+1 AND type_id=$1,`, troops_id).Scan(&cost)
	if err != nil {
		return err
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var troop_cost int64 = cost
	if int64(elixir) < (troop_cost) {
		return errors.New("Not enough elixir")
	}

	_, err = tx.Exec(ctx,
		`UPDATE villages
			SET elixir = elixir - $1
			WHERE id = $2`, troop_cost, village_id)

	if err != nil {
		return err
	}
	var upgradeTimeSec int
	err = r.DB.QueryRow(ctx,
		`SELECT upgrade_time_sec from troops_base_metadata
	WHERE level= $1+1 AND type_id=$2`, level, troops_id).Scan(&upgradeTimeSec)
	if err != nil {
		return err
	}

	upgradeEnd := time.Now().Add(
		time.Duration(upgradeTimeSec) * time.Second,
	)

	_, err = tx.Exec(
		ctx,
		`UPDATE troops_village
		SET upgrade_ends_at = $1
		WHERE village_id = $2
		AND troops_id = $3`,
		upgradeEnd,
		village_id,
		troops_id,
	)

	if err != nil {
		return err
	}
	return tx.Commit(ctx)

}

func (r *VillageRepository) CompleteTroopUpgrades(village_id int64) error {
	ctx := context.Background()

	rows, err := r.DB.Query(ctx,
		`SELECT troops_id
		FROM troops_village
		WHERE village_id = $1
		AND upgrade_ends_at IS NOT NULL
		AND upgrade_ends_at <= NOW()`, village_id)

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {

		var troop_id int64

		err := rows.Scan(&troop_id)
		if err != nil {
			return err
		}

		_, err = r.DB.Exec(ctx,
			`UPDATE troops_village
			SET level = level + 1,
				upgrade_ends_at = NULL
			WHERE village_id = $1
			AND troops_id = $2`, village_id, troop_id)

		if err != nil {
			return err
		}
	}

	return nil
}
