package repositories

import (
	"context"
	"errors"
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
	ctx := context.Background()
	rows, err := r.DB.Query(ctx,
		`SELECT bv.id, bv.building_id, bv.level, bv.upgrade_ends_at, bv.x, bv.y,
		bm.size_x, bm.size_y
		From buildings_village bv 
		JOIN buildings_metadata bm ON bm.id = bv.building_id
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
			&building.SizeX,
			&building.SizeY,
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
func InitialBuildings(ctx context.Context, db DBExecutor, village_id int64) error {
	_, err := db.Exec(ctx,
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

func CanPlace(ctx context.Context, db DBExecutor, village_id int64, size_x int, size_y int, x int, y int, skip_instance_id int64) (bool, error) {
	if x+size_x-1 > 49 || y+size_y-1 > 49 {
		return false, nil
	}

	rows, err := db.Query(ctx,
		`SELECT bv.id, bv.x, bv.y, bm.size_x, bm.size_y
		FROM buildings_village bv
		JOIN buildings_metadata bm ON bm.id = bv.building_id
		WHERE bv.village_id = $1`, village_id)

	if err != nil {
		return false, err
	}
	defer rows.Close()

	newX1 := x
	newY1 := y
	newX2 := x + size_x - 1
	newY2 := y + size_y - 1

	for rows.Next() {
		var instance_id int64
		var current_x, current_y, current_sizex, current_sizey int

		err := rows.Scan(
			&instance_id,
			&current_x,
			&current_y,
			&current_sizex,
			&current_sizey,
		)
		if err != nil {
			return false, err
		}

		if instance_id == skip_instance_id {
			continue
		}

		currentX1 := current_x
		currentY1 := current_y
		currentX2 := current_x + current_sizex - 1
		currentY2 := current_y + current_sizey - 1

		if !(newX2 < currentX1 || newX1 > currentX2 || newY2 < currentY1 || newY1 > currentY2) {
			return false, nil
		}
	}

	if err := rows.Err(); err != nil {
		return false, err
	}

	return true, nil
}
func BuildingUpgrade(ctx context.Context, db DBExecutor, village_id int64, building_instance_id int64) error {
	err := CompleteUpgrades(ctx, db, village_id)
	if err != nil {
		return err
	}

	var building_id int64
	var level int
	var upgrade_ends_at pgtype.Timestamp

	err = db.QueryRow(ctx,
		`SELECT building_id, level, upgrade_ends_at
		FROM buildings_village
		WHERE id = $1
		AND village_id = $2
		FOR UPDATE`, building_instance_id, village_id).Scan(&building_id, &level, &upgrade_ends_at)

	if err != nil {
		return err
	}

	if upgrade_ends_at.Valid {
		return errors.New("Building already upgrading")
	}

	var upgradingCount int
	err = db.QueryRow(ctx,
		`SELECT COUNT(*)
		FROM buildings_village
		WHERE village_id = $1
		AND upgrade_ends_at IS NOT NULL
		AND id != $2`, village_id, building_instance_id).Scan(&upgradingCount)

	if err != nil {
		return err
	}

	if upgradingCount > 0 {
		return errors.New("Another building is already upgrading")
	}

	var townhall_level int
	var resources VillageResources

	err = db.QueryRow(ctx,
		`SELECT townhall_level, gold, elixir
		FROM villages
		WHERE id = $1
		FOR UPDATE`, village_id).Scan(&townhall_level, &resources.Gold, &resources.Elixir)

	if err != nil {
		return err
	}

	if building_id == 1 && level >= 4 {
		return errors.New("Townhall has reached its maximum level of 4")
	}
	if building_id != 1 {
		if level+1 > townhall_level {
			return errors.New("building cannot exceed townhall level")
		}
	}

	var metadata BuildingMetadata

	err = db.QueryRow(ctx,
		`SELECT upgrade_cost, cost_type, upgrade_time_sec
		FROM buildings_metadata
		WHERE id = $1`,
		building_id).Scan(&metadata.UpgradeCost, &metadata.CostType, &metadata.UpgradeTimeSec)

	if err != nil {
		return err
	}

	switch metadata.CostType {

	case "Gold":

		if resources.Gold < metadata.UpgradeCost {
			return errors.New("Not enough gold")
		}

		_, err = db.Exec(ctx,
			`UPDATE villages
			SET gold = gold - $1
			WHERE id = $2`, metadata.UpgradeCost, village_id)

	case "Elixir":

		if resources.Elixir < metadata.UpgradeCost {
			return errors.New("Not enough elixir")
		}

		_, err = db.Exec(ctx,
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

	_, err = db.Exec(
		ctx,
		`UPDATE buildings_village
		SET upgrade_ends_at = $1
		WHERE id = $2`, finish_time, building_instance_id)

	return err
}

func CompleteUpgrades(ctx context.Context, db DBExecutor, village_id int64) error {
	var completedTownhall int

	err := db.QueryRow(ctx,
		`SELECT COUNT(*)
		FROM buildings_village
		WHERE village_id = $1
		AND building_id = 1
		AND upgrade_ends_at IS NOT NULL
		AND upgrade_ends_at <= NOW()`,
		village_id,
	).Scan(&completedTownhall)

	if err != nil {
		return err
	}

	if completedTownhall > 0 {
		_, err = db.Exec(ctx,
			`UPDATE villages
			SET townhall_level = LEAST(townhall_level + 1, 4)
			WHERE id = $1`,
			village_id,
		)

		if err != nil {
			return err
		}
	}
	_, err = db.Exec(ctx,
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

	return nil
}

func (r *VillageRepository) BuildingUpgradeTX(village_id int64, building_instance_id int64) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = BuildingUpgrade(ctx, tx, village_id, building_instance_id)
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

func AddBuilding(ctx context.Context, db DBExecutor, village_id int64, building_id int64, x int, y int) error {
	if x < 0 || x > 49 || y < 0 || y > 49 {
		return errors.New("invalid position")
	}

	var townhall_level int

	err := db.QueryRow(ctx,
		`SELECT townhall_level
		FROM villages
		WHERE id = $1
		FOR UPDATE`, village_id).Scan(&townhall_level)

	if err != nil {
		return err
	}

	var cost int
	var costType string
	var sizex, sizey int

	err = db.QueryRow(ctx,
		`SELECT upgrade_cost, cost_type, size_x, size_y
		FROM buildings_metadata
		WHERE id = $1`, building_id).Scan(&cost, &costType, &sizex, &sizey)

	if err != nil {
		return err
	}

	if isDefense(building_id) {

		var unlock_level int

		err = db.QueryRow(ctx,
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

	err = db.QueryRow(ctx,
		`SELECT max_quantity
		FROM building_limits
		WHERE building_id = $1
		AND townhall_level = $2`, building_id, townhall_level).Scan(&max_quantity)

	if err != nil {
		return err
	}

	var current_quantity int

	err = db.QueryRow(ctx,
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

	canPlace, err := CanPlace(ctx, db, village_id, sizex, sizey, x, y, 0)
	if err != nil {
		return err
	}

	if !canPlace {
		return errors.New("invalid placement")
	}

	if costType == "Gold" {

		cmd, err := db.Exec(ctx,
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

		cmd, err := db.Exec(ctx,
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

		_, err = db.Exec(ctx,
			`INSERT INTO buildings_village
    	(village_id, building_id, level, x, y, last_collected_at)
    	VALUES ($1,$2,1,$3,$4,NOW())`, village_id, building_id, x, y)

	} else {

		_, err = db.Exec(
			ctx,
			`INSERT INTO buildings_village
			(village_id,building_id,level,x,y)
			VALUES
			($1,$2,1,$3,$4)`, village_id, building_id, x, y)
	}

	if err != nil {
		return err
	}

	return nil

}

func (r *VillageRepository) AddBuildingTX(village_id int64, building_id int64, x int, y int) error {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = AddBuilding(ctx, tx, village_id, building_id, x, y)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *VillageRepository) GetShopBuildings(village_id int64) ([]models.ShopBuilding, error) {
	ctx := context.Background()

	var townhall_level int
	err := r.DB.QueryRow(ctx,
		`SELECT townhall_level
        FROM villages
        WHERE id = $1`, village_id).Scan(&townhall_level)

	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx,
		`SELECT bm.id, bm.name, bm.upgrade_cost, bm.cost_type, bm.size_x, bm.size_y,
			bl.max_quantity,
		COUNT(bv.id) AS current_count
		FROM building_limits bl
		JOIN buildings_metadata bm
			ON bm.id = bl.building_id
		LEFT JOIN buildings_village bv
			ON bv.building_id = bm.id
			AND bv.village_id = $1
		LEFT JOIN defense_metadata dm
        	ON dm.type_id = bm.id AND dm.level = 1
		WHERE bl.townhall_level = $2 AND bm.id != 1
		AND (dm.unlock_level IS NULL OR dm.unlock_level <= $2)
		GROUP BY
			bm.id,
			bm.name,
			bm.upgrade_cost,
			bm.cost_type,
			bm.size_x,
			bm.size_y,
			bl.max_quantity
		HAVING COUNT(bv.id) < bl.max_quantity
		ORDER BY bm.id`, village_id, townhall_level)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.ShopBuilding

	for rows.Next() {
		var b models.ShopBuilding

		err := rows.Scan(
			&b.BuildingID,
			&b.Name,
			&b.Cost,
			&b.CostType,
			&b.SizeX,
			&b.SizeY,
			&b.MaxQuantity,
			&b.CurrentCount,
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
