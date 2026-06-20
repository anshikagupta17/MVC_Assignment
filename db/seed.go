package db

import (
	seed "github.com/anshikagupta17/MVC_Assignment/db/seeds"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedAll(conn *pgxpool.Pool) error {

	if err := seed.SeedTestUsers(conn); err != nil {
		return err
	}

	if err := seed.SeedBuildings(conn); err != nil {
		return err
	}

	if err := seed.SeedDefense(conn); err != nil {
		return err
	}

	if err := seed.SeedMines(conn); err != nil {
		return err
	}

	if err := seed.SeedStorage(conn); err != nil {
		return err
	}

	if err := seed.SeedArmy(conn); err != nil {
		return err
	}

	if err := seed.SeedLimits(conn); err != nil {
		return err
	}

	if err := seed.SeedTroops(conn); err != nil {
		return err
	}

	if err := seed.SeedTroopsLevel(conn); err != nil {
		return err
	}

	if err := seed.SeedTestVillages(conn); err != nil {
		return err
	}

	return nil
}
