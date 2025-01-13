package usecase

import (
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Placeholder for a const block, typically used for defining multiple constants.
const ()

// TaskUsecase defines methods for processing OCR and translations, as well as generating downloadable links from Markdown.
type TaskUsecase interface {
	ProcessOCRAndTranslate(username string, fileContent []byte, lang string) (string, error)
	CreateDownloadLinkWithMdString(mdString string) (string, error)
}

// TaskUsecaseImpl is the implementation of task-related operations using repository and service dependencies.
// It manages user tasks, integrates OCR, text translation, and S3 storage services.
type TaskUsecaseImpl struct {
	ur   repository.UserRepository
	tr   repository.TaskRepository
	ocrc service.OCRClient
	s3s  service.S3StorageService
	ts   service.TranslateService
}

// NewTaskUsecase initializes and returns a new TaskUsecaseImpl instance with required repositories and services.
func NewTaskUsecase(
	ur repository.UserRepository, tr repository.TaskRepository, ocrc service.OCRClient, s3s service.S3StorageService,
	ts service.TranslateService,
) *TaskUsecaseImpl {
	return &TaskUsecaseImpl{ur: ur, tr: tr, ocrc: ocrc, s3s: s3s, ts: ts}
}

// ProcessOCRAndTranslate performs OCR on the input file, subtracts user balance based on pages, and translates the text.
func (t *TaskUsecaseImpl) ProcessOCRAndTranslate(username string, fileContent []byte, lang string) (string, error) {
	ocrResponse, err := t.ocrc.ProcessOCR(fileContent, lang)
	if err != nil || ocrResponse == nil {
		log.Println("Error during OCR processing:", err)
		return "", errors.New("failed to process OCR")
	}

	// Merge and clean OCR response lines
	cleanedText := mergeAndCleanStrings(ocrResponse.Lines)

	// Decrease user balance based on the number of pages
	numPages := int(ocrResponse.PageNum)
	if err = t.ur.DecreaseBalance(username, numPages); err != nil {
		log.Printf("Error decreasing balance for user %s: %v", username, err)
		return "", err
	}

	// Translate the cleaned text
	translatedResponse, err := t.ts.TranslateText(cleanedText)
	if err != nil {
		log.Println("Error during text translation:", err)
		return "", err
	}

	return translatedResponse.Lines, nil
}

// mergeAndCleanStrings takes a slice of strings, merges them, and cleans the resulting string using text cleaning utils.
func mergeAndCleanStrings(lines []string) string {
	var builder strings.Builder
	for _, line := range lines {
		builder.WriteString(line)
	}
	return utils.TextCleaning(builder.String())
}

// s3KeyPrefix specifies the prefix path for storing objects in the S3 bucket.
// tempFilePattern defines the naming pattern for temporary files used in the application.
// presignedURLExpiry sets the expiration duration for presigned URLs to 1 hour.
const (
	s3KeyPrefix        = "mds/"
	tempFilePattern    = "respMd-*.md"
	presignedURLExpiry = time.Hour
)

// CreateDownloadLinkWithMdString generates a presigned download link for a file created from the provided markdown string.
// It temporarily creates a file with the given markdown content, uploads it to S3, and generates a presigned URL for access.
func (t *TaskUsecaseImpl) CreateDownloadLinkWithMdString(mdString string) (string, error) {
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
	if err := t.s3s.UploadFileToS3(bucketName, s3Key, mdTmpFile.Name(), 1); err != nil {
		return "", errors.Wrap(err, "failed to upload file to S3")
	}

	downLink, err := t.s3s.GeneratePresignedURL(bucketName, s3Key, presignedURLExpiry)
	if err != nil {
		return "", errors.Wrap(err, "failed to generate presigned URL")
	}

	return downLink, nil
}

// createTempFileWithContent creates a temporary file with the specified content and name pattern, then returns the file.
// Ensures the file is closed after writing and includes proper error handling for file operations.
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
