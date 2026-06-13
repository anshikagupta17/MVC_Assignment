package seed

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func SeedLimits(conn *pgx.Conn) error {
	ctx := context.Background()

	var count int
	err := conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM building_limits`,
	).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = conn.Exec(ctx, `
	INSERT INTO building_limits
	(building_id, townhall_level, max_quantity)
	VALUES
	(1,1,1),(1,2,1),(1,3,1),(1,4,1),

	(2,1,1),(2,2,2),(2,3,2),(2,4,2),

	(3,1,0),(3,2,1),(3,3,2),(3,4,2),

	(4,1,10),(4,2,20),(4,3,30),(4,4,40),

	(5,1,0),(5,2,0),(5,3,1),(5,4,2),

	(6,1,1),(6,2,2),(6,3,3),(6,4,4),

	(7,1,1),(7,2,2),(7,3,3),(7,4,4),

	(8,1,1),(8,2,1),(8,3,2),(8,4,2),

	(9,1,1),(9,2,1),(9,3,2),(9,4,2),

	(10,1,1),(10,2,1),(10,3,2),(10,4,2)
	`)

	return err
}
