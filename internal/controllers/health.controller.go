package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zeni-42/Mhawk/internal/database"
)

func HealthCheck(c *gin.Context) {
	pingPG := database.GetPing()

	m := map[string]interface{}{
		"status": "OK",
		"service": "Mhawk",
		"version": "0.1.0",
		"timestamps": time.Now(),
		"dependencies": map[string]interface{}{
			"Postgres": pingPG, 
		}, 
	}

	c.JSON(200, m)
}