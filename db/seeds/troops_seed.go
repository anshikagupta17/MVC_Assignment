package seed

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func SeedTroops(conn *pgx.Conn) error {
	ctx := context.Background()
	var count int
	err := conn.QueryRow(ctx,
		`SELECT count(*) FROM troops_base_metadata`).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err = conn.Exec(ctx,
		`INSERT INTO troops_base_metadata
		(id,name,cost_type,training_cost,speed,unlock_level,upgrade_time_sec,housing_space)
		VALUES
		(1,'Barbarian','Elixir',50,18,1,1080,1),
		(2,'Archer','Elixir',85,24,2,1080,1),
		(3,'Giant','Elixir',500,12,3,1440,5),
		(4,'Wizard','Elixir',1500,16,4,0,4),
		(5,'Goblin','Elixir',25,32,4,0,1);)`)
	return err
}

func SeedTroopsLevel(conn *pgx.Conn) error {
	ctx := context.Background()
	var count int
	err := conn.QueryRow(ctx,
		`SELECT count(*) FROM troops_level_metadata`).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err = conn.Exec(ctx,
		`INSERT INTO troops_level_metadata
		(type_id,level,damage,max_health,upgrade_cost)
		VALUES

		(1,1,9,45,50),
		(1,2,12,54,100),
		(1,3,15,65,500),
		(1,4,18,85,1000),

		(2,1,7,25,60),
		(2,2,10,30,180),
		(2,3,13,36,450),
		(2,4,17,43,1100),

		(3,1,11,300,120),
		(3,2,14,360,350),
		(3,3,19,430,900),
		(3,4,24,520,2000),

		(4,1,50,80,300),
		(4,2,65,95,800),
		(4,3,85,115,2000),
		(4,4,110,140,4000),
		
		(5,1,9,20,40),
		(5,2,12,24,120),
		(5,3,16,30,300),
		(5,4,21,38,800);`)

	return err

}
