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
	IsNew			bool		`json:"isnew"`
	Avatar			string		`json:"avatar"`
	IsPro           bool      	`json:"ispro"`
	CreatedAt       time.Time 	`json:"createdat"`
	UpdatedAt       time.Time 	`json:"updatedat"`
}