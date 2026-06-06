package repositories

import (
	"context"

	"github.com/anshikagupta17/MVC_Assignment/core"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	DB *pgx.Conn
}

func (r *UserRepository) CreateUser(username, hashed_pass string) (int64, error) {
	ctx := context.Background()
	var id int64
	err := r.DB.QueryRow(ctx,
		`INSERT INTO users (username, pass_hash)
		 VALUES ($1, $2)
		 RETURNING id`,
		username,
		hashed_pass,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) GetUser(username string) (core.User, error) {
	ctx := context.Background()
	err := r.DB.QueryRow(ctx,
		`SELECT FROM users (ID, username, pass_hash)
		VALUES ($1, $2, $3)
		RETURNING `,
	)
}
