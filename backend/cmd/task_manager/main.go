package main

import (
	"database/sql"
	"github.com/oOSomnus/transflate/internal/task_manager/handlers"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/cmd/task_manager/config"
	"github.com/oOSomnus/transflate/pkg/middleware"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task Manager Service] ")
}

func main() {
	gin.SetMode(gin.DebugMode)
	config.ConnectDB()

	if config.DB == nil {
		log.Fatal("Database connection failed")
	}

	r := gin.Default()
	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/submit", handlers.TaskSubmit)
	port := ":8080"
	log.Printf("Starting server on %s", port)
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(config.DB)

	defer func() {
		err := service.CloseOcrGRPCConn()
		if err != nil {
			log.Println(err)
		}
	}()

	defer func() {
		err := service.CloseTransGrpcConn()
		if err != nil {
			log.Println(err)
		}
	}()

	err := r.Run(port)
	if err != nil {
		log.Fatal(err)
	}

}
