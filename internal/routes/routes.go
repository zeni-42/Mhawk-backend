package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/controllers"
)

func Router(router *gin.Engine) {
	api := router.Group("/api/v1")

	api.GET("/health", controllers.HealthCheck)
}