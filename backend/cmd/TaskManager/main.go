package main

import (
	"github.com/oOSomnus/transflate/services/TaskManager/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/cmd/TaskManager/config"
	"github.com/oOSomnus/transflate/pkg/middleware"
)

func main() {
	config.ConnectDB()

	if config.DB == nil {
		log.Fatal("Database connection failed")
	}

	r := gin.Default()
	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	port := ":8080"
	log.Printf("Starting server on %s", port)
	r.Run(port)

}
