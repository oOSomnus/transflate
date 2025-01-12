package repository

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9" // Redis v9 Go library
)

type TaskRepository interface {
	GetTaskState(ctx context.Context, username, taskId string) (int, string, error)
	SetTaskState(ctx context.Context, username, taskId string, state int, ttl time.Duration) error
	FetchAllTask(ctx context.Context, username string) (map[string]map[string]interface{}, error)
	UpdateTaskLink(ctx context.Context, username, taskId, link string) error
}

type TaskRepositoryImpl struct {
	redisClient *redis.Client
}

// NewTaskRepository creates a new instance of TaskRepository.
func NewTaskRepository(client *redis.Client) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{redisClient: client}
}

// GetTaskState retrieves the integer value of a task's status and link from Redis by username and taskId.
func (t *TaskRepositoryImpl) GetTaskState(ctx context.Context, username, taskId string) (int, string, error) {
	key := "task:" + username
	rawData, err := t.redisClient.HGet(ctx, key, taskId).Result()
	if errors.Is(err, redis.Nil) {
		return 0, "", errors.New("task does not exist")
	} else if err != nil {
		return 0, "", err
	}

	var taskData map[string]interface{}
	if err := json.Unmarshal([]byte(rawData), &taskData); err != nil {
		return 0, "", errors.New("failed to parse task data")
	}

	state, err := strconv.Atoi(taskData["status"].(string))
	if err != nil {
		return 0, "", errors.New("failed to parse task state as an integer")
	}

	link := ""
	if l, ok := taskData["link"].(string); ok {
		link = l
	}

	return state, link, nil
}

// SetTaskState sets the integer value of a task's status in Redis for a given username and taskId.
func (t *TaskRepositoryImpl) SetTaskState(
	ctx context.Context, username, taskId string, state int, ttl time.Duration,
) error {
	key := "task:" + username
	taskData := map[string]interface{}{
		"status": strconv.Itoa(state),
		"link":   "",
	}
	rawData, err := json.Marshal(taskData)
	if err != nil {
		return errors.New("failed to serialize task data")
	}

	err = t.redisClient.HSet(ctx, key, taskId, rawData).Err()
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

// FetchAllTask retrieves all tasks and their statuses and links for a given username.
func (t *TaskRepositoryImpl) FetchAllTask(ctx context.Context, username string) (
	map[string]map[string]interface{}, error,
) {
	key := "task:" + username
	result, err := t.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	tasks := make(map[string]map[string]interface{})
	for taskId, rawData := range result {
		var taskData map[string]interface{}
		if err := json.Unmarshal([]byte(rawData), &taskData); err != nil {
			return nil, errors.New("failed to parse task data")
		}
		tasks[taskId] = taskData
	}

	return tasks, nil
}

// UpdateTaskLink updates the link of a task in Redis for a given username and taskId.
func (t *TaskRepositoryImpl) UpdateTaskLink(ctx context.Context, username, taskId, link string) error {
	key := "task:" + username
	rawData, err := t.redisClient.HGet(ctx, key, taskId).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("task does not exist")
	} else if err != nil {
		return err
	}

	var taskData map[string]interface{}
	if err := json.Unmarshal([]byte(rawData), &taskData); err != nil {
		return errors.New("failed to parse task data")
	}

	taskData["link"] = link
	updatedData, err := json.Marshal(taskData)
	if err != nil {
		return errors.New("failed to serialize updated task data")
	}

	err = t.redisClient.HSet(ctx, key, taskId, updatedData).Err()
	if err != nil {
		return err
	}

	return nil
}
