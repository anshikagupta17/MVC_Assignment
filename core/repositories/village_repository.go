package repositories

import (
	"context"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VillageRepository struct {
	DB *pgxpool.Pool
}

func (r *VillageRepository) CreateVillage(user_id int64) (int64, error) {
	ctx := context.Background()
	var village_id int64
	err := r.DB.QueryRow(ctx,
		`INSERT INTO villages (
			user_id,
			townhall_level,
			gold,
			elixir,
			housing_space,
			trophies,
			layout
		)
		VALUES (
			$1,
			1,
			750,
			750,
			20,
			0,
			'{}'
		) RETURNING id`,
		user_id,
	).Scan(&village_id)
	if err != nil {
		return 0, err
	}

	return village_id, nil
}

func (r *VillageRepository) GetVillage(userID int64) (models.Village, error) {
	ctx := context.Background()

	var v models.Village

	err := r.DB.QueryRow(ctx,
		`SELECT id, user_id, townhall_level, gold, elixir, housing_space, trophies, layout
		 FROM villages
		 WHERE user_id = $1`,
		userID,
	).Scan(
		&v.ID,
		&v.UserId,
		&v.TownhallLevel,
		&v.Gold,
		&v.Elixir,
		&v.HousingSpace,
		&v.Trophies,
		&v.Layout,
	)

	return v, err
}

func (r *VillageRepository) VillageState(user_id int64) (models.VillageResponse, error) {
	ctx := context.Background()
	var village models.VillageResponse
	err := r.DB.QueryRow(ctx,
		`SELECT id, gold, elixir, townhall_level, layout, trophies
		FROM villages
		where user_id=$1`, user_id).Scan(&village.ID, &village.Gold, &village.Elixir, &village.TownhallLevel, &village.Layout, &village.Trophies)
	if err != nil {
		return models.VillageResponse{}, err
	}

	err = CompleteUpgrades(ctx, r.DB, village.ID)
	if err != nil {
		return models.VillageResponse{}, err
	}

	err = r.CompleteTroopUpgrades(village.ID)
	if err != nil {
		return models.VillageResponse{}, err
	}

	buildings, err := r.VillageBuildings(village.ID)
	if err != nil {
		return models.VillageResponse{}, err
	}

	village.Buildings = buildings

	return village, nil

}
