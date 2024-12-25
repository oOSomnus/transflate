package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/cmd/TaskManager/config"
	"github.com/oOSomnus/transflate/cmd/TaskManager/handlers"
)

func main() {
	config.ConnectDB()

	if config.DB == nil {
		log.Fatal("Database connection failed")
	}

	r := gin.Default()
	r.POST("/login", handlers.Login)

	r.Run(":8080")
}
