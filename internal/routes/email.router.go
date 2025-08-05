package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/controllers"
)

func EmailRoutes(r *gin.RouterGroup){
	r.POST("/send", controllers.SendEmail)
}