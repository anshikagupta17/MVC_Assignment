package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedArmy(conn *pgxpool.Pool) error {
	ctx := context.Background()

	var count int
	err := conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM army_metadata`,
	).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = conn.Exec(ctx, `
	INSERT INTO army_metadata
	(type_id, level, capacity)
	VALUES
	(10, 1, 20),
	(10, 2, 30),
	(10, 3, 35),
	(10, 4,40)
	`)

	return err
}
