package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zeni-42/Mhawk/internal/routes"
)


func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(".env file not loaded")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routes.Router(r)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT not found")
	}
	fmt.Println("[GIN] SERVER:" + port)
	r.Run(":"+port)
}