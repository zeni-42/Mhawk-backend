package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeni-42/Mhawk/internal/database"
	"github.com/zeni-42/Mhawk/internal/models"
)

func FindApiKey(apikey string) (*models.ApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) 
	defer cancel()
	var apiKey models.ApiKey

	psql := `
		SELECT id, key_name, api_key
		FROM apikeys 
		WHERE api_key = $1;
	`

	err := database.DB.QueryRow(ctx, psql, apikey).Scan(
		&apiKey.Id,
		&apiKey.KeyName,
		&apiKey.ApiKey,
	)
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func SaveAPIKey(apikey models.ApiKey) uuid.UUID {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) 
	defer cancel()

	psql := `
		INSET INTO apikeys (
			key_name, api_key
		) VALUES (
			$1, $2
		);
	`

	var InsertedId uuid.UUID

	if err := database.DB.QueryRow(ctx, psql, apikey.KeyName, apikey.ApiKey).Scan(&InsertedId); err != nil {
		return uuid.Nil
	}

	insertedApiKey, err := FindApiKey(apikey.ApiKey)
	if err != nil {
		return uuid.Nil
	}

	deRefApiId := *insertedApiKey

	return deRefApiId.Id
}