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

const ()

// TaskUsecase defines the contract for processing OCR, translating content, and returning the result as a string.
type TaskUsecase interface {
	ProcessOCRAndTranslate(username string, fileContent []byte, lang string) (string, error)
	CreateDownloadLinkWithMdString(mdString string) (string, error)
}

// TaskUsecaseImpl provides the implementation for task-related business logic utilizing a UserRepository instance.
type TaskUsecaseImpl struct {
	ur   repository.UserRepository
	tr   repository.TaskRepository
	ocrc service.OCRClient
	s3s  service.S3StorageService
	ts   service.TranslateService
}

// NewTaskUsecase creates and initializes a new TaskUsecaseImpl with the provided UserRepository.
func NewTaskUsecase(
	ur repository.UserRepository, tr repository.TaskRepository, ocrc service.OCRClient, s3s service.S3StorageService,
	ts service.TranslateService,
) *TaskUsecaseImpl {
	return &TaskUsecaseImpl{ur: ur, tr: tr, ocrc: ocrc, s3s: s3s, ts: ts}
}

// ProcessOCRAndTranslate processes a file using OCR, decreases user balance, and translates the extracted text.
// Parameters: username (string) - The name of the user requesting the process.
// fileContent ([]byte) - The content of the file to be processed via OCR.
// lang (string) - The language code used for OCR processing.
// Returns: Translated text (string) if successful, or an error when a failure occurs during processing.
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

// mergeAndCleanStrings concatenates a slice of strings and applies text cleaning to the resulting string.
func mergeAndCleanStrings(lines []string) string {
	var builder strings.Builder
	for _, line := range lines {
		builder.WriteString(line)
	}
	return utils.TextCleaning(builder.String())
}

const (
	s3KeyPrefix        = "mds/"
	tempFilePattern    = "respMd-*.md"
	presignedURLExpiry = time.Hour
)

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
