package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeni-42/Mhawk/internal/database"
	"github.com/zeni-42/Mhawk/internal/models"
)

func SaveEmail(email models.Email) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	psql := `
		INSERT INTO emails (
			user_id, api_id, "to", subject, body, is_html, html 
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		) RETURNING id;
	`

	var id uuid.UUID

	if err := database.DB.QueryRow(ctx, psql, email.UserId, email.ApiId, email.To, email.Subject, email.Body, email.IsHtml, email.Html).Scan(&id); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}