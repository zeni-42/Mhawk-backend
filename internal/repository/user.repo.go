package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zeni-42/Mhawk/internal/database"
	"github.com/zeni-42/Mhawk/internal/models"
)

func FindUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user models.User

	psql := `
		SELECT id, fullname, email 
		FROM users 
		WHERE email = $1;
	`

	err := database.DB.QueryRow(ctx, psql, email).Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(userData models.User) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	userData.Id = uuid.New()

	psql := `
		INSERT INTO users (
			id, fullname, email, password
		) VALUES (
			$1, $2, $3, $4
		);
	`

	_, err := database.DB.Exec(ctx, psql, userData.Id, userData.Fullname, userData.Email, userData.Password)
	if  err != nil {
		return uuid.Nil, err
	}

	return userData.Id, nil
}

func FindUserById(id uuid.UUID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user models.User

	psql := `
		SELECT id, fullname, email, is_pro, is_organization, created_at, updated_at
		FROM users
		WHERE id = $1;
	`

	err := database.DB.QueryRow(ctx, psql, id).Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
		&user.IsPro,
		&user.IsOrganization,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}