package db

import (
	"context"
	"log"
	"os"

	seed "github.com/anshikagupta17/MVC_Assignment/db/seeds"
	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDB() {
	db_url := "postgres://anshikagupta@localhost:5432/assignment_mvc?sslmode=disable"
	log.Println(os.Getenv("DB_URI"))

	conn, err := pgx.Connect(context.Background(), db_url)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("DB not responding:", err)
	}

	Conn = conn
	log.Println("DB connected")

	err = seed.SeedBuildings(Conn)
	if err != nil {
		log.Fatal(err)
	}
}
