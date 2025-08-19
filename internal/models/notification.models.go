package models

import (
	"time"

	"github.com/google/uuid"
)

type Noification struct {
	ID			uuid.UUID	`json:"id"`
	SenderId	uuid.UUID	`json:"sender_id"`
	RecieverId	uuid.UUID	`json:"reciever_id"`
	IsSystem	string		`json:"is_system"`
	Title		string		`json:"title"`
	Message		string		`json:"message"`
	Data		string		`json:"data"`
	IsRead		bool		`json:"is_read"`
	ReadAt		time.Time	`json:"read_at"`
	CreatedAt	time.Time	`json:"created_at"`
}