package db

import (
	"context"
	"log"
	"os"

	"github.com/anshikagupta17/MVC_Assignment/db"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
)

var Conn *pgxpool.Pool

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	db_url := os.Getenv("DB_URI")
	if db_url == "" {
		log.Fatal("DB_URI not set")
	}

	conn, err := pgxpool.New(context.Background(), db_url)
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
		err = db.SeedAll(conn)
		if err != nil {
			log.Fatal(err)
		}
	}
}
