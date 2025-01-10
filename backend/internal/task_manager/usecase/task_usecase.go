package usecase

import (
	"errors"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"strings"
)

// ProcessOCRAndTranslate processes the OCR on a file, translates the extracted text, and handles user balance decrement.
// Parameters: usernameStr (string) - The username for balance deduction; fileContent ([]byte) - File data for OCR;
// lang (string) - The target OCR language.
// Returns: Translated text (string) or an error if the process fails.
// ProcessOCRAndTranslate processes the OCR on a file, translates the extracted text, and handles user balance decrement.
func ProcessOCRAndTranslate(username string, fileContent []byte, lang string) (string, error) {
	ocrResponse, err := service.ProcessOCR(fileContent, lang)
	if err != nil || ocrResponse == nil {
		log.Println("Error during OCR processing:", err)
		return "", errors.New("failed to process OCR")
	}

	// Merge and clean OCR response lines
	cleanedText := mergeAndCleanStrings(ocrResponse.Lines)

	// Decrease user balance based on the number of pages
	numPages := int(ocrResponse.PageNum)
	if err = DecreaseBalance(username, numPages); err != nil {
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

// mergeAndCleanStrings merges multiple strings into one and applies text cleaning.
func mergeAndCleanStrings(lines []string) string {
	var builder strings.Builder
	for _, line := range lines {
		builder.WriteString(line)
	}
	return utils.TextCleaning(builder.String())
}
