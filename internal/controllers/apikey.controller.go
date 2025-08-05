package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/zeni-42/Mhawk/internal/database"
	"github.com/zeni-42/Mhawk/internal/models"
	"github.com/zeni-42/Mhawk/internal/repository"
	"github.com/zeni-42/Mhawk/internal/utils/response"
)

type userData struct {
	UserId		uuid.UUID 		`json:"id"`
	KeyName		string			`json:"keyName"`
}

func findUserFromDB(id uuid.UUID) (*models.User, error) {
	user, err := repository.FindUserById(id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func KeyGenerator() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func GenerateNewApiKey(context *gin.Context) {
	var api models.ApiKey
	var userData userData
	var err error

	if strings.TrimSpace(userData.KeyName) == "" {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Key name cannot be empty"))
		return
	}

	if err := context.BindJSON(&userData); err != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid data"))
		return
	}

	// Checking cache for the user data
	_, err = database.GetUserDataFromRedis(userData.UserId.String());
	if err != nil {
		log.Printf("%v", err)

		// If not found get user from DB
		_, DBerr := findUserFromDB(userData.UserId)
		if DBerr != nil {
			if errors.Is(DBerr, pgx.ErrNoRows) {
				context.IndentedJSON(http.StatusNotFound, response.Error(DBerr, http.StatusNotFound, "User not found"))
				return
			}
		}
		context.IndentedJSON(http.StatusInternalServerError, response.Error(DBerr, http.StatusInternalServerError, "Database error"))
		return
	}

	key := KeyGenerator()
	if key == "" {
		context.IndentedJSON(http.StatusFailedDependency, response.Error(nil, http.StatusFailedDependency, "Failed to generate API key"))
		return
	}

	const keyPrefix = "mhawk+"

	apiKey := userData.KeyName + keyPrefix+ key

	existingApiKey, err :=  repository.FindApiKey(apiKey)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Error checking for existing API key"))
		return
	}

	if existingApiKey != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(nil, http.StatusInternalServerError, "API key already exists"))
		return
	}

	api.KeyName = userData.KeyName
	api.ApiKey = apiKey

	apiID := repository.SaveAPIKey(api)
	if err := repository.UpdateUserApiId(userData.UserId, apiID); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save API ID"))
		return
	}

	context.IndentedJSON(http.StatusOK, response.Success(apiKey, http.StatusOK, "API key"))
}