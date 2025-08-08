package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/controllers"
)

func ApiKeyRouter(r *gin.RouterGroup) {
	r.POST("/", controllers.GenerateNewApiKey)
	r.GET("/:id", controllers.GetUserApiKeys)
	r.DELETE("/:id", controllers.DeleteAPI)
	r.PUT("/:id", controllers.ToggleActive)
}