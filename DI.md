```go
// db.go
package config

import (
	"database/sql"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/spf13/viper"
	"log"
	"github.com/go-redis/redis/v8"
)

type DatabaseConfig interface {
	Connect() (*sql.DB, error)
}

type RedisConfig interface {
	Connect() (*redis.Client, error)
}

type PostgresConfig struct{}

type RedisConfigImpl struct{}

func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{}
}

func (p *PostgresConfig) Connect() (*sql.DB, error) {
	username := viper.GetString("pg.username")
	password := viper.GetString("pg.password")
	host := viper.GetString("pg.host")
	port := "5432"
	dbname := "postgres"
	sslmode := "disable"

	dsn := "postgres://" + username + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=" + sslmode
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to the PostgreSQL database!")
	return db, nil
}

func NewRedisConfig() *RedisConfigImpl {
	return &RedisConfigImpl{}
}

func (r *RedisConfigImpl) Connect() (*redis.Client, error) {
	addr := viper.GetString("redis.addr")
	password := viper.GetString("redis.password")
	db := viper.GetInt("redis.db")

	client := redis.NewClient(
		&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}
	)

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	log.Println("Connected to the Redis server!")
	return client, nil
}

// repository/user_repo.go
package repository

import (
"context"
"database/sql"
"errors"
"time"
"github.com/go-redis/redis/v8"
)
type UserRepository interface {
	FindUserByUsername(username string) (string, error)
	IfUserExists(username string) (bool, error)
	CreateUser(username, password string) error
	GetBalance(username string) (int, error)
	DecreaseBalance(username string, balance int) error
}

type UserRepositoryImpl struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewUserRepository(db *sql.DB, redis *redis.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{DB: db, Redis: redis}
}

// Other repository methods remain unchanged

// usecase/user_usecase.go
package usecase

import (
"github.com/oOSomnus/transflate/internal/task_manager/repository"
"github.com/go-redis/redis/v8"
"golang.org/x/crypto/bcrypt"
)
type UserUsecase struct {
	Repo  repository.UserRepository
	Redis *redis.Client
}

func NewUserUsecase(repo repository.UserRepository, redis *redis.Client) *UserUsecase {
	return &UserUsecase{Repo: repo, Redis: redis}
}

func (u *UserUsecase) Authenticate(username, password string) (bool, error) {
	pwdHash, err := u.Repo.FindUserByUsername(username)
	if err != nil {
		return false, err
	}

	if bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password)) != nil {
		return false, errors.New("invalid credentials")
	}
	return true, nil
}

func (u *UserUsecase) CacheUserSession(username string, sessionData string) error {
	ctx := context.Background()
	err := u.Redis.Set(ctx, "session:"+username, sessionData, 24*time.Hour).Err()
	return err
}

// handlers/user_handler.go
package handlers

import (
"github.com/gin-gonic/gin"
"github.com/oOSomnus/transflate/internal/task_manager/usecase"
"net/http"
)
type UserHandler struct {
	Usecase *usecase.UserUsecase
}

func NewUserHandler(u *usecase.UserUsecase) *UserHandler {
	return &UserHandler{Usecase: u}
}

func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	isAuthenticated, err := h.Usecase.Authenticate(req.Username, req.Password)
	if err != nil || !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	h.Usecase.CacheUserSession(req.Username, "session_data_here")

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// main.go
package main

import (
"github.com/gin-gonic/gin"
"github.com/oOSomnus/transflate/cmd/task_manager/config"
"github.com/oOSomnus/transflate/internal/task_manager/handlers"
"github.com/oOSomnus/transflate/internal/task_manager/repository"
"github.com/oOSomnus/transflate/internal/task_manager/usecase"
)

func main() {
	dbConfig := config.NewPostgresConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		panic(err)
	}

	redisConfig := config.NewRedisConfig()
	redisClient, err := redisConfig.Connect()
	if err != nil {
		panic(err)
	}

	repo := repository.NewUserRepository(db, redisClient)
	usecase := usecase.NewUserUsecase(repo, redisClient)
	handler := handlers.NewUserHandler(usecase)

	r := gin.Default()
	r.POST("/login", handler.Login)
	// 注册其他路由...

	r.Run(":8080")

	---
	package
	main

	import (
		"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/cmd/task_manager/config"
	"github.com/oOSomnus/transflate/internal/task_manager/handlers"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"github.com/oOSomnus/transflate/internal/task_manager/usecase"
	)

	type App struct {
		Router   *gin.Engine
		DB       *sql.DB
		Redis    *redis.Client
		Handlers struct {
			UserHandler *handlers.UserHandler
			// 其他 Handler...
		}
	}

	func
	NewApp()(*App, error)
	{
		dbConfig := config.NewPostgresConfig()
		db, err := dbConfig.Connect()
		if err != nil {
			return nil, err
		}

		redisConfig := config.NewRedisConfig()
		redisClient, err := redisConfig.Connect()
		if err != nil {
			return nil, err
		}

		repo := repository.NewUserRepository(db, redisClient)
		usecase := usecase.NewUserUsecase(repo, redisClient)

		app := &App{
			Router: gin.Default(),
			DB:     db,
			Redis:  redisClient,
		}

		app.Handlers.UserHandler = handlers.NewUserHandler(usecase)
		// 初始化其他 Handler...

		return app, nil
	}

	func
	main()
	{
		app, err := NewApp()
		if err != nil {
			panic(err)
		}

		r := app.Router
		r.POST("/login", app.Handlers.UserHandler.Login)
		// 注册其他路由...

		r.Run(":8080")
	}

}


```