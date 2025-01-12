// redis.go
package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// RedisService defines methods for interacting with a Redis client and managing its lifecycle.
type RedisService interface {
	GetClient() *redis.Client
	GetContext() context.Context
	Close()
}

// RedisClient represents a wrapper around *redis.Client and context.Context for interacting with a Redis database.
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisClient initializes a new RedisClient instance with connection details from configuration and tests the connection.
func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			DB:       0, // 使用默认数据库
		},
	)

	// Test connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic("Failed to connect to redis: " + err.Error())
	}

	return &RedisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}

// GetClient returns the Redis client instance managed by the RedisClient structure.
func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

// GetContext returns the context associated with the RedisClient instance.
func (r *RedisClient) GetContext() context.Context {
	return r.ctx
}

// Close terminates the connection to the Redis server held by the RedisClient.
func (r *RedisClient) Close() {
	r.client.Close()
}
