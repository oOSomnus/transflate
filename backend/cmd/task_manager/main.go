package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/cmd/task_manager/config"
	"github.com/oOSomnus/transflate/internal/task_manager/handlers"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/pkg/middleware"
	"github.com/spf13/viper"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task Manager Service] ")
}

func main() {
	// viper config
	env := os.Getenv("TRANSFLATE_ENV")
	if env == "" {
		env = "local"
	}
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	gin.SetMode(gin.DebugMode)
	config.ConnectDB()

	if config.DB == nil {
		log.Fatal("Database connection failed")
	}

	r := gin.Default()
	r.SetTrustedProxies([]string{"172.18.0.0/16", "localhost"})
	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{viper.GetString("cors.allow-origin")},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
			},
		),
	)
	log.Printf("cors setting: %s", viper.GetString("cors.allow-origin"))
	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/submit", handlers.TaskSubmit)
	auth.GET("/user/info", handlers.Info)
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
