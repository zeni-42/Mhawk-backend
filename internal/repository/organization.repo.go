package repository

import (
	"context"
	"time"

	"github.com/zeni-42/Mhawk/internal/database"
	"github.com/zeni-42/Mhawk/internal/models"
)

func FindOrganizationByDomanin(domain string) (models.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var existingOrg models.Organization

	psql := `
		SELECT id, founder, name, description, domain, created_at, updated_at
		FROM organization
		WHERE domain = $1;
	`

	if err := database.DB.QueryRow(ctx, psql, domain).Scan(
		&existingOrg.ID,
		&existingOrg.Founder,
		&existingOrg.Name,
		&existingOrg.Description,
		&existingOrg.Domain,
		&existingOrg.CreatedAt,
		&existingOrg.UpdatedAt,
		); err != nil {
			return models.Organization{}, err
	}

	return existingOrg, nil
}

func CreateOrganization(org models.Organization) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	psql := `
	INSERT INTO organization (
		founder, name, description, domain
	) VALUES (
		$1, $2, $3, $4
	);`

	rows, err := database.DB.Exec(ctx, psql, org.Founder, org.Name, org.Description, org.Domain)
	if err != nil {
		return 0, err
	}

	return rows.RowsAffected(), nil
}