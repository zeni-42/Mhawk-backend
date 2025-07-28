package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailStatusCode string

const (
	StatusSent      EmailStatusCode = "sent"
	StatusDelivered EmailStatusCode = "delivered"
	StatusFailed    EmailStatusCode = "failed"
	StatusQueued    EmailStatusCode = "queued"
	StatusBounced   EmailStatusCode = "bounced"
)

type Log struct {
	Id             uuid.UUID       `json:"id"`
	EmailId        uuid.UUID       `json:"emailId"`
	OrganizationId uuid.UUID       `json:"organizationId"`
	Status         EmailStatusCode `json:"status"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}