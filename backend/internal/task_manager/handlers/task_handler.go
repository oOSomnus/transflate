package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/internal/task_manager/usecase"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"net/http"
	"path/filepath"
)

// init initializes the log package with specific flags and a custom prefix for task handler logging.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task handler] ")
}

// TaskHandler defines an interface for handling task-related operations.
// TaskSubmit processes the submission of a task from the request context.
// TaskStatusCheckHandler retrieves the status of a task based on the request context.
type TaskHandler interface {
	TaskSubmit(c *gin.Context)
	TaskStatusCheckHandler(c *gin.Context)
}

// TaskHandlerImpl handles task-related operations, connecting the use case and task status service layers.
type TaskHandlerImpl struct {
	Usecase           usecase.TaskUsecase
	TaskStatusService service.TaskStatusService
}

// NewTaskHandler initializes and returns a new instance of TaskHandlerImpl with the provided usecase and service.
func NewTaskHandler(u usecase.TaskUsecase, tss service.TaskStatusService) *TaskHandlerImpl {
	return &TaskHandlerImpl{Usecase: u, TaskStatusService: tss}
}

// TaskSubmit handles the submission of a task, including file upload, processing, status updates, and download link generation.
func (h *TaskHandlerImpl) TaskSubmit(c *gin.Context) {
	log.Println("Processing new task submission...")

	usernameStr, err := getAuthenticatedUsername(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, "User not authorized to submit task")
		return
	}

	fileContent, fileName, err := handleFileUpload(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	lang := c.DefaultPostForm("lang", "eng")

	taskId, err := h.TaskStatusService.CreateNewTask(usernameStr, fileName)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create new task")
		return
	}

	log.Printf("Created new task with ID %s", taskId)
	c.JSON(http.StatusOK, gin.H{"data": "ok"})
	go func() {
		err = h.TaskStatusService.UpdateTaskStatus(usernameStr, taskId, service.Translating)
		if err != nil {
			log.Printf(
				"Error updating task status: %v", err,
			)
			handleTaskStatusError(usernameStr, taskId, h.TaskStatusService)
			return
		}
		// Process OCR and Translation
		transResponse, err := h.Usecase.ProcessOCRAndTranslate(usernameStr, fileContent, lang)
		if err != nil {
			log.Printf("Error processing OCR and translation: %v", err)
			handleTaskStatusError(usernameStr, taskId, h.TaskStatusService)
			return
		}
		err = h.TaskStatusService.UpdateTaskStatus(usernameStr, taskId, service.Uploading)
		if err != nil {
			log.Printf("Error updating task status: %v", err)
			handleTaskStatusError(usernameStr, taskId, h.TaskStatusService)
			return
		}
		// Create download link
		downLink, err := h.Usecase.CreateDownloadLinkWithMdString(transResponse)
		if err != nil {
			log.Printf("Error generating download link: %v", err)
			handleTaskStatusError(usernameStr, taskId, h.TaskStatusService)
			return
		}
		err = h.TaskStatusService.UpdateTaskStatus(usernameStr, taskId, service.Done)
		if err != nil {
			log.Printf("Error updating task status: %v", err)
			handleTaskStatusError(usernameStr, taskId, h.TaskStatusService)
			return
		}

		if err = h.TaskStatusService.UpdateTaskDownloadLink(taskId, downLink); err != nil {
			log.Printf("Error updating task download link: %v", err)
			handleTaskStatusError(usernameStr, taskId, h.TaskStatusService)
			return
		}
	}()
}

// getAuthenticatedUsername retrieves the authenticated username from the given context.
// Returns an error if the username is not found or is of an invalid type.
func getAuthenticatedUsername(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", fmt.Errorf("username not found")
	}
	usernameStr, ok := username.(string)
	if !ok {
		return "", fmt.Errorf("invalid username type")
	}
	return usernameStr, nil
}

// handleFileUpload processes a file upload from a multipart form, ensuring it is a valid PDF file and returns its content.
// Returns an error if the file is missing, not a PDF, or fails during content reading.
func handleFileUpload(c *gin.Context) ([]byte, string, error) {
	file, err := c.FormFile("document")
	if err != nil {
		return nil, "", fmt.Errorf("invalid document")
	}

	if filepath.Ext(file.Filename) != ".pdf" {
		return nil, "", fmt.Errorf("only PDF files are allowed")
	}

	fileContent, err := utils.OpenFile(file)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read file content")
	}
	return fileContent, file.Filename, nil
}

// handleError sends a JSON error response with the given HTTP status code and message, and logs the error message.
func handleError(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.JSON(statusCode, gin.H{"error": message})
}

// handleTaskStatusError updates the task status to error state and logs any error encountered during the update process.
func handleTaskStatusError(username string, taskId string, taskHandler service.TaskStatusService) {
	err := taskHandler.UpdateTaskStatus(username, taskId, service.Error)
	if err != nil {
		log.Printf("Error updating task status: %v", err)
		return
	}
}

// TaskStatusCheckHandler handles the retrieval of all tasks for an authenticated user and returns the results as JSON.
// Responds with 401 if the user is unauthorized or 500 if there is an error retrieving the tasks.
// If successful, responds with 200 and the task data in the response body.
func (h *TaskHandlerImpl) TaskStatusCheckHandler(c *gin.Context) {
	log.Println("Processing Task Status Check...")

	usernameStr, err := getAuthenticatedUsername(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, "User not authorized to submit task")
		return
	}
	results, err := h.TaskStatusService.GetAllTask(usernameStr)
	if err != nil {
		log.Printf("Error getting all task: %v", err)
		handleError(c, http.StatusInternalServerError, "Failed to get all task")
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}
