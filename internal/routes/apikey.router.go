package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/controllers"
)

func ApiKeyRouter(r *gin.RouterGroup) {
	r.POST("/create", controllers.GenerateNewApiKey)
}