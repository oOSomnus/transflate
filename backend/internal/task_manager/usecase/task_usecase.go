package usecase

import (
	"errors"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"strings"
)

/*
ProcessOCRAndTranslate performs OCR on the given file content, translates the extracted text, and manages user balance.

Parameters:
  - usernameStr (string): The username of the user for whom the OCR and translation is being performed.
  - fileContent ([]byte): The content of the file to be processed through OCR.
  - lang (string): The language code for OCR processing.

Returns:
  - (string): The translated text from the OCR-processed file.
  - (error): An error if the OCR processing, text translation, or balance decrement fails.
*/
func ProcessOCRAndTranslate(usernameStr string, fileContent []byte, lang string) (string, error) {
	ocrResponse, err := service.ProcessOCR(fileContent, lang)
	if err != nil {
		log.Println("ocr error")
		return "", err
	}
	if ocrResponse == nil {
		log.Println("ocrResponse is nil")
		return "", errors.New("ocrResponse is nil")
	}
	respLines := ocrResponse.Lines
	log.Println("File processed successfully.")
	//merge strings
	var builder strings.Builder
	for _, line := range respLines {
		builder.WriteString(line)
	}
	mergedString := builder.String()
	mergedString = utils.TextCleaning(mergedString)
	//decrease balance
	numPages := int(ocrResponse.PageNum)
	err = DecreaseBalance(usernameStr, numPages)
	if err != nil {
		log.Println("Error decreasing balance.")
		return "", err
	}
	log.Printf("Username: %s", usernameStr)
	//translate
	transResponse, err := service.TranslateText(mergedString)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return transResponse.Lines, nil
}
