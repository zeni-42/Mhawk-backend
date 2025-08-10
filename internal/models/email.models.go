package models

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	Id             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"userid"`
	OrganizationId uuid.UUID `json:"organizationid"`
	TemplateId     uuid.UUID `json:"templateid"`
	ApiId          uuid.UUID `json:"apiid"`
	LogId          uuid.UUID `json:"logid"`
	Recievers      []string  `json:"recievers"`
	IsBulk         bool      `json:"isbulk"`
	CreatedAt      time.Time `json:"createdat"`
	UpdatedAt      time.Time `json:"updatedat"`
}