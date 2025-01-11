package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/redis/go-redis/v9" // Redis v9 Go library
)

type TaskRepository interface {
	GetTaskState(ctx context.Context, uuid string) (int64, error)
	SetTaskState(ctx context.Context, uuid string, state int64) error
}

type TaskRepositoryImpl struct {
	redisClient *redis.Client
}

// NewTaskRepository creates a new instance of TaskRepository.
func NewTaskRepository(client *redis.Client) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{redisClient: client}
}

// GetTaskState retrieves the integer state of a task from Redis by its UUID.
func (t *TaskRepositoryImpl) GetTaskState(ctx context.Context, uuid string) (int64, error) {
	val, err := t.redisClient.Get(ctx, uuid).Result()
	if errors.Is(err, redis.Nil) {
		return 0, errors.New("task state does not exist")
	} else if err != nil {
		return 0, err
	}

	// Convert the retrieved value to an integer
	state, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.New("failed to parse task state as an integer")
	}

	return state, nil
}

// SetTaskState sets or updates the integer state of a task in Redis identified by its UUID.
func (t *TaskRepositoryImpl) SetTaskState(ctx context.Context, uuid string, state int64) error {
	err := t.redisClient.Set(ctx, uuid, state, 0).Err() // Setting TTL as 0 means no expiration
	if err != nil {
		return err
	}
	return nil
}
