package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/controllers"
)

func Router(router *gin.Engine) {
	// Parent routes
	api := router.Group("/api/v1")
	api.GET("/health", controllers.HealthCheck)
	
	// User Routes
	userGroup := api.Group("/users")
	UserRoutes(userGroup)

	// Email Routes
	emailGroup := api.Group("/email")
	EmailRoutes(emailGroup)

	// ApiKey Routes
	apiKeyGroup := api.Group("/apikey")
	ApiKeyRouter(apiKeyGroup)

	// Organization Routes
	organizationGroups := api.Group("/orgs")
	OrganizationRoutes(organizationGroups)
}