package repositories

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type BuildingMetadata struct {
	UpgradeCost    int
	CostType       string
	UpgradeTimeSec int
}

type VillageResources struct {
	Gold   int
	Elixir int
}

func (r *VillageRepository) VillageBuildings(village_id int64) ([]models.VillageBuilding, error) {
	err := r.CompleteUpgrades(village_id)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	rows, err := r.DB.Query(ctx,
		`SELECT id, building_id, level, upgrade_ends_at, x, y
		From buildings_village
		WHERE village_id= $1`, village_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.VillageBuilding
	for rows.Next() {
		var building models.VillageBuilding
		var upgradeEndsAt pgtype.Timestamp

		err := rows.Scan(
			&building.ID,
			&building.BuildingId,
			&building.Level,
			&upgradeEndsAt,
			&building.X,
			&building.Y,
		)

		if err != nil {
			return nil, err
		}

		if upgradeEndsAt.Valid {
			t := upgradeEndsAt.Time
			building.UpgradeEndsAt = &t
		}

		result = append(result, building)
	}
	return result, nil
}
func (r *VillageRepository) InitialBuildings(village_id int64) error {
	ctx := context.Background()
	_, err := r.DB.Exec(ctx,
		`INSERT INTO buildings_village
		(village_id, building_id, level, x, y)
		VALUES
		($1, 1, 1, 10, 10),
		($1, 6, 1, 5, 5),
		($1, 7, 1, 15, 5),
		($1, 10, 1, 10, 15)`, village_id)

	return err
}

func (r *VillageRepository) MoveBuilding(village_id, building_instance_id int64, x, y int) error {
	ctx := context.Background()
	_, err := r.DB.Exec(
		ctx,
		`UPDATE buildings_village
		SET x = $3,
			y = $4
		WHERE village_id = $1
		AND id = $2`,
		village_id,
		building_instance_id,
		x,
		y,
	)

	return err
}

func (r *VillageRepository) CanPlaceBuilding(village_id, building_instance_id int64, x, y int) (bool, error) {
	ctx := context.Background()
	var buildingID int64

	err := r.DB.QueryRow(
		ctx,
		`SELECT building_id
		FROM buildings_village
		WHERE id = $1`, building_instance_id).Scan(&buildingID)

	if err != nil {
		return false, err
	}

	var sizex, sizey int

	err = r.DB.QueryRow(
		ctx,
		`SELECT size_x, size_y
		FROM buildings_metadata
		WHERE id = $1`, buildingID).Scan(&sizex, &sizey)

	if err != nil {
		return false, err
	}

	rows, err := r.DB.Query(ctx,
		`SELECT id, building_id, x, y
		FROM buildings_village
		WHERE village_id = $1`, village_id)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	newX1 := x
	newY1 := y
	newX2 := x + sizex - 1
	newY2 := y + sizey - 1

	for rows.Next() {

		var instance_id int64
		var b_id int
		var current_x, current_y int

		err := rows.Scan(
			&instance_id,
			&b_id,
			&current_x,
			&current_y,
		)
		if err != nil {
			return false, err
		}
		if instance_id == building_instance_id {
			continue
		}

		var current_sizex, current_sizey int
		err = r.DB.QueryRow(ctx,
			`SELECT size_x, size_y 
			FROM buildings_metadata 
			WHERE id=$1`, b_id).Scan(&current_sizex, &current_sizey)

		if err != nil {
			return false, err
		}

		currentX1 := current_x
		currentX2 := current_x + current_sizex - 1
		currentY1 := current_y
		currentY2 := current_y + current_sizey - 1

		if !(newX2 < currentX1 || newX1 > currentX2 || newY2 < currentY1 || newY1 > currentY2) {
			return false, nil
		}
	}

	return true, nil
}

func (r *VillageRepository) CanPlaceNewBuilding(village_id int64, building_id int64, x int, y int) (bool, error) {
	ctx := context.Background()
	var sizex, sizey int

	err := r.DB.QueryRow(ctx,
		`SELECT size_x, size_y
		FROM buildings_metadata
		WHERE id = $1`, building_id).Scan(&sizex, &sizey)

	if err != nil {
		return false, err
	}

	rows, err := r.DB.Query(
		ctx,
		`SELECT id, building_id, x, y
		FROM buildings_village
		WHERE village_id = $1`, village_id,
	)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	newX1 := x
	newY1 := y
	newX2 := x + sizex - 1
	newY2 := y + sizey - 1

	for rows.Next() {

		var instance_id int64
		var CurrentBuildingId int64
		var current_x, current_y int

		err := rows.Scan(
			&instance_id,
			&CurrentBuildingId,
			&current_x,
			&current_y,
		)

		if err != nil {
			return false, err
		}

		var current_size_x, current_size_y int

		err = r.DB.QueryRow(ctx,
			`SELECT size_x, size_y
			FROM buildings_metadata
			WHERE id = $1`, CurrentBuildingId).Scan(&current_size_x, &current_size_y)

		if err != nil {
			return false, err
		}

		currentX1 := current_x
		currentY1 := current_y
		currentX2 := current_x + current_size_x - 1
		currentY2 := current_y + current_size_y - 1

		if !(newX2 < currentX1 || newX1 > currentX2 || newY2 < currentY1 || newY1 > currentY2) {
			return false, nil
		}
	}

	return true, nil

}

func (r *VillageRepository) BuildingUpgrade(village_id int64, building_instance_id int64) error {
	ctx := context.Background()

	var building_id int64
	var level int
	var upgrade_ends_at pgtype.Timestamp

	err := r.DB.QueryRow(ctx,
		`SELECT building_id, level, upgrade_ends_at
		FROM buildings_village
		WHERE id = $1
		AND village_id = $2`, building_instance_id, village_id).Scan(&building_id, &level, &upgrade_ends_at)

	if err != nil {
		return err
	}

	if upgrade_ends_at.Valid {
		return errors.New("Building already upgrading")
	}

	var upgradingCount int
	err = r.DB.QueryRow(ctx,
		`SELECT COUNT(*)
    	FROM buildings_village
    	WHERE village_id = $1
    	AND upgrade_ends_at IS NOT NULL`, village_id).Scan(&upgradingCount)

	if err != nil {
		return err
	}

	if upgradingCount > 0 {
		return errors.New("Another building is already upgrading")
	}

	var townhall_level int
	err = r.DB.QueryRow(ctx,
		`SELECT townhall_level
		FROM villages
		WHERE id = $1`, village_id).Scan(&townhall_level)

	log.Println("village_id:", village_id, "building_instance_id:", building_instance_id)

	if err != nil {
		return err
	}

	if level+1 > townhall_level {
		return errors.New("building cannot exceed townhall level")
	}

	var metadata BuildingMetadata

	err = r.DB.QueryRow(ctx,
		`SELECT upgrade_cost, cost_type, upgrade_time_sec
		FROM buildings_metadata
		WHERE id = $1`,
		building_id).Scan(&metadata.UpgradeCost, &metadata.CostType, &metadata.UpgradeTimeSec)

	if err != nil {
		return err
	}

	var resources VillageResources

	err = r.DB.QueryRow(
		ctx,
		`SELECT gold, elixir
		FROM villages
		WHERE id = $1`, village_id).Scan(&resources.Gold, &resources.Elixir)

	if err != nil {
		return err
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	switch metadata.CostType {

	case "Gold":

		if resources.Gold < metadata.UpgradeCost {
			return errors.New("Not enough gold")
		}

		_, err = tx.Exec(ctx,
			`UPDATE villages
			SET gold = gold - $1
			WHERE id = $2`, metadata.UpgradeCost, village_id)

	case "Elixir":

		if resources.Elixir < metadata.UpgradeCost {
			return errors.New("Not enough elixir")
		}

		_, err = tx.Exec(ctx,
			`UPDATE villages
			SET elixir = elixir - $1
			WHERE id = $2`, metadata.UpgradeCost, village_id)

	default:
		return errors.New("Invalid resource type")
	}

	if err != nil {
		return err
	}

	finish_time := time.Now().Add(
		time.Duration(metadata.UpgradeTimeSec) * time.Second,
	)

	_, err = tx.Exec(
		ctx,
		`UPDATE buildings_village
		SET upgrade_ends_at = $1
		WHERE id = $2`, finish_time, building_instance_id)

	return tx.Commit(ctx)
}

func (r *VillageRepository) CompleteUpgrades(village_id int64) error {

	ctx := context.Background()
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx,
		`UPDATE buildings_village
		SET
			level = level + 1,
			upgrade_ends_at = NULL
		WHERE village_id = $1
		AND upgrade_ends_at IS NOT NULL
		AND upgrade_ends_at <= NOW()`, village_id)

	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func isDefense(id int64) bool {
	switch id {
	case 2, 3, 4, 5:
		return true
	default:
		return false
	}
}

func (r *VillageRepository) AddBuilding(village_id int64, building_id int64, x int, y int) error {
	ctx := context.Background()
	var townhall_level int

	err := r.DB.QueryRow(ctx,
		`SELECT townhall_level
		FROM villages
		WHERE id = $1`, village_id).Scan(&townhall_level)

	if err != nil {
		return err
	}

	var cost int
	var costType string

	err = r.DB.QueryRow(ctx,
		`SELECT upgrade_cost, cost_type
		FROM buildings_metadata
		WHERE id = $1`, building_id).Scan(&cost, &costType)

	if err != nil {
		return err
	}

	if isDefense(building_id) {

		var unlock_level int

		err = r.DB.QueryRow(ctx,
			`SELECT unlock_level
			FROM defense_metadata
			WHERE type_id = $1 AND level = 1`, building_id).Scan(&unlock_level)

		if err != nil {
			return err
		}

		if townhall_level < unlock_level {
			return errors.New("Building locked")
		}
	}

	var max_quantity int

	err = r.DB.QueryRow(ctx,
		`SELECT max_quantity
		FROM building_limits
		WHERE building_id = $1
		AND townhall_level = $2`, building_id, townhall_level).Scan(&max_quantity)

	if err != nil {
		return err
	}

	var current_quantity int

	err = r.DB.QueryRow(ctx,
		`SELECT COUNT(*)
		FROM buildings_village
		WHERE village_id = $1
		AND building_id = $2`, village_id, building_id).Scan(&current_quantity)

	if err != nil {
		return err
	}

	if current_quantity >= max_quantity {
		return errors.New("max quantity reached")
	}

	canPlace, err := r.CanPlaceNewBuilding(village_id, building_id, x, y)

	if err != nil {
		return err
	}

	if !canPlace {
		return errors.New("invalid placement")
	}

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if costType == "Gold" {

		cmd, err := tx.Exec(ctx,
			`UPDATE villages
			SET gold = gold - $1
			WHERE id = $2
			AND gold >= $1`, cost, village_id)

		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return errors.New("Not enough gold")
		}
	}

	if costType == "Elixir" {

		cmd, err := tx.Exec(ctx,
			`UPDATE villages
			SET elixir = elixir - $1
			WHERE id = $2 AND elixir >= $1`, cost, village_id)

		if err != nil {
			return err
		}

		if cmd.RowsAffected() == 0 {
			return errors.New("Not enough elixir")
		}
	}

	if building_id == 6 || building_id == 7 {

		_, err = tx.Exec(ctx,
			`INSERT INTO buildings_village
    	(village_id, building_id, level, x, y, last_collected_at)
    	VALUES ($1,$2,1,$3,$4,NOW())`, village_id, building_id, x, y)

	} else {

		_, err = tx.Exec(
			ctx,
			`INSERT INTO buildings_village
			(village_id,building_id,level,x,y)
			VALUES
			($1,$2,1,$3,$4)`, village_id, building_id, x, y)
	}

	if err != nil {
		return err
	}

	return tx.Commit(ctx)

}
