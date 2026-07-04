package repositories

import (
	"context"
	"errors"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func CreateUser(ctx context.Context, db DBExecutor, username, hashed_pass string) (int64, error) {
	var id int64
	err := db.QueryRow(ctx,
		`INSERT INTO users (username, pass_hash)
		 VALUES ($1, $2)
		 RETURNING id`, username, hashed_pass).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, errors.New("Username already taken")
		}
		return 0, err
	}

	return id, nil
}

func (r *UserRepository) RegisterTX(username string, hashed_pass string) (int64, error) {
	ctx := context.Background()

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	user_id, err := CreateUser(ctx, tx, username, hashed_pass)
	if err != nil {
		return 0, err
	}

	village_id, err := CreateVillage(ctx, tx, user_id)
	if err != nil {
		return 0, err
	}

	err = InitialBuildings(ctx, tx, village_id)
	if err != nil {
		return 0, err
	}

	return user_id, tx.Commit(ctx)
}

func (r *UserRepository) GetUser(username string) (models.User, error) {
	ctx := context.Background()
	var user models.User
	err := r.DB.QueryRow(ctx,
		`SELECT id, username, pass_hash
		FROM users
		WHERE username=$1 `,
		username).Scan(&user.ID, &user.UserName, &user.PassWord)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
