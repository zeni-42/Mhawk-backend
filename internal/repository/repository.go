package repository

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitTables(client *pgxpool.Pool) {
	if client == nil {
		log.Fatalf("DB client is nil can't proceed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exist bool
	err := client.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'users'
		);
	`).Scan(&exist)

	if err != nil {
		log.Fatalf("Failed to check if 'users' table exists: %v", err)
	}

	if !exist {
		query := 
		`CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			fullname TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			refresh_token TEXT,
			free_token INTEGER DEFAULT 50,
			is_pro BOOLEAN DEFAULT FALSE,
			is_organization BOOLEAN DEFAULT FALSE,
			api_id UUID,
			organization_id UUID,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);`
		
		_ , err := client.Exec(ctx, query)
		if err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
	}
}