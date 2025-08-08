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

	apikey.ExpireDate = time.Now().AddDate(0, 1, 0)

	psql := `
		INSERT INTO apikeys (
			user_id, key_name, api_key, description, environment, expired_date
		) VALUES (
			$1, $2, $3, $4, $5, $6
		);
	`

	if _, err := database.DB.Exec(ctx, psql, apikey.UserId, apikey.KeyName, apikey.ApiKey, apikey.Description, apikey.Environment, apikey.ExpireDate); err != nil {
		return uuid.Nil, err
	}

	insertedApiKey, err := FindApiKey(apikey.ApiKey)
	if err != nil {
		return uuid.Nil, err
	}

	deRefApiId := *insertedApiKey

	return deRefApiId.Id, nil
}

func FindAllApisFromUserId (id uuid.UUID) ([]models.ApiKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	psql := `
		SELECT id, user_id, key_name, api_key, description, is_active, environment, token, expired_date, created_at
		FROM apikeys
		WHERE user_id = $1;
	`

	rows, err := database.DB.Query(ctx, psql, id);
	if err != nil {
		return  nil, err
	}
	defer rows.Close()

	var apiKeys []models.ApiKey

	for rows.Next() {

		var api models.ApiKey

		err := rows.Scan(
			&api.Id,
			&api.UserId,
			&api.KeyName,
			&api.ApiKey,
			&api.Description,
			&api.IsActive,
			&api.Environment,
			&api.Token,
			&api.ExpireDate,
			&api.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		apiKeys = append(apiKeys, api)
	}

	if rows.Err() != nil {
		return  nil, rows.Err()
	}

	return apiKeys, nil
}