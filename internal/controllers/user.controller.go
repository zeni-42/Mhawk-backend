package controllers

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/zeni-42/Mhawk/internal/models"
	"github.com/zeni-42/Mhawk/internal/repository"
	"github.com/zeni-42/Mhawk/internal/utils/cloudinary"
	"github.com/zeni-42/Mhawk/internal/utils/response"
	"github.com/zeni-42/Mhawk/internal/utils/token"
	"golang.org/x/crypto/bcrypt"
)

// Uses to handle user registration and save in DB
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

// Used to handle user login and set cookies
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

	if err := repository.UpdateRefreshToken(registeredUser.Id, Rtoken); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Database error"))
		return
	}

	context.SetCookie("accessToken", AToken, 60 * 60 * 24 * 10, "/", "", true, true)
	context.SetCookie("refreshToken", Rtoken, 60 * 60 * 24 * 30, "/", "", true, true)

	context.IndentedJSON(http.StatusOK, response.Success(nil, http.StatusOK, "User logged in"))
}

type User struct {
	UserId string `json:"userId"`
}

// Used to log out the user, removes refreshtoken in DB and clear cookies
func LogoutUser(context *gin.Context) {
	var user User 

	if err := context.BindJSON(&user); err != nil {
		log.Printf("BindJSON error: %v", err)
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid data"))
		return
	}
	id, err := uuid.Parse(user.UserId)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid UUID format"))
		return
	}

	if err := repository.UpdateRefreshToken(id, nil); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Database error"))
		return
	}

	context.SetCookie("accessToken", "", -1, "/", "", true, true)
	context.SetCookie("refreshToken", "", -1, "/", "", true, true)

	context.IndentedJSON(http.StatusOK, response.Success(nil, http.StatusOK, "User logged out"))
}

// Used to update the default image and if the user wants to update the image later on this function will be used.
func UpdateAvatar(context *gin.Context) {
	var user User
	
	if err := context.BindJSON(&user.UserId); err != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid data"))
		return
	}
	
	id, err := uuid.Parse(user.UserId)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, response.Error(err, http.StatusBadRequest, "Invalid UUID format"))
		return
	}
	
	file, _ := context.FormFile("avatar")
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save avatar"))
		return
	}
	path, _ := os.Getwd()
	publicPath := filepath.Join(path, "public")

	if err := os.MkdirAll(publicPath, os.ModePerm); err != nil {
		log.Println("Failed to create public directory:", err)
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Server error"))
		return
	}

	filename := filepath.Base(file.Filename)
	savePath := filepath.Join(publicPath, filename)

	if err := context.SaveUploadedFile(file, savePath); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save file"))
		return
	}

	avatarUrl, err := cloudinary.UploadOnCloudinary(file)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to upload image"))
		return
	}

	if err := repository.UpdateUserAvatar(id, avatarUrl); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, response.Error(err, http.StatusInternalServerError, "Failed to save image"))
		return
	}

	context.IndentedJSON(http.StatusOK, response.Success(nil, http.StatusOK, "Avatar added"))
}