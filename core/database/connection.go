package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

func connection() *pgx.Conn {
	connect := "host=localhost port=5432 user=postgres dbname=MVC_Assignment sslmode=disable"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, connect)

	if err != nil {
		log.Fatal("Connection not made: ", err)

	}

	err = conn.Ping(ctx)

	if err != nil {
		log.Fatal("DB not connected: ", err)
	}

	return conn
}
