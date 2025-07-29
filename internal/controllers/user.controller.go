package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/zeni-42/Mhawk/internal/models"
	"github.com/zeni-42/Mhawk/internal/repository"
	"github.com/zeni-42/Mhawk/internal/utils/response"
	"github.com/zeni-42/Mhawk/internal/utils/token"
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
		return
	}

	if existingUser != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Email is taken"))
		return
	}

	byteHashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	user.Password = string(byteHashedPass)

	userId, err := repository.CreateUser(user); 
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save user"))
		return
	}

	userPointer, err := repository.FindUserById(userId)
	if errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save user"))
		return
	}

	user = *userPointer

	context.IndentedJSON(http.StatusCreated, response.Success(user, http.StatusCreated, "User Registered"))
}

func LoginUser(context *gin.Context) {
	var user models.User
	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid data"))
	}

	registeredUser, err := repository.FindUserByEmail(user.Email)
	if errors.Is(err, pgx.ErrNoRows) {
		context.IndentedJSON(http.StatusNotFound, response.Error(nil, http.StatusNotFound, "User not found"))
		return
	}

	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Database error"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(registeredUser.Password), []byte(user.Password)); err != nil {
		context.IndentedJSON(http.StatusUnauthorized, response.Error(nil, http.StatusUnauthorized, "Invalid credentails"))
		return
	}

	uData := *registeredUser

	AToken := token.GetAccessToken(uData)
	Rtoken := token.GetRefreshToken(uData)

	context.SetCookie("accessToken", AToken, 60 * 60 * 24 * 10, "/", "", true, true)
	context.SetCookie("refreshToken", Rtoken, 60 * 60 * 24 * 30, "/", "", true, true)

	context.IndentedJSON(http.StatusOK, response.Success(nil, http.StatusOK, "User logged in"))
}