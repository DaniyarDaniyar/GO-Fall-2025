package config

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectDB() *sqlx.DB {
	dsn := "postgres://postgres:postgres123@localhost:5432/practice5?sslmode=disable"

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("DB connection fail: %v", err)
	}

	log.Println("DB connected")
	return DB
}
