package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectPG() (*pgxpool.Pool, error) {
	dbUrl := os.Getenv("POSTGRES_URL")
	if dbUrl == "" {
		log.Fatalf("POSTGRES_URL not found")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("DB connection failed %v", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("DB ping failed %v", err)
	}

	log.Println("DB CONNECTED")
	return DB, nil
}

func GetPing() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if DB != nil {
		if err := DB.Ping(ctx); err != nil {
			return false
		}
		return true
	}
	return false
}

func DisconnectPG() error {
	if DB != nil {
		DB.Close()
		return nil
	}
	return fmt.Errorf("DB client is empty")
}