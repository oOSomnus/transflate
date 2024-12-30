package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/oOSomnus/transflate/internal/TranslateService/domain"
	"io"
	"log"
	"net/http"
	"os"
)

/*
TranslateChunk translates a given text chunk into Chinese using OpenAI's GPT API.

Parameters:
  - chunk (string): The text to be translated.

Returns:
  - (string): The translated text in Chinese if the request is successful.
  - (error): An error if there are issues with environment configuration, HTTP request creation, API response, or JSON unmarshalling.
*/
func TranslateChunk(prevContext string, chunk string) (string, error) {
	apiURL := "https://api.openai.com/v1/chat/completions"
	//maxTokens := 1000
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	requestBody := domain.OpenAIRequest{
		Model: "gpt-4",
		Messages: []domain.OpenAIMessage{
			{Role: "system", Content: "You are an expert in text translation.\n"},
			{Role: "system", Content: "I'll provide you with some text, and please translate it into Chinese.\n"},
			{Role: "system", Content: "In order to provide you with some context, below are some words from the previous text.\n"},
			{Role: "system", Content: prevContext},
			{Role: "system", Content: "Below are the text you need to translate, please try your best to translate it.\n"},
			{Role: "system", Content: "Please only provide the translated context. Please don't include any of your own words.\n"},
			{Role: "user", Content: chunk},
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var openAIResp domain.OpenAIResponse
	if err := json.Unmarshal(respBody, &openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) > 0 {
		return openAIResp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from OpenAI API")
}
