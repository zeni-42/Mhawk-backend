package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/zeni-42/Mhawk/internal/database"
	"github.com/zeni-42/Mhawk/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(".env file not loaded")
	}

	database.ConnectPG()
	defer func () {
		if err := database.DisconnectPG(); err != nil {
			log.Println("[GIN] Disconnection failed")
		}
		log.Println("DB DISCONNECTED")
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	routes.Router(r)

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT not found")
	}

	server := &http.Server{
		Addr: ":" + port,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func () {
		log.Println("SERVER:", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println("Failed to start server")
		}
	} ()

	<- stop
	log.Println("Shutdown init")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed")
	}

	log.Println("Server stopped")
}