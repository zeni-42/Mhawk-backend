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

func SaveAPIKey(apikey models.ApiKey) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second) 
	defer cancel()

	psql := `
		INSERT INTO apikeys (
			user_id, key_name, api_key, description, env
		) VALUES (
			$1, $2, $3, $4, $5
		);
	`

	if _, err := database.DB.Exec(ctx, psql, apikey.UserId, apikey.KeyName, apikey.ApiKey, apikey.Description, apikey.Environment); err != nil {
		return uuid.Nil, err
	}

	insertedApiKey, err := FindApiKey(apikey.ApiKey)
	if err != nil {
		return uuid.Nil, err
	}

	deRefApiId := *insertedApiKey

	return deRefApiId.Id, nil
}