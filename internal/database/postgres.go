package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func Connect() {
	dbUrl := os.Getenv("POSTGRES_URL")
	if dbUrl == "" {
		log.Fatalf("POSTGRES_URL not found")
	}

	var err error
	DB, err = sql.Open("pgx", dbUrl)
	if err != nil {
		log.Fatalf("DB connection failed %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("DB ping failed %v", err)
	}

	fmt.Println("[GIN] DB CONNECTED")
}

func GetPing() bool {
	if DB != nil {
		if err := DB.Ping(); err != nil {
			return false
		}
		return true
	}
	return false
}