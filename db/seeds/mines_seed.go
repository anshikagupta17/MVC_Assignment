package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedMines(conn *pgxpool.Pool) error {
	ctx := context.Background()

	var count int
	err := conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM mines_metadata`,
	).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = conn.Exec(ctx, `
	INSERT INTO mines_metadata
	(type_id, level, production_rate)
	VALUES
	(6, 1, 0.05),
	(6, 2, 0.1),
	(6, 3, 0.15),
	(6, 4, 0.2),
	(7, 1, 0.05),
	(7, 2, 0.1),
	(7, 3, 0.15),
	(7, 4, 0.2);
	`)

	return err
}
