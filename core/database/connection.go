package db
import {
	"log"
	_ "github.com/lib/pq"
	"database/sql"
}

func connection() *sql.DB {
	connect:= "host=localhost port=5432 username=postgres dbname=MVC_Assignment sslmode=disable"

	db, err:= sql.Open("postgres", connect)

	if err!=nil {
		log.Fatal("Connection not made: ", err)

	}

	err= db.Ping()

	if err!=nil {
		log.Fatal("DB not connected: ", err)
	}

	return db
}
