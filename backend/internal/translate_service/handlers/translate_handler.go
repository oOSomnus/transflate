package handlers

import (
	"github.com/oOSomnus/transflate/internal/translate_service/domain"
	"log"
)

/*
TranslateChunk translates a given text chunk into Chinese using OpenAI's GPT API.

Parameters:
  - chunk (string): The text to be translated.

Returns:
  - (string): The translated text in Chinese if the request is successful.
  - (error): An error if there are issues with environment configuration, HTTP request creation, API response, or JSON unmarshalling.
*/
func TranslateChunk(
	prevContext string, chunk string, translator domain.Translator,
) (string, error) {
	result, err := translator.Translate(prevContext, chunk)
	if err != nil {
		log.Println("Failed to translate chunk:", err)
		return "", err
	}

	return result, nil
}
