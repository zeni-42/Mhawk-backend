package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/controllers"
)

func UserRoutes(r *gin.RouterGroup) {
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)
	r.POST("/logout", controllers.LogoutUser)
}