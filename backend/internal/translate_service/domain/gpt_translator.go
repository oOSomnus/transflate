package domain

import (
	"context"
	"errors"
	"github.com/oOSomnus/transflate/pkg/utils"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"time"
)

type GPTTranslator struct {
	client *openai.Client
}

/*
NewGPTTranslator creates a new instance of GPTTranslator.

Parameters:
  - apiKey (string): The API key used to authenticate with the OpenAI API.

Returns:
  - (*GPTTranslator): A pointer to the initialized GPTTranslator instance.
*/
func NewGPTTranslator() *GPTTranslator {
	utils.LoadEnv()
	apiKey := utils.GetEnv("OPENAI_API_KEY")
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &GPTTranslator{
		client: client,
	}

}

/*
Translate performs a translation of the given text into Chinese using the OpenAI API.

Parameters:
  - prevContext (string): Additional contextual information to assist the translation.
  - text (string): The English text to be translated.

Returns:
  - (string): The translated text in markdown format.
  - (error): An error if the translation process fails or no response is received from the OpenAI API.
*/
func (g *GPTTranslator) Translate(prevContext, text string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	chatCompletion, err := g.client.Chat.Completions.New(
		ctx,
		openai.ChatCompletionNewParams{
			Messages: openai.F(
				[]openai.ChatCompletionMessageParamUnion{
					openai.SystemMessage(
						"You are a professional translator. Translate the following English text into Chinese. Ignore random characters or symbols, and focus on the meaningful content. Provide the result in markdown format. You don't need to translate the previous context.",
					),
					openai.UserMessage("Previous context for reference: " + prevContext),
					openai.UserMessage("Text to translate: " + text),
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
	return "", errors.New("no response from OpenAI API")
}
