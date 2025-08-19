package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID        	uuid.UUID	`json:"id"`
	Founder   	uuid.UUID	`json:"founder"`
	Name		string		`json:"name"`
	Description	string		`json:"description"`
	Domain		string		`json:"domain"`
	CreatedAt 	time.Time	`json:"created_at"`
	UpdatedAt 	time.Time	`json:"updated_at"`
}

type OrganizationAPIKey struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	APIID          uuid.UUID `json:"api_id"`
}

type OrganizationMember struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
}