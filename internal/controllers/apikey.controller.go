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

// Generate an cryptograthically secure API key
func KeyGenerator() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func GenerateNewApiKey(context *gin.Context) {
	var api models.ApiKey
	var err error

	// Bind and validate incoming JSON data
	if err := context.BindJSON(&api); err != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid data"))
		return
	}

	// // Checking for empty fields
	if strings.TrimSpace(api.KeyName) == "" || strings.TrimSpace(api.Description) == "" {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Missing fields"))
		return
	}

	// Checking cache for the user data
	_, err = database.GetUserDataFromRedis(api.UserId.String());
	if err != nil {
		log.Println("User not found in cache")

		// If not found get user from DB
		_, DBerr :=  repository.FindUserById(api.UserId)
		if DBerr != nil {
			if errors.Is(DBerr, pgx.ErrNoRows) {
				context.IndentedJSON(http.StatusNotFound, response.Error(DBerr, http.StatusNotFound, "User not found"))
				return
			}
			context.IndentedJSON(http.StatusInternalServerError, response.Error(DBerr, http.StatusInternalServerError, "Database error"))
			return
		}
	}

	key := KeyGenerator()
	if key == "" {
		context.IndentedJSON(http.StatusFailedDependency, response.Error(nil, http.StatusFailedDependency, "Failed to generate API key"))
		return
	}

	const keyPrefix = "follow_@zeni-42_on_github_"

	apiKey := keyPrefix + key

	existingApiKey, err := repository.FindApiKey(apiKey)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Error checking for existing API key"))
		return
	}

	if existingApiKey != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(nil, http.StatusInternalServerError, "API key already exists"))
		return
	}

	api.ApiKey = apiKey

	if _, err := repository.SaveAPIKey(api); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save API ID"))
		return
	}

	context.IndentedJSON(http.StatusOK, response.Success(apiKey, http.StatusOK, "API key"))
}

func GetUserApiKeys(context *gin.Context) {
	idParams := context.Param("id")
	if idParams == "" {
		context.IndentedJSON(http.StatusBadRequest, response.Error(nil, http.StatusBadRequest, "Missing values"))
		return
	}

	parsedUUID := uuid.MustParse(idParams)

	_, err := repository.FindUserById(parsedUUID);
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusNotFound, response.Error(err, http.StatusNotFound, "User not found"))
		return
	}

	apiKeys, err := repository.FindAllApisFromUserId(parsedUUID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusNotFound, response.Error(err, http.StatusNotFound, "No API key found"))
		return
	}

	context.IndentedJSON(http.StatusOK, response.Success(apiKeys, http.StatusOK, "All API keys"))
}

func DeleteAPI(context *gin.Context) {
	apiId := context.Param("id")
	if apiId == "" {
		context.IndentedJSON(http.StatusBadRequest, response.Error(nil, http.StatusBadRequest, "Missing values"))
		return
	}

	parsedId := uuid.MustParse(apiId)

	rowsAffected, err := repository.FindByIdAndDeleteAPI(parsedId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			context.IndentedJSON(http.StatusNotFound, response.Error(err, http.StatusNotFound, "API not found"))
			return
		}
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Database error"))
		return
	}

	if rowsAffected == 0 {
		context.IndentedJSON(http.StatusNotFound, response.Error(nil, http.StatusNotFound, "API ID does not exist"))
		return
	}

	context.IndentedJSON(http.StatusOK, response.Success(rowsAffected, http.StatusOK, "API removed"))
}

func ToggleActive(context *gin.Context) {
	apiId := context.Param("id")
	if apiId == "" {
		context.IndentedJSON(http.StatusBadRequest, response.Error(nil, http.StatusBadRequest, "Missing values"))
		return
	}

	parsedAPIid := uuid.MustParse(apiId)

	rowsAffected, err := repository.FindByIdAndToggleActive(parsedAPIid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			context.IndentedJSON(http.StatusNotFound, response.Error(err, http.StatusNotFound, "API not found"))
			return
		}
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Database error"))
		return
	}

	if rowsAffected == 0 {
		context.IndentedJSON(http.StatusNotFound, response.Error(nil, http.StatusNotFound, "API ID does not exist"))
		return
	}

	context.IndentedJSON(http.StatusOK, response.Success(rowsAffected, http.StatusOK, "API updated"))
}