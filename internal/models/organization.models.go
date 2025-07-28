package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id        uuid.UUID `json:"id"`
	Founder   uuid.UUID `json:"founder"`
	Members   []string  `json:"members"`
	FreeToken int32     `json:"freeToken"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
