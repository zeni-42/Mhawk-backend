package models

import (
	"time"

	"github.com/google/uuid"
)

type EmailStatus string;

const (
	EmailStatusPending		EmailStatus = "Pending"
	EmailStatusSent			EmailStatus = "Sent"
	EmailStatusDeliverd		EmailStatus = "Delivered"
	EmailStatusFailed		EmailStatus = "Failed"
)

type Email struct {
	Id             	uuid.UUID 	`json:"id"`
	UserId         	uuid.UUID 	`json:"userid"`
	ApiId          	uuid.UUID 	`json:"apiid"`
	To				string		`json:"to"`
	Subject			string		`json:"subject"`
	Body			string		`json:"body"`
	Html			string		`json:"html"`
	IsHtml			bool		`json:"ishtml"`
	Status			EmailStatus	`json:"status"`
	IsTemplate		bool		`json:"istemplate"`	
	TemplateId     	uuid.UUID 	`json:"templateid"`
	IsBulk         	bool      	`json:"isbulk"`
	LogId          	uuid.UUID 	`json:"logid"`
	CreatedAt      	time.Time 	`json:"createdat"`
	UpdatedAt      	time.Time 	`json:"updatedat"`
}
