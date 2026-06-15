package seed

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedStorage(conn *pgxpool.Pool) error {
	ctx := context.Background()

	var count int
	err := conn.QueryRow(ctx,
		`SELECT COUNT(*) FROM storage_metadata`).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	_, err = conn.Exec(ctx, `
	INSERT INTO storage_metadata
	(type_id, level, max_capacity)
	VALUES
	(8, 1, 1500),
	(8, 2, 3000),
	(8, 3, 6000),
	(8, 4, 12000),
	(9, 1, 1500),
	(9, 2, 3000),
	(9, 3, 6000),
	(9, 4, 12000)
	`)

	return err
}
