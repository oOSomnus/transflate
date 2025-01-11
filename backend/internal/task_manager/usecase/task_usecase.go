package usecase

import (
	"errors"
	"github.com/oOSomnus/transflate/internal/task_manager/repository"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"strings"
)

// TaskUsecase defines the contract for processing OCR, translating content, and returning the result as a string.
type TaskUsecase interface {
	ProcessOCRAndTranslate(username string, fileContent []byte, lang string) (string, error)
}

// TaskUsecaseImpl provides the implementation for task-related business logic utilizing a UserRepository instance.
type TaskUsecaseImpl struct {
	ur repository.UserRepository
}

// NewTaskUsecase creates and initializes a new TaskUsecaseImpl with the provided UserRepository.
func NewTaskUsecase(ur repository.UserRepository) *TaskUsecaseImpl {
	return &TaskUsecaseImpl{ur: ur}
}

// ProcessOCRAndTranslate processes a file using OCR, decreases user balance, and translates the extracted text.
// Parameters: username (string) - The name of the user requesting the process.
// fileContent ([]byte) - The content of the file to be processed via OCR.
// lang (string) - The language code used for OCR processing.
// Returns: Translated text (string) if successful, or an error when a failure occurs during processing.
func (t *TaskUsecaseImpl) ProcessOCRAndTranslate(username string, fileContent []byte, lang string) (string, error) {
	ocrResponse, err := service.ProcessOCR(fileContent, lang)
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
	translatedResponse, err := service.TranslateText(cleanedText)
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
