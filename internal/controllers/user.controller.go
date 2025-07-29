package controllers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/zeni-42/Mhawk/internal/models"
	"github.com/zeni-42/Mhawk/internal/repository"
	"github.com/zeni-42/Mhawk/internal/utils/response"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(context *gin.Context) {
	var user models.User

	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid data"))
	}

	existingUser, err := repository.FindUserByEmail(user.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Somthing went wrong"))
	}

	if existingUser != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Email is taken"))
	}

	byteHashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	user.Password = string(byteHashedPass)

	userId, err := repository.CreateUser(user); 
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save user"))
	}

	userPointer, err := repository.FindUserById(userId)
	if errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save user"))
	}

	user = *userPointer

	context.IndentedJSON(http.StatusCreated, response.Success(user, http.StatusCreated, "User Registered"))
}