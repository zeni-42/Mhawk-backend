package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id				uuid.UUID	`json:"id"`
	Fullname		string		`json:"fullname"`
	Email			string		`json:"email"`
	Password        string    	`json:"password"`
	RefreshToken    string    	`json:"refreshtoken"`
	FreeToken       int32     	`json:"freetoken"`
	IsPro           bool      	`json:"ispro"`
	IsOrganization  bool      	`json:"isorganization"`
	ApiId           uuid.UUID 	`json:"apiid"`
	OrganizationId  uuid.UUID 	`json:"organizationid"`
	CreatedAt       time.Time 	`json:"createdat"`
	UpdatedAt       time.Time 	`json:"updatedat"`
}