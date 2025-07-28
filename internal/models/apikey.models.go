package models

import (
	"time"

	"github.com/google/uuid"
)

type ApiKey struct {
	Id          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"userId"`
	ApiKey      uuid.UUID `json:"apiKey"`
	UsedToken   int32     `json:"usedToken"`
	ExpireDate  time.Time `json:"expireDate"`
	RefreshDate time.Time `json:"refreshDate"`
	CreatedAt   time.Time `json:"createdAt"`
}