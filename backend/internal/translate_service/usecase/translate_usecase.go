package usecase

import (
	"github.com/oOSomnus/transflate/internal/translate_service/domain"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"runtime"
	"strings"
	"sync"
)

const (
	maxWordsPerChunk = 2000 // 最大每块单词数
	contextWordCount = 50   // 上下文单词限制
)

// TranslateText translates a long string into another language.
func TranslateText(longString string) (string, error) {
	chunks := utils.SplitString(longString, maxWordsPerChunk)

	// Initialize parallel processing workers
	numWorkers := runtime.NumCPU() * 2
	workersPool := make(chan struct{}, numWorkers)

	translatedChunks, translationErrors := initResults(len(chunks))

	translator := domain.NewGPTTranslator()
	var wg sync.WaitGroup

	for i, chunk := range chunks {
		wg.Add(1)
		workersPool <- struct{}{}
		go processChunk(i, chunks, chunk, translator, translatedChunks, translationErrors, workersPool, &wg)
	}

	wg.Wait()

	reportErrors(translationErrors)

	finalTranslation := strings.Join(translatedChunks, "\n")
	return finalTranslation, nil
}

// processChunk handles the translation of an individual chunk.
func processChunk(
	index int,
	chunks []string,
	chunk string,
	translator domain.Translator,
	translatedChunks []string,
	translationErrors []error,
	workersPool chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	defer func() { <-workersPool }() // 释放 worker

	prevContext := getPreviousContext(index, chunks)
	result, err := translator.Translate(prevContext, chunk)

	translatedChunks[index] = result
	translationErrors[index] = err
}

// getPreviousContext retrieves the previous chunk's context for translation.
func getPreviousContext(index int, chunks []string) string {
	if index == 0 {
		return ""
	}
	return utils.GetLastNWords(chunks[index-1], contextWordCount)
}

// initResults initializes the arrays for storing translation results and errors.
func initResults(size int) ([]string, []error) {
	return make([]string, size), make([]error, size)
}

// reportErrors logs all translation errors.
func reportErrors(errors []error) {
	for i, err := range errors {
		if err != nil {
			log.Printf("Error translating chunk %d: %v\n", i, err)
		}
	}
}
