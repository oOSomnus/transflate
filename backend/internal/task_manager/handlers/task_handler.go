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

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Task handler] ")
}

// TaskSubmit handles the submission of a task by an authenticated user, processes the uploaded PDF file, and generates a downloadable link.
// It validates the session, ensures the uploaded file is a PDF, reads its contents, and performs OCR and translation.
// The result is converted to a markdown-based response and uploaded to generate a presigned download link.
// Responds with the download link or appropriate error in JSON format.
func TaskSubmit(c *gin.Context) {
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
	transResponse, err := usecase.ProcessOCRAndTranslate(usernameStr, fileContent, lang)
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

// getAuthenticatedUsername validates the session and retrieves the username.
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

// handleFileUpload validates the uploaded file, ensuring it's a PDF, and returns its content.
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

// handleError sends a JSON error response with the given status and message.
func handleError(c *gin.Context, statusCode int, message string) {
	log.Println(message)
	c.JSON(statusCode, gin.H{"error": message})
}

// CreateDownloadLinkWithMdString generates a downloadable presigned URL from a given markdown string.
// It writes the input string to a temporary file, uploads the file to S3, and creates a presigned URL for downloading.
const (
	s3KeyPrefix        = "mds/"
	tempFilePattern    = "respMd-*.md"
	presignedURLExpiry = time.Hour
)

// CreateDownloadLinkWithMdString generates a downloadable presigned URL from a given markdown string.
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

// createTempFileWithContent creates a temporary file, writes the provided content, and returns the file.
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
