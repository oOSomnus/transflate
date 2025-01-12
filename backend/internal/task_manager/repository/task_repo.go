package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9" // Redis v9 Go library
)

type TaskRepository interface {
	GetTaskState(ctx context.Context, username, taskId string) (int, error)
	SetTaskState(ctx context.Context, username, taskId string, state int, ttl time.Duration) error
	FetchAllTask(ctx context.Context, username string) (map[string]int, error)
}

type TaskRepositoryImpl struct {
	redisClient *redis.Client
}

// NewTaskRepository creates a new instance of TaskRepository.
func NewTaskRepository(client *redis.Client) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{redisClient: client}
}

// GetTaskState retrieves the integer value of a task's status from Redis by username and taskId.
func (t *TaskRepositoryImpl) GetTaskState(ctx context.Context, username, taskId string) (int, error) {
	key := "task:" + username
	status, err := t.redisClient.HGet(ctx, key, taskId).Result()
	if errors.Is(err, redis.Nil) {
		return 0, errors.New("task does not exist")
	} else if err != nil {
		return 0, err
	}

	state, err := strconv.Atoi(status)
	if err != nil {
		return 0, errors.New("failed to parse task state as an integer")
	}

	return state, nil
}

// SetTaskState sets the integer value of a task's status in Redis for a given username and taskId.
func (t *TaskRepositoryImpl) SetTaskState(
	ctx context.Context, username, taskId string, state int, ttl time.Duration,
) error {
	key := "task:" + username
	err := t.redisClient.HSet(ctx, key, taskId, state).Err()
	if err != nil {
		return err
	}

	// Set expiration for the entire key
	err = t.redisClient.Expire(ctx, key, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

// FetchAllTask retrieves all tasks and their statuses for a given username.
func (t *TaskRepositoryImpl) FetchAllTask(ctx context.Context, username string) (map[string]int, error) {
	key := "task:" + username
	result, err := t.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	tasks := make(map[string]int)
	for taskId, status := range result {
		state, err := strconv.Atoi(status)
		if err != nil {
			return nil, errors.New("failed to parse task state as an integer")
		}
		tasks[taskId] = state
	}

	return tasks, nil
}
