package seed

import (
	"context"
	"fmt"

	"github.com/anshikagupta17/MVC_Assignment/core/auth"
	"github.com/jackc/pgx/v5"
)

func SeedTestUsers(conn *pgx.Conn) error {
	ctx := context.Background()

	password_hash, err := auth.HashPass("Assignment4")
	if err != nil {
		return err
	}

	for i := 1; i <= 10; i++ {
		username := fmt.Sprintf("test%d", i)

		_, err := conn.Exec(
			ctx,
			`INSERT INTO users(username, pass_hash)
			VALUES($1,$2)
			ON CONFLICT(username) DO NOTHING`, username, password_hash)

		if err != nil {
			return err
		}
	}

	return nil
}
