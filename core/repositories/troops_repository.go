package repositories

import (
	"context"
	"encoding/json"
	"log"

	"github.com/anshikagupta17/MVC_Assignment/core/models"
	"github.com/jackc/pgx/v5"
)

type TroopsConn struct {
	DB *pgx.Conn
}

func (r *TroopsConn) SaveArmy(villageID int, troop models.Troop) error {
	ctx := context.Background()

	data, err := json.Marshal(troop)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO troops_village (village_id, army)
		VALUES ($1, $2)
		ON CONFLICT (village_id)
		DO UPDATE SET army = EXCLUDED.army
	`

	_, err = r.DB.Exec(ctx, query, villageID, data)
	if err != nil {
		log.Println("SaveArmy error:", err)
		return err
	}

	return nil
}
