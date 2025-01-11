// redis.go
package config

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisService interface {
	GetClient() *redis.Client
	GetContext() context.Context
	Close()
}

// RedisClient Impl
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

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

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

func (r *RedisClient) GetContext() context.Context {
	return r.ctx
}

func (r *RedisClient) Close() {
	r.client.Close()
}
