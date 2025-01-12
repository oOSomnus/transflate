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

// init sets logging configuration with timestamp, microseconds precision, and a prefix for task handler logs.
func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task handler] ")
}

// TaskHandler defines an interface for handling task-related operations, including task submission through HTTP context.
type TaskHandler interface {
	TaskSubmit(c *gin.Context)
}

// TaskHandlerImpl is a struct that implements task-related HTTP handlers by leveraging the TaskUsecase interface.
type TaskHandlerImpl struct {
	Usecase           usecase.TaskUsecase
	TaskStatusService service.TaskStatusService
}

// NewTaskHandler creates and returns a new instance of TaskHandlerImpl with the provided TaskUsecase instance.
func NewTaskHandler(u usecase.TaskUsecase, tss service.TaskStatusService) *TaskHandlerImpl {
	return &TaskHandlerImpl{Usecase: u, TaskStatusService: tss}
}

// TaskSubmit handles task submission by processing an uploaded file, performing OCR and translation, and returning a download link.
func (h *TaskHandlerImpl) TaskSubmit(c *gin.Context) {
	log.Println("Processing new task submission...")

	usernameStr, err := getAuthenticatedUsername(c)
	if err != nil {
		handleError(c, http.StatusUnauthorized, "User not authorized to submit task")
		return
	}

	fileContent, err := handleFileUpload(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error())
		return
	}

	lang := c.DefaultPostForm("lang", "eng")

	// Process OCR and Translation
	transResponse, err := h.Usecase.ProcessOCRAndTranslate(usernameStr, fileContent, lang)
	if err != nil {
		log.Printf("Error processing OCR and translation: %v", err)
		handleError(c, http.StatusBadRequest, "Failed to process OCR")
		return
	}

	// Create download link
	downLink, err := h.Usecase.CreateDownloadLinkWithMdString(transResponse)
	if err != nil {
		log.Printf("Error generating download link: %v", err)
		handleError(c, http.StatusInternalServerError, "Failed to generate download link")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": downLink})
}

// getAuthenticatedUsername retrieves the username of the authenticated user from the given Gin context.
// It returns an error if the username is not found or is of an invalid type.
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

// handleFileUpload processes a file upload from a gin.Context and only accepts PDF files, returning its content as a byte slice.
// Returns an error if the uploaded file is invalid, not a PDF, or if reading the file content fails.
func handleFileUpload(c *gin.Context) ([]byte, error) {
	file, err := c.FormFile("document")
	if err != nil {
		return nil, fmt.Errorf("invalid document")
	}

	if filepath.Ext(file.Filename) != ".pdf" {
		return nil, fmt.Errorf("only PDF files are allowed")
	}

	fileContent, err := utils.OpenFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content")
	}
	return fileContent, nil
}

// handleError sends a JSON response with a given HTTP status code and error message, and logs the error message.
func handleError(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.JSON(statusCode, gin.H{"error": message})
}
