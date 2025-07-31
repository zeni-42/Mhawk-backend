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
		SELECT id, fullname, email, password
		FROM users 
		WHERE email = $1;
	`

	err := database.DB.QueryRow(ctx, psql, email).Scan(
		&user.Id,
		&user.Fullname,
		&user.Email,
		&user.Password,
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

func UpdateRefreshToken(id uuid.UUID, token interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	psql := `
		UPDATE users 
		SET refresh_token = $2 
		WHERE id = $1;
	`

	_ , err := database.DB.Exec(ctx, psql, id.String(), token)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserAvatar(id uuid.UUID, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	psql := `
		UPDATE users
		SET avatar = $1
		WHERE id = $2
	`

	_, err := database.DB.Exec(ctx, psql, url, id.String())
	if err != nil {
		return err
	}

	return nil
}