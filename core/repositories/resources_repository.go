package repositories

import (
	"context"
	"time"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type CollectedResources struct {
	Gold   int
	Elixir int
}

func (r *VillageRepository) CollectResources(village_id int64) (CollectedResources, error) {
	ctx := context.Background()
	var gold, elixir int
	err := r.DB.QueryRow(ctx,
		`SELECT gold, elixir
		FROM villages
		WHERE id = $1`, village_id).Scan(&gold, &elixir)
	if err != nil {
		return CollectedResources{}, err
	}
	rows, err := r.DB.Query(ctx,
		`SELECT b.id, b.level, b.building_id, b.last_collected_at, b.upgrade_ends_at,
		m.type_id, m.production_rate
		FROM buildings_village b
		JOIN mines_metadata m ON m.type_id=b.building_id AND m.level=b.level
		WHERE b.village_id = $1 AND b.building_id IN (6,7)`, village_id)

	if err != nil {
		return CollectedResources{}, err
	}
	defer rows.Close()

	var mines []models.Mine
	now := time.Now()

	for rows.Next() {
		var m models.Mine

		var lastCollected pgtype.Timestamp

		err := rows.Scan(
			&m.ID,
			&m.Level,
			&m.BuildingId,
			&lastCollected,
			&m.UpgradeEndsAt,
			&m.TypeId,
			&m.ProductionRate,
		)
		if err != nil {
			return CollectedResources{}, err
		}
		if lastCollected.Valid {
			m.LastCollected = lastCollected.Time
		} else {
			m.LastCollected = now
		}
		mines = append(mines, m)

	}

	new_gold := 0
	new_elixir := 0
	for _, m := range mines {
		if m.UpgradeEndsAt != nil {
			continue
		}

		sec := int(now.Sub(m.LastCollected).Seconds())
		if sec <= 0 {
			continue
		}

		gain := int(m.ProductionRate * float64(sec))
		switch m.TypeId {
		case 6:
			new_gold += gain

		case 7:
			new_elixir += gain
		}
	}
	s_rows, err := r.DB.Query(ctx,
		`SELECT b.building_id,
		s.max_capacity
		FROM buildings_village b
		JOIN storage_metadata s ON s.type_id = b.building_id AND s.level = b.level
		WHERE village_id = $1 AND b.building_id IN (8,9)`, village_id)

	if err != nil {
		return CollectedResources{}, err
	}
	defer s_rows.Close()
	gold_cap := 0
	elixir_cap := 0

	for s_rows.Next() {

		var building_id int64
		var capacity int

		err := s_rows.Scan(
			&building_id,
			&capacity,
		)

		if err != nil {
			return CollectedResources{}, err
		}

		switch building_id {

		case 8:
			gold_cap += capacity

		case 9:
			elixir_cap += capacity
		}
	}
	if gold+new_gold > gold_cap {
		new_gold = gold_cap - gold
		if new_gold < 0 {
			new_gold = 0

		}
	}

	if elixir+new_elixir > elixir_cap {
		new_elixir = elixir_cap - elixir
		if new_elixir < 0 {
			new_elixir = 0
		}
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return CollectedResources{}, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE villages
		SET gold = gold + $1,
			elixir = elixir + $2
		WHERE id = $3`, new_gold, new_elixir, village_id)

	if err != nil {
		return CollectedResources{}, err
	}

	_, err = tx.Exec(ctx,
		`UPDATE buildings_village
		SET last_collected_at = $1
		WHERE village_id = $2 AND building_id IN (6,7) AND upgrade_ends_at IS NULL
		`, now, village_id)

	if err != nil {
		return CollectedResources{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return CollectedResources{}, err
	}

	return CollectedResources{
		Gold:   new_gold,
		Elixir: new_elixir}, nil
}
