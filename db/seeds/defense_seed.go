package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedDefense(conn *pgxpool.Pool) error {
	ctx := context.Background()

	var count int
	err := conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM defense_metadata`,
	).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = conn.Exec(ctx, `
	INSERT INTO defense_metadata
	(type_id, level, unlock_level, damage, max_health, range)
	VALUES
	(1, 1, 1, 0, 1500, 0),
	(1, 2, 1, 0, 1600, 0),
	(1, 3, 1, 0, 1850, 0),
	(1, 4, 1, 0, 2100, 0),
	(2, 1, 1, 9, 420, 9),
	(2, 2, 1, 11, 470, 9),
	(2, 3, 1, 15, 520, 9),
	(2, 4, 1, 19, 570, 9),
	(3, 1, 2, 11, 380, 10),
	(3, 2, 2, 15, 420,10),
	(3, 3, 2, 19, 460, 10),
	(4, 1, 2, 0, 300, 0),
	(4, 2, 2, 0, 500, 0),
	(4, 3, 2, 0, 700, 0),
	(5, 1, 3, 4, 400, 11),
	(5, 2, 3, 5, 450, 11)
	`)

	return err
}
