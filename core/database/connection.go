package db

import (
	"context"
	"log"
	"os"

	seed "github.com/anshikagupta17/MVC_Assignment/db/seeds"
	"github.com/jackc/pgx/v5"

	"github.com/joho/godotenv"
)

var Conn *pgx.Conn

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db_url := os.Getenv("DB_URI")
	if db_url == "" {
		log.Fatal("DB_URI not set")
	}

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

	if os.Getenv("SEED_DB") == "true" {
		err = seed.SeedBuildings(Conn)
		if err != nil {
			log.Fatal(err)
		}
	}
}
