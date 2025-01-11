package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/cmd/task_manager/config"
	"github.com/oOSomnus/transflate/internal/task_manager/handlers"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/internal/task_manager/usecase"
	"github.com/oOSomnus/transflate/pkg/middleware"
	"github.com/spf13/viper"
	"log"
	"os"
)

const (
	defaultPort        = ":8080"
	defaultEnvironment = "local"
	configType         = "yaml"
	pgUsernameKey      = "pg.username"
	pgPasswordKey      = "pg.password"
	corsAllowOriginKey = "cors.allow-origin"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task Manager Service] ")
}

func main() {
	initializeConfig()
	r := initializeServer()
	log.Printf("Starting server on %s", defaultPort)
	if err := r.Run(defaultPort); err != nil {
		log.Fatal(err)
	}
}

// initializeServer handles all the setup logic for the server
func initializeServer() *gin.Engine {
	gin.SetMode(viper.GetString("gin.mode"))

	// Check database configuration
	if viper.GetString(pgUsernameKey) == "" || viper.GetString(pgPasswordKey) == "" {
		log.Fatalf("Database username or password is missing in the configuration")
	}

	// Establish database connection
	pgConfig := config.NewPostgresConfig()
	dbConnection, err := pgConfig.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	r := gin.Default()
	setupTrustedProxies(r)
	setupCORS(r)

	// Initialize repositories, use cases, and handlers
	userRepo := repository.NewUserRepository(dbConnection)
	userHandler := handlers.NewUserHandler(usecase.NewUserUsecase(userRepo))
	taskHandler := handlers.NewTaskHandler(usecase.NewTaskUsecase(userRepo))

	// Register routes
	r.POST("/login", userHandler.Login)
	r.POST("/register", userHandler.Register)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/submit", taskHandler.TaskSubmit)
	auth.GET("/user/info", userHandler.Info)

	// Handle closing resources
	closeServices(dbConnection)

	return r
}

// initializeConfig sets up the configuration based on the environment
func initializeConfig() {
	environment := os.Getenv("TRANSFLATE_ENV")
	if environment == "" {
		environment = defaultEnvironment
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", environment))
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}

// setupTrustedProxies configures trusted proxies for the server
func setupTrustedProxies(r *gin.Engine) {
	trustedProxies := []string{"172.18.0.0/16", "127.0.0.1"}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}
}

// setupCORS configures CORS settings
func setupCORS(r *gin.Engine) {
	allowOrigin := viper.GetString(corsAllowOriginKey)
	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{allowOrigin},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
			},
		),
	)
	log.Printf("CORS setting: %s", allowOrigin)
}

// closeServices ensures all resources are properly closed on shutdown
func closeServices(db *sql.DB) {
	defer func() {
		if db != nil && db.Close() != nil {
			log.Println("Error closing database connection")
		}
	}()

	defer func() {
		if err := service.CloseOcrGRPCConn(); err != nil {
			log.Println(err)
		}
	}()

	defer func() {
		if err := service.CloseTransGrpcConn(); err != nil {
			log.Println(err)
		}
	}()
}
