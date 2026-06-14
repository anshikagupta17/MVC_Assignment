package seed

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TestVillage struct {
	TownhallLevel int
	Trophies      int
	HousingSpace  int
}

func SeedTestVillages(conn *pgx.Conn) error {
	ctx := context.Background()

	rows, err := conn.Query(ctx,
		`SELECT id
		FROM users
		WHERE username LIKE 'test%'
		ORDER BY id`)

	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []int64

	for rows.Next() {

		var userID int64

		err := rows.Scan(&userID)
		if err != nil {
			return err
		}

		userIDs = append(userIDs, userID)
	}

	villages := []TestVillage{
		{1, 100, 20},
		{1, 150, 20},
		{1, 180, 20},
		{2, 220, 30},
		{2, 260, 30},
		{2, 300, 30},
		{3, 350, 40},
		{3, 420, 40},
		{4, 500, 50},
		{4, 600, 50},
	}

	for i, userID := range userIDs {

		if i >= len(villages) {
			break
		}

		village := villages[i]

		var count int

		err := conn.QueryRow(ctx,
			`SELECT COUNT(*)
			FROM villages
			WHERE user_id = $1`, userID).Scan(&count)

		if err != nil {
			return err
		}

		if count > 0 {
			continue
		}

		_, err = conn.Exec(ctx,
			`INSERT INTO villages(
			user_id,
			townhall_level,
			gold,
			elixir,
			housing_space,
			trophies,
			layout
		)
		VALUES(
			$1,
			$2,
			50000,
			50000,
			$3,
			$4,
			$5
		)`, userID, village.TownhallLevel, village.HousingSpace, village.Trophies, `{}`)

		if err != nil {
			return err
		}
	}

	return nil
}
