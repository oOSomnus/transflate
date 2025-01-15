package domain

import (
	"context"
	"errors"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/spf13/viper"
	"time"
)

// GPTTranslator is a struct that provides translation capabilities using OpenAI's API.
// It wraps a client for interacting with OpenAI's services.
type GPTTranslator struct {
	client *openai.Client
}

// NewGPTTranslator initializes and returns a new instance of GPTTranslator with an OpenAI client using the API key from config.
func NewGPTTranslator() *GPTTranslator {
	//utils.LoadEnv()
	apiKey := viper.GetString("openai.api.key")
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &GPTTranslator{
		client: client,
	}

}

// Translate uses GPT-4 Turbo to translate an English text input into Chinese, excluding prior context and irrelevant symbols.
// prevContext provides optional reference data, while text represents the content to translate.
// Returns the translated text in markdown format or an error if the translation request fails.
func (g *GPTTranslator) Translate(prevContext, text string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	chatCompletion, err := g.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F(
				[]openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(
						"You are a professional translator. Translate the following text into Chinese. Ignore random characters or symbols, and focus on the meaningful content. Provide the result in markdown format and try your best separating paragraphs. You don't need to translate the previous context.",
					),
					openai.UserMessage("Previous context for reference: " + prevContext),
					openai.UserMessage("Text to translate: " + text),
				},
			),
			Model: openai.F(openai.ChatModelGPT4Turbo),
			//MaxCompletionTokens: openai.Int(3000),
		},
	)
	if err != nil {
		return "", err
	}

	translateResponse := chatCompletion.Choices[0].Message.Content
	if len(translateResponse) > 0 {
		return translateResponse, nil
	}
	return "", errors.New("no response from OpenAI API")
}
