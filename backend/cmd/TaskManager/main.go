package main

import (
	"database/sql"
	"github.com/oOSomnus/transflate/internal/TaskManager/handlers"
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
	r.POST("/submit", handlers.TaskSubmit)
	port := ":8080"
	log.Printf("Starting server on %s", port)
	r.Run(port)
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(config.DB)
}
