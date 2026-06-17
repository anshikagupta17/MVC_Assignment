package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedVillageBuildings(conn *pgxpool.Pool, village_id int64, townhall_level int) error {
	ctx := context.Background()

	type Building struct {
		ID    int
		Level int
		X     int
		Y     int
	}

	var buildings []Building

	switch townhall_level {

	case 1:
		buildings = []Building{
			{1, 1, 10, 10},
			{6, 1, 5, 5},
			{7, 1, 15, 5},
			{10, 1, 10, 15},
		}

	case 2:
		buildings = []Building{
			{1, 2, 10, 10},
			{6, 2, 5, 5},
			{7, 2, 15, 5},
			{10, 2, 10, 15},
			{2, 1, 5, 15},
		}

	case 3:
		buildings = []Building{
			{1, 3, 10, 10},
			{6, 3, 5, 5},
			{7, 3, 15, 5},
			{10, 3, 10, 15},
			{2, 2, 5, 15},
			{3, 1, 15, 15},
			{8, 1, 5, 10},
			{9, 1, 15, 10},
		}

	default:
		buildings = []Building{
			{1, 4, 10, 10},
			{6, 4, 5, 5},
			{7, 4, 15, 5},
			{10, 4, 10, 15},
			{2, 3, 5, 15},
			{3, 2, 15, 15},
			{5, 1, 10, 5},
			{8, 2, 5, 10},
			{9, 2, 15, 10},
		}
	}

	for _, b := range buildings {

		_, err := conn.Exec(ctx,
			`INSERT INTO buildings_village
			(village_id, building_id, level, x, y)
			VALUES ($1,$2,$3,$4,$5)`,
			village_id,
			b.ID,
			b.Level,
			b.X,
			b.Y,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
