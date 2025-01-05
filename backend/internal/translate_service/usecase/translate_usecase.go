package usecase

import (
	"github.com/oOSomnus/transflate/internal/translate_service/domain"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"runtime"
	"strings"
	"sync"
)

/*
TranslateText translates a long string into another language by partitioning the input and processing the chunks in parallel.

Parameters:
  - longString (string): The input text to be translated.

Returns:
  - (string): The translated text, reconstructed from the processed chunks.
  - (error): An error if any issues occur during translation.
*/
func TranslateText(longString string) (string, error) {
	maxWords := 2000
	// partition string
	chunks := utils.SplitString(longString, maxWords)

	// translate in parallel
	var wg sync.WaitGroup
	results := make([]string, len(chunks))
	errors := make([]error, len(chunks))
	// Acuiring current CPU nums
	numCPU := runtime.NumCPU()
	workersPool := make(chan struct{}, numCPU*2)
	translator := domain.NewGPTTranslator()
	for i, chunk := range chunks {
		wg.Add(1)
		workersPool <- struct{}{}
		go func(i int, chunk string) {
			defer wg.Done()
			defer func() { <-workersPool }()
			prevContext := ""
			if i != 0 {
				prevContext = utils.GetLastNWords(chunks[i-1], 50)
			}
			result, err := translateChunk(prevContext, chunk, translator)
			results[i] = result
			errors[i] = err
		}(i, chunk)
	}

	wg.Wait()

	// error checking
	for i, err := range errors {
		if err != nil {
			log.Printf("Error translating chunk %d: %v\n", i, err)
		}
	}

	// combine results
	finalTranslation := strings.Join(results, "\n")
	return finalTranslation, nil
}

/*
TranslateChunk translates a given text chunk into Chinese using OpenAI's GPT API.

Parameters:
  - chunk (string): The text to be translated.

Returns:
  - (string): The translated text in Chinese if the request is successful.
  - (error): An error if there are issues with environment configuration, HTTP request creation, API response, or JSON unmarshalling.
*/
func translateChunk(
	prevContext string, chunk string, translator domain.Translator,
) (string, error) {
	result, err := translator.Translate(prevContext, chunk)
	if err != nil {
		log.Println("Failed to translate chunk:", err)
		return "", err
	}
	return result, nil
}
