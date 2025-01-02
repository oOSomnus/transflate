package usecase

import (
	"errors"
	"github.com/oOSomnus/transflate/internal/task_manager/service"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"strings"
)

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
		log.Println("Error processing ocr response.")
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
