package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/zeni-42/Mhawk/internal/repository"
	"github.com/zeni-42/Mhawk/internal/utils/response"
)

type EmailForm struct {
	UserId		uuid.UUID	`json:"userid"`
	ApiKeyId	uuid.UUID	`json:"apikeyid"`	
	To 			string		`json:"to"`
	Subject		string		`json:"subject"`
	Body		string		`json:"body"`
	IsHTml		bool		`json:"isHtml"`
	Html		string		`json:"html"`
}

func SendEmail(context *gin.Context) {
	var userForm EmailForm

	if err := context.BindJSON(&userForm); err != nil{
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest,  "Invalid data"))
		return
	}

	_, err := repository.FindUserById(userForm.UserId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			context.IndentedJSON(http.StatusNotFound, response.Error(err, http.StatusNotFound, "User not found"))
			return 
		}
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "USER table error"))
		return 
	}

	apikey, err := repository.FindAPIUsingId(userForm.ApiKeyId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			context.IndentedJSON(http.StatusNotFound, response.Error(err, http.StatusNotFound, "API key not found"))
			return 
		}
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "API table error"))
		return 
	}

	url := os.Getenv("EMAIL_SERVER")
	if url == "" {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(nil, http.StatusInternalServerError, "Missing env value"))
		return
	}

	rowAffected, err := repository.UpdateAPIToken(apikey)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to updated token"))
		return
	}

	if rowAffected == 0 {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Row affected is 0"))
		return
	}

	bodyMap := map[string]interface{}{
		"to": userForm.To,
		"subject": userForm.Subject,
		"body": userForm.Body,
		"isHtml": false,
		"html": "",
	}

	body, err := json.Marshal(bodyMap)
	if err != nil {
		log.Println("Failed to prepare body")
		return
	}

	r, err := http.NewRequest(http.MethodPost, url + "/send", bytes.NewBuffer(body))
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(nil, http.StatusInternalServerError, "Failed to call EMAIL SERVICE"))
		return 
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to call EMAIL service"))
		return
	}

	defer res.Body.Close()

	resBody, _ := io.ReadAll(r.Body)
	context.IndentedJSON(http.StatusOK, response.Success(resBody, http.StatusOK, "Email sent successfully"))
}