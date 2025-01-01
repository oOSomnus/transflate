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
						"You are an expert in text translation" +
							"I'll provide you with some text, and please translate it into Chinese." +
							"In order to provide you with some context, below are some words from the previous text.",
					),
					openai.SystemMessage(prevContext),
					openai.SystemMessage(
						"Below are the text you need to translate, please try your best to translate it." +
							"Please don't include any of your own thought." +
							"Notice that there are some random characters or symbols, it is due to the inaccuracy of ocr. Please provide translate based on your understanding of the text and sounds like native speaker." +
							"Please provide your translated text in markdown format.",
					),
					openai.UserMessage(chunk),
				},
			),
			Model:               openai.F(openai.ChatModelGPT3_5Turbo),
			MaxCompletionTokens: openai.Int(1000),
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
