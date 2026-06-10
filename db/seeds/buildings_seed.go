package seed

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func SeedBuildings(conn *pgx.Conn) error {
	_, err := conn.Exec(
		context.Background(),
		`
		INSERT INTO buildings_metadata
		(id, size_x, size_y, name, upgrade_cost, cost_type, upgrade_time_sec)
		VALUES
		(1, 4, 4, 'Town Hall', 1000, 'Gold', 60),
		(6, 3, 3, 'Gold Mine', 150, 'Gold', 30),
		(7, 3, 3, 'Elixir Mine', 150, 'Elixir', 30),
		(10, 5, 5, 'Army Camp', 250, 'Elixir', 45)
		ON CONFLICT (id) DO NOTHING`)

	return err
}
