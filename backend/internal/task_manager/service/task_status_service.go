package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"log"
	"strings"
	"time"
)

// TaskStatusService provides methods to manage and query the status of tasks, including updating, retrieving, and creating tasks.
type TaskStatusService interface {
	UpdateTaskStatus(string, string, int) error
	GetTaskStatus(string, string) (int, error)
	CreateNewTask(string) (string, error)
	GetAllTask(string) (map[string]map[string]interface{}, error)
	UpdateTaskDownloadLink(string, string, string) error
}

// TaskStatusServiceImpl provides methods to manage task states via a TaskRepository.
// It includes functionalities for creating, updating, retrieving, and fetching tasks.
type TaskStatusServiceImpl struct {
	tr repository.TaskRepository
}

// NewTaskStatusService initializes and returns a new instance of TaskStatusServiceImpl with the provided TaskRepository.
func NewTaskStatusService(tr repository.TaskRepository) *TaskStatusServiceImpl {
	return &TaskStatusServiceImpl{tr: tr}
}

// InvalidTaskId indicates that the provided Task ID is invalid.
// NotAuthorized represents an authorization failure error message.
// ErrorAccessingData signifies an error occurred while accessing data.
const (
	InvalidTaskId      = "invalid Task ID"
	NotAuthorized      = "not Authorized"
	ErrorAccessingData = "error accessing data"
)

// TaskReceived represents the state where a task has been received and not yet processed.
// ProcessingImages indicates that the task is currently processing images.
// ProcesingText signifies that the task is processing textual data.
// Done denotes that the task has been completed successfully.
// Error represents the state where an error occurred in task processing.
const (
	TaskReceived = 0
	Translating  = 1
	Uploading    = 2
	Done         = 3
	Error        = 9
)

// UpdateTaskStatus updates the status of the specified task if the username matches and returns an error if any issue occurs.
func (tss *TaskStatusServiceImpl) UpdateTaskStatus(username string, taskID string, status int) error {
	// taskId format: username-UUID
	idUsername, taskUUID, err := parseTaskID(taskID)
	if err != nil {
		log.Printf("error parsing task id: %v", err)
		return err
	}
	if idUsername != username {
		log.Println("Not Authorized")
		return errors.New(NotAuthorized)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = tss.tr.SetTaskState(ctx, idUsername, taskUUID, status, 12*time.Hour)
	if err != nil {
		log.Printf("Error handling update: %v", err)
		return errors.New(ErrorAccessingData)
	}
	return nil
}

// GetTaskStatus retrieves the status of a specific task for a given username and task ID.
// Returns the task status as an integer and an error if any occurs during the operation.
func (tss *TaskStatusServiceImpl) GetTaskStatus(username string, taskID string) (int, error) {
	idUsername, taskUUID, err := parseTaskID(taskID)
	if err != nil {
		log.Printf("error parsing task id: %v", err)
		return 0, err
	}
	if idUsername != username {
		log.Println("Not Authorized")
		return 0, errors.New(NotAuthorized)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	status, _, err := tss.tr.GetTaskState(ctx, idUsername, taskUUID)
	if err != nil {
		log.Printf("Error handling update: %v", err)
		return 0, errors.New(ErrorAccessingData)
	}
	return status, nil
}

// CreateNewTask generates a new task ID for the specified username, initializes its state, and returns the new task ID or an error.
func (tss *TaskStatusServiceImpl) CreateNewTask(username string) (string, error) {
	newId := uuid.New()
	taskId := username + "-" + newId.String()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := tss.tr.SetTaskState(ctx, username, newId.String(), TaskReceived, 12*time.Hour)
	if err != nil {
		log.Printf("Error handling update: %v", err)
		return "", errors.New(ErrorAccessingData)
	}
	return taskId, nil
}

// parseTaskID splits a task ID into its username and UUID components, returning an error if the format is invalid.
func parseTaskID(taskID string) (string, string, error) {
	elems := strings.Split(taskID, "-")
	if len(elems) != 2 {
		return "", "", errors.New(InvalidTaskId)
	}
	return elems[0], elems[1], nil
}

// GetAllTask fetches all tasks along with their status for the specified username and returns them as a map.
func (tss *TaskStatusServiceImpl) GetAllTask(username string) (map[string]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	allTasks, err := tss.tr.FetchAllTask(ctx, username)
	if err != nil {
		log.Printf("Error fetching all tasks: %v", err)
		return nil, errors.New(ErrorAccessingData)
	}
	return allTasks, nil
}

// UpdateTaskDownloadLink updates the download link of a specific task for the given username and task ID. Returns an error if failed.
func (tss *TaskStatusServiceImpl) UpdateTaskDownloadLink(username string, taskID string, link string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := tss.tr.UpdateTaskLink(ctx, username, taskID, link); err != nil {
		log.Printf("Error updating task link: %v", err)
		return errors.New(ErrorAccessingData)
	}
	return nil
}
