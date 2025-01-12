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

// TaskStatusService represents a service for managing and retrieving task statuses.
// Provides methods to update, retrieve, and create task statuses.
type TaskStatusService interface {
	UpdateTaskStatus(string, string, int) error
	GetTaskStatus(string, string) (int, error)
	CreateNewTask(string) (string, error)
}

// TaskStatusServiceImpl is a service implementation that interacts with TaskRepository to manage task states.
type TaskStatusServiceImpl struct {
	tr repository.TaskRepository
}

// InvalidTaskId represents the error message for an invalid task ID.
// NotAuthorized indicates the error message when the user is not authorized.
// ErrorAccessingData is the error message for issues encountered while accessing data.
const (
	InvalidTaskId      = "invalid Task ID"
	NotAuthorized      = "not Authorized"
	ErrorAccessingData = "error accessing data"
)

// TaskReceived represents the initial state when a task is received.
// ProcessingImages represents the state where the task involves image processing.
// ProcesingText represents the state where the task involves text processing.
// Done represents the state when the task is completed.
const (
	TaskReceived     = 0
	ProcessingImages = 1
	ProcesingText    = 2
	Done             = 3
)

// UpdateTaskStatus updates the status of a task for a given username and taskID.
// Returns an error if the taskID is invalid, unauthorized, or if there is an issue accessing data.
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
	qkey := getQueryKey(username, taskUUID)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = tss.tr.SetTaskState(ctx, qkey, status, 12*time.Hour)
	if err != nil {
		log.Printf("Error handling update: %v", err)
		return errors.New(ErrorAccessingData)
	}
	return nil
}

// GetTaskStatus retrieves the status of a task for a given username and taskID.
// Returns an integer status code and an error if any issues occur during execution.
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
	qkey := getQueryKey(username, taskUUID)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	status, err := tss.tr.GetTaskState(ctx, qkey)
	if err != nil {
		log.Printf("Error handling update: %v", err)
		return 0, errors.New(ErrorAccessingData)
	}
	return status, nil
}

// CreateNewTask generates a new task ID by combining the provided username with a newly generated UUID.
func (tss *TaskStatusServiceImpl) CreateNewTask(username string) string {
	newId := uuid.New()
	return username + "-" + newId.String()
}

// parseTaskID splits the provided taskID string into username and task UUID. It returns an error if the format is invalid.
func parseTaskID(taskID string) (string, string, error) {
	elems := strings.Split(taskID, "-")
	if len(elems) != 2 {
		return "", "", errors.New(InvalidTaskId)
	}
	return elems[0], elems[1], nil
}

// getQueryKey generates a unique query key based on the username and the task UUID.
func getQueryKey(username string, taskUUID string) string {
	return "task" + ":" + username + ":" + taskUUID
}
