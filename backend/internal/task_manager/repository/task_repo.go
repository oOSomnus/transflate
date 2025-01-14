package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrTaskNotFound indicates that the requested task was not found or has expired in the data store.
var (
	ErrTaskNotFound = errors.New("task not found or expired")
)

// TaskRepository defines methods for managing and interacting with user tasks and their associated states and metadata.
type TaskRepository interface {
	// SetTaskState: If the filename does not exist, create and set it. If it already exists, the original filename will not be overwritten. Set new status and TTL every time
	SetTaskState(ctx context.Context, username, taskId string, status int, filename string, ttl time.Duration) error

	// GetTaskState: Get the status and filename of task
	GetTaskState(ctx context.Context, username, taskId string) (int, string, error)

	// FetchAllTask: Get all taskIds at once -> {status, filename, link...}, Does not contain expired keys
	FetchAllTask(ctx context.Context, username string) (map[string]map[string]interface{}, error)

	// UpdateTaskLink: Download link for update tasks
	UpdateTaskLink(ctx context.Context, username, taskId, link string) error
}

// RedisTaskRepository interacts with Redis to manage task-related data for users.
// It provides methods to create, update, fetch, and manage task states and details.
// The repository ensures task data is stored and retrieved efficiently using Redis structures.
type RedisTaskRepository struct {
	client *redis.Client
}

// NewRedisTaskRepository initializes and returns a new RedisTaskRepository instance with the provided Redis client.
func NewTaskRepository(client *redis.Client) *RedisTaskRepository {
	return &RedisTaskRepository{
		client: client,
	}
}

// buildTaskKey generates a Redis key for a specific task by combining the username and task ID.
func buildTaskKey(username, taskId string) string {
	return fmt.Sprintf("task:%s:%s", username, taskId)
}

// buildUserSetKey generates a Redis key for the task ID set associated with a specific user based on the username.
func buildUserSetKey(username string) string {
	return fmt.Sprintf("tasks:%s", username)
}

// SetTaskState sets the state of a task, initializing or updating status, filename, and TTL in Redis storage.
func (r *RedisTaskRepository) SetTaskState(
	ctx context.Context, username, taskId string, status int, filename string, ttl time.Duration,
) error {
	key := buildTaskKey(username, taskId)

	// First determine whether the key originally existed
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}

	if exists == 0 {
		// key does not exist -> create and insert filename
		hm := map[string]interface{}{
			"status":   status,
			"filename": filename,
			"link":     "", // 初始空
		}
		if err := r.client.HSet(ctx, key, hm).Err(); err != nil {
			return err
		}
		// Maintains a collection of current user task IDs
		if err := r.client.SAdd(ctx, buildUserSetKey(username), taskId).Err(); err != nil {
			return err
		}
	} else {
		// key already exists -> update status without overwriting the original filename
		// update status
		if err := r.client.HSet(ctx, key, "status", status).Err(); err != nil {
			return err
		}

	}

	// Set separate TTL (from now on)
	if err := r.client.Expire(ctx, key, ttl).Err(); err != nil {
		return err
	}
	return nil
}

// GetTaskState retrieves the status and filename of a task for the given username and taskId from Redis.
// It returns the task's status as an integer, the filename as a string, or an error if the task does not exist or fails retrieval.
func (r *RedisTaskRepository) GetTaskState(ctx context.Context, username, taskId string) (int, string, error) {
	key := buildTaskKey(username, taskId)

	// Determine whether key exists
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return 0, "", err
	}
	if exists == 0 {
		return 0, "", ErrTaskNotFound
	}

	// Get status, filename
	values, err := r.client.HMGet(ctx, key, "status", "filename").Result()
	if err != nil {
		return 0, "", err
	}
	// values[0] -> status, values[1] -> filename
	if len(values) < 2 || values[0] == nil || values[1] == nil {
		return 0, "", ErrTaskNotFound
	}

	statusVal, ok := values[0].(string)
	if !ok {
		return 0, "", ErrTaskNotFound
	}
	filenameVal, ok := values[1].(string)
	if !ok {
		return 0, "", ErrTaskNotFound
	}

	// Convert status to int
	// Note here: What goes into HSet/HMSet is interface{}, and what comes out through HGet is string, which needs to be converted.
	var statusInt int
	_, err = fmt.Sscanf(statusVal, "%d", &statusInt)
	if err != nil {
		return 0, "", err
	}

	return statusInt, filenameVal, nil
}

// FetchAllTask retrieves all tasks for a specified username from Redis and returns them as a nested map structure.
// It checks for task existence, removes expired tasks, and parses task fields into appropriate types. Returns an error if any Redis operation fails.
func (r *RedisTaskRepository) FetchAllTask(ctx context.Context, username string) (
	map[string]map[string]interface{}, error,
) {
	result := make(map[string]map[string]interface{})

	// First get all taskIds under the user name
	taskIds, err := r.client.SMembers(ctx, buildUserSetKey(username)).Result()
	if err != nil {
		return nil, err
	}

	for _, taskId := range taskIds {
		key := buildTaskKey(username, taskId)

		// Determine whether key still exists
		exists, err := r.client.Exists(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if exists == 0 {
			// May have expired or been deleted -> removed from collection
			if remErr := r.client.SRem(ctx, buildUserSetKey(username), taskId).Err(); remErr != nil {
				log.Printf("failed to remove expired task from set: %v", remErr)
			}
			continue
		}

		// Get hash data
		vals, err := r.client.HGetAll(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		if len(vals) == 0 {
			// Description key does not exist or is empty -> remove it altogether
			_ = r.client.SRem(ctx, buildUserSetKey(username), taskId).Err()
			continue
		}

		// Convert fields such as status to appropriate types
		statusInt := 0
		fmt.Sscanf(vals["status"], "%d", &statusInt)

		tmp := map[string]interface{}{
			"status":   statusInt,
			"filename": vals["filename"],
			"link":     vals["link"],
		}
		result[taskId] = tmp
	}

	return result, nil
}

// UpdateTaskLink updates the task's link field for the specified username and taskId in the Redis storage.
// Returns an error if the task is not found or Redis operation fails.
func (r *RedisTaskRepository) UpdateTaskLink(ctx context.Context, username, taskId, link string) error {
	key := buildTaskKey(username, taskId)

	// Determine whether key exists
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return ErrTaskNotFound
	}

	//Direct HSet "link" field
	if err := r.client.HSet(ctx, key, "link", link).Err(); err != nil {
		return err
	}
	return nil
}
