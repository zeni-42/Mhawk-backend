package models

import (
	"time"

	"github.com/google/uuid"
)
type ApiKey struct {
	Id          uuid.UUID 	`json:"id"`
	UserId		uuid.UUID	`json:"userid"`
	KeyName		string		`json:"keyname"`	
	ApiKey      string 		`json:"apikey"`
	Description	string		`json:"description"`
	IsActive	bool		`json:"isactive"`
	Environment	string		`json:"environment"`
	UsedToken   int32     	`json:"usedtoken"`
	ExpireDate  time.Time 	`json:"expiredate"`
	RefreshDate time.Time 	`json:"refreshdate"`
	CreatedAt   time.Time 	`json:"createdat"`
	UpdatedAt	time.Time	`json:"updatedat"`
}