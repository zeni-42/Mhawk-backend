package models

import (
	"time"

	"github.com/google/uuid"
)

type Template struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body"`
	AuthorId       uuid.UUID `json:"authorId"`
	OrganizationId uuid.UUID `json:"organizationId"`
	IsRestricted   bool      `json:"isRestricted"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
