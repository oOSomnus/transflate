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
	"strings"
)

// defaultPort defines the default port the server will run on.
// defaultEnvironment specifies the default environment for the application.
// configType indicates the configuration file type used.
// corsAllowOriginKey represents the key for allowed CORS origins in configurations.
// pgUsernameKey defines the key for the PostgreSQL username in configurations.
// pgPasswordKey defines the key for the PostgreSQL password in configurations.
// allowMethods specifies the HTTP methods allowed for CORS.
// allowHeaders specifies the HTTP headers allowed for CORS.
const (
	defaultPort        = ":8080"
	defaultEnvironment = "local"
	configType         = "yaml"
	//pgPort             = "5432"
	//pgDBName           = "postgres"
	corsAllowOriginKey = "cors.allow-origin"
	pgUsernameKey      = "pg.username"
	pgPasswordKey      = "pg.password"
	allowMethods       = "GET,POST,PUT,DELETE,OPTIONS"
	allowHeaders       = "Content-Type,Authorization"
)

// init configures the logger with standard flags, microseconds precision, and a custom prefix for the service.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task Manager Service] ")
}

// main is the entry point of the application where configuration is initialized, resources are set up, and the server starts.
func main() {
	initializeConfig()

	// Setup resources
	dbConnection := setupDatabaseConnection()
	defer cleanupResources(dbConnection)
	redisClient := setupRedisClient()
	defer redisClient.Close()

	r := initializeServer(dbConnection, redisClient)

	log.Printf("Starting server on %s", defaultPort)
	if err := r.Run(defaultPort); err != nil {
		log.Fatal(err)
	}
}

// setupDatabaseConnection initializes and returns a database connection using the Postgres configuration settings.
func setupDatabaseConnection() *sql.DB {
	pgConfig := config.NewPostgresConfig()
	dbConnection, err := pgConfig.Connect()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	return dbConnection
}

// setupRedisClient initializes and returns a new instance of RedisClient from the configuration.
func setupRedisClient() *config.RedisClient {
	return config.NewRedisClient()
}

// initializeServer initializes and configures a Gin engine with middleware, repositories, handlers, and routes.
// It takes a database connection and a Redis client as parameters and returns the configured *gin.Engine.
// The function sets trusted proxies, configures CORS, initializes services, and sets up routes for HTTP handling.
// It also ensures proper resource cleanup for services upon termination.
func initializeServer(db *sql.DB, redisClient *config.RedisClient) *gin.Engine {
	gin.SetMode(viper.GetString("gin.mode"))
	verifyDatabaseCredentials()

	r := gin.Default()

	setupTrustedProxies(r, []string{"172.18.0.0/16", "127.0.0.1"})
	configureCORS(r, viper.GetString(corsAllowOriginKey))

	// Setup repositories
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(redisClient.GetClient())

	// Setup handlers
	userHandler := handlers.NewUserHandler(usecase.NewUserUsecase(userRepo))
	s3Service, ocrService, translateService := initializeServices()

	defer cleanupServiceResources(ocrService, translateService)

	taskStatusService := service.NewTaskStatusService(taskRepo)

	taskUsecase := usecase.NewTaskUsecase(userRepo, taskRepo, ocrService, s3Service, translateService)
	taskHandler := handlers.NewTaskHandler(taskUsecase, taskStatusService)

	setupRoutes(r, userHandler, taskHandler)

	return r
}

// initializeServices initializes and returns instances of S3StorageServiceImpl, OCRService, and TranslateServiceImpl.
// It logs and exits the application if any of the services fail to initialize.
func initializeServices() (service.S3StorageService, service.OCRClient, service.TranslateService) {
	s3Service, err := service.NewS3StorageService()
	if err != nil {
		log.Fatalf("S3 service initialization failed: %v", err)
	}

	ocrService := service.NewOCRService()

	translateService, err := service.NewTranslateService()
	if err != nil {
		log.Fatalf("Translate service initialization failed: %v", err)
	}

	return s3Service, ocrService, translateService
}

// cleanupServiceResources ensures the proper closure of resources for OCR and Translate services to release gRPC connections.
func cleanupServiceResources(ocrService service.OCRClient, translateService service.TranslateService) {
	if err := ocrService.Close(); err != nil {
		log.Println("OCR service close error:", err)
	}
	if err := translateService.CloseTransGrpcConn(); err != nil {
		log.Println("Translate service close error:", err)
	}
}

// configureCORS configures Cross-Origin Resource Sharing (CORS) settings for a gin.Engine instance.
// It allows requests from the specified origin and sets HTTP methods, headers, and credentials options.
func configureCORS(r *gin.Engine, allowOrigin string) {
	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins:     []string{allowOrigin},
				AllowMethods:     parseCSV(allowMethods),
				AllowHeaders:     parseCSV(allowHeaders),
				AllowCredentials: true,
			},
		),
	)
	log.Printf("CORS setting: %s", allowOrigin)
}

// setupTrustedProxies configures the trusted proxy IPs for the provided Gin engine instance.
// It terminates the application on error setting the proxies.
func setupTrustedProxies(r *gin.Engine, proxies []string) {
	if err := r.SetTrustedProxies(proxies); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}
}

// setupRoutes initializes HTTP routes and associates them with their corresponding handlers and middleware.
// r is the Gin engine used to define HTTP routes and middleware.
// userHandler handles user-related endpoints like login, register, and user info.
// taskHandler handles task-related endpoints, such as task submission.
func setupRoutes(r *gin.Engine, userHandler *handlers.UserHandlerImpl, taskHandler *handlers.TaskHandlerImpl) {
	r.POST("/login", userHandler.Login)
	r.POST("/register", userHandler.Register)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/submit", taskHandler.TaskSubmit)
	auth.GET("/user/info", userHandler.Info)
}

// verifyDatabaseCredentials ensures the presence of database username and password in the application configuration.
// Logs a fatal error and terminates the application if either credential is missing.
func verifyDatabaseCredentials() {
	if viper.GetString(pgUsernameKey) == "" || viper.GetString(pgPasswordKey) == "" {
		log.Fatalf("Database credentials are missing in the configuration")
	}
}

// initializeConfig reads and applies environment-specific configuration files using Viper.
// Defaults to "local" environment if TRANSFLATE_ENV is not set.
// Terminates the application if the configuration file cannot be read.
func initializeConfig() {
	env := os.Getenv("TRANSFLATE_ENV")
	if env == "" {
		env = defaultEnvironment
	}

	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType(configType)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Config file reading error: %v", err)
	}
}

// cleanupResources safely closes the provided database connection and logs any errors that occur during the close operation.
func cleanupResources(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Println("Database closing error:", err)
	}
}

// parseCSV splits a comma-separated string into a slice of strings.
func parseCSV(input string) []string {
	return strings.Split(input, ",")
}
