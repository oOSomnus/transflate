package handlers

import (
	"context"
	"errors"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"

	"log"
	"os"
	"time"
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
	prevContext string, chunk string,
) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
	ctx, cancel := context.WithTimeout(
		context.Background(), 60*time.Second,
	)
	defer cancel()
	chatCompletion, err := client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F(
				[]openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(
						"You are a professional translator. Translate the following English text into Chinese. Ignore random characters or symbols, and focus on the meaningful content. Provide the result in markdown format.",
					),
					openai.UserMessage("Previous context for reference: " + prevContext),
					openai.UserMessage("Text to translate: " + chunk),
				},
			),
			Model:               openai.F(openai.ChatModelGPT4Turbo),
			MaxCompletionTokens: openai.Int(3000),
		},
	)
	if err != nil {
		return "", err
	}

	translateResponse := chatCompletion.Choices[0].Message.Content

	if len(translateResponse) > 0 {
		return translateResponse, nil
	}
	log.Println("no response from OpenAI API")
	return "", errors.New("no response from OpenAI API")
}
