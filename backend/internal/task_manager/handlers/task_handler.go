package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/internal/task_manager/usecase"
	"github.com/oOSomnus/transflate/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
	Usecase usecase.TaskUsecase
}

// NewTaskHandler creates and returns a new instance of TaskHandlerImpl with the provided TaskUsecase instance.
func NewTaskHandler(u usecase.TaskUsecase) *TaskHandlerImpl {
	return &TaskHandlerImpl{Usecase: u}
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
	downLink, err := CreateDownloadLinkWithMdString(transResponse)
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

// s3KeyPrefix defines the prefix used for S3 object keys.
// tempFilePattern specifies the pattern for temporary file naming.
// presignedURLExpiry sets the expiration duration for presigned URLs.
const (
	s3KeyPrefix        = "mds/"
	tempFilePattern    = "respMd-*.md"
	presignedURLExpiry = time.Hour
)

// CreateDownloadLinkWithMdString generates a presigned S3 download link for a Markdown string by uploading it as a file.
func CreateDownloadLinkWithMdString(mdString string) (string, error) {
	bucketName := viper.GetString("s3.bucket.name")

	mdTmpFile, err := createTempFileWithContent(mdString, tempFilePattern)
	if err != nil {
		return "", errors.Wrap(err, "error creating temp file with content")
	}
	defer func() {
		err := os.Remove(mdTmpFile.Name())
		if err != nil {
			log.Printf("failed to remove temp file: %v", err)
		}
	}()

	s3Key := s3KeyPrefix + filepath.Base(mdTmpFile.Name())
	if err := service.UploadFileToS3(bucketName, s3Key, mdTmpFile.Name(), 1); err != nil {
		return "", errors.Wrap(err, "failed to upload file to S3")
	}

	downLink, err := service.GeneratePresignedURL(bucketName, s3Key, presignedURLExpiry)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate presigned URL")
	}

	return downLink, nil
}

// createTempFileWithContent creates a temporary file with the specified content and pattern, returning the file or an error.
func createTempFileWithContent(content, pattern string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, errors.Wrap(err, "error creating temp file")
	}

	// Write content to the file
	if _, err := tempFile.Write([]byte(content)); err != nil {
		tempFile.Close()
		return nil, errors.Wrap(err, "error writing to temp file")
	}

	// Ensure the file is returned in a closed state
	if err := tempFile.Close(); err != nil {
		return nil, errors.Wrap(err, "error closing temp file")
	}

	return tempFile, nil
}
