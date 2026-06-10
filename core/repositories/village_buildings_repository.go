package repositories

import (
	"context"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
)

func (r *VillageRepository) VillageBuildings(village_id int64) ([]models.VillageBuilding, error) {
	ctx := context.Background()
	rows, err := r.DB.Query(ctx,
		`SELECT id, building_id, level, quantity, upgrade_ends_at, x, y
		From buildings_village
		WHERE village_id= $1`, village_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.VillageBuilding
	for rows.Next() {
		var building models.VillageBuilding
		err := rows.Scan(
			&building.ID,
			&building.BuildingId,
			&building.Level,
			&building.Quantity,
			&building.UpgradeEndsAt,
			&building.X,
			&building.Y,
		)

		if err != nil {
			return nil, err
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

		var instanceID int64
		var b_id int
		var current_x, current_y int

		err := rows.Scan(
			&instanceID,
			&b_id,
			&current_x,
			&current_y,
		)
		if err != nil {
			return false, err
		}
		if instanceID == building_instance_id {
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

func (r *VillageRepository) BuildingUpgrade(village_id, building_id int64) {}
