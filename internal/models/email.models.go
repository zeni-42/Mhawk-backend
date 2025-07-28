package models

import (
	"time"

	"github.com/google/uuid"
)

type Email struct {
	Id             uuid.UUID `json:"id"`
	UserId         uuid.UUID `json:"userId"`
	OrganizationId uuid.UUID `json:"organizationId"`
	TemplateId     uuid.UUID `json:"templateId"`
	ApiId          uuid.UUID `json:"apiId"`
	LogId          uuid.UUID `json:"logId"`
	Recievers      []string  `json:"recievers"`
	IsBulk         bool      `json:"isBulk"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}