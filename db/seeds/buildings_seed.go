package seed

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func SeedBuildings(conn *pgx.Conn) error {
	ctx := context.Background()

	var count int
	err := conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM buildings_metadata`,
	).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = conn.Exec(ctx, `
	INSERT INTO buildings_metadata
	(id, name, size_x, size_y, upgrade_cost, cost_type, upgrade_time_sec)
	VALUES
	(1,'Townhall', 4, 4, 2000, 'Gold', 600),
	(2,'Cannon', 2, 2, 200, 'Gold', 60),
	(3,'Archer Tower', 2,2, 250, 'Gold', 70),
	(4,'Wall', 1,1, 35, 'Gold', 30),
	(5,'Mortar', 3,3, 600, 'Gold', 300),
	(6,'Gold Mine', 3, 3, 150, 'Elixir', 90),
	(7,'Elixir Mine', 3, 3, 150, 'Gold', 90),
	(8,'Gold Storage', 3, 3, 200, 'Gold', 120),
	(9,'Elixir Storage', 3, 3, 200, 'Gold', 120),
	(10,'Army Camp', 3, 3, 200, 'Elixir', 120)
	`)

	return err
}
