package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9" // Redis v9 Go library
)

type TaskRepository interface {
	GetTaskState(ctx context.Context, uuid string) (int, error)
	SetTaskState(ctx context.Context, uuid string, state int, ttl time.Duration) error
}

type TaskRepositoryImpl struct {
	redisClient *redis.Client
}

// NewTaskRepository creates a new instance of TaskRepository.
func NewTaskRepository(client *redis.Client) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{redisClient: client}
}

// GetTaskState retrieves the integer value of a task from Redis by its UUID.
func (t *TaskRepositoryImpl) GetTaskState(ctx context.Context, uuid string) (int, error) {
	val, err := t.redisClient.Get(ctx, uuid).Result()
	if errors.Is(err, redis.Nil) {
		return 0, errors.New("task state does not exist")
	} else if err != nil {
		return 0, err
	}

	// Convert the retrieved value to an integer
	state, err := strconv.Atoi(val)
	if err != nil {
		return 0, errors.New("failed to parse task state as an integer")
	}

	return state, nil
}

func (t *TaskRepositoryImpl) SetTaskState(ctx context.Context, uuid string, state int, ttl time.Duration) error {
	err := t.redisClient.SetNX(ctx, uuid, state, ttl).Err()
	if err != nil {
		if !errors.Is(err, redis.Nil) {
			return err
		}
		// Key exists, update the value
		err = t.redisClient.Set(ctx, uuid, state, ttl).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
