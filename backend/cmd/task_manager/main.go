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

const (
	defaultPort         = ":8080"
	defaultEnvironment  = "local"
	configType          = "yaml"
	ginModeConfigKey    = "gin.mode"
	corsAllowOriginKey  = "cors.allow-origin"
	jwtSecretConfigKey  = "jwt.secret"
	trustedProxiesCIDR1 = "172.18.0.0/16"
	trustedProxiesCIDR2 = "127.0.0.1"
)

func main() {
	initializeConfig()
	initializeServerDependencies()

	r := gin.Default()
	setupTrustedProxies(r)
	setupCORS(r)

	registerRoutes(r)

	log.Printf("Starting server on %s", defaultPort)
	deferResourceCleanup()
	if err := r.Run(defaultPort); err != nil {
		log.Fatal(err)
	}
}

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

func initializeServerDependencies() {
	gin.SetMode(viper.GetString(ginModeConfigKey))
	config.ConnectDB()
	if config.DB == nil {
		log.Fatal("Database connection failed")
	}
}

func setupTrustedProxies(r *gin.Engine) {
	trustedProxies := []string{trustedProxiesCIDR1, trustedProxiesCIDR2}
	if err := r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Error setting trusted proxies: %v", err)
	}
}

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

func registerRoutes(r *gin.Engine) {
	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.Register)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/submit", handlers.TaskSubmit)
	auth.GET("/user/info", handlers.Info)
}

func deferResourceCleanup() {
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}(config.DB)

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
