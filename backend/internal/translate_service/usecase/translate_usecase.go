package usecase

import (
	"github.com/oOSomnus/transflate/internal/translate_service/domain"
	"github.com/oOSomnus/transflate/pkg/utils"
	"log"
	"runtime"
	"strings"
	"sync"
)

// maxWordsPerChunk defines the maximum number of words allowed per chunk.
// contextWordCount specifies the limit on the number of words for context.
const (
	maxWordsPerChunk = 2000 // 最大每块单词数
	contextWordCount = 50   // 上下文单词限制
)

// TranslateText splits a long string into smaller chunks, translates each chunk in parallel, and returns the full translation.
func TranslateText(longString string) (string, error) {
	chunks := utils.SplitString(longString, maxWordsPerChunk)

	// Initialize parallel processing workers
	maxNumTokens := max(runtime.NumCPU()*2, 10)
	apiTokens := make(chan struct{}, maxNumTokens)

	translatedChunks, translationErrors := initResults(len(chunks))

	translator := domain.NewGPTTranslator()
	var wg sync.WaitGroup

	for i, chunk := range chunks {
		wg.Add(1)
		go processChunk(i, chunks, chunk, translator, translatedChunks, translationErrors, apiTokens, &wg)
	}

	wg.Wait()

	reportErrors(translationErrors)

	finalTranslation := strings.Join(translatedChunks, "\n")
	return finalTranslation, nil
}

var mutex = &sync.Mutex{}

// processChunk processes a single text chunk by translating it using the provided Translator and updates the results concurrently.
// index specifies the position of the chunk in the chunks slice.
// chunks contains all text chunks to be processed.
// chunk is the specific text chunk being processed.
// translator is an instance of a domain.Translator to handle the translation task.
// translatedChunks is an array to store the translated output of each chunk.
// translationErrors is an array to store any error encountered during translation of each chunk.
// workersPool is a channel used to manage the pool of goroutines processing chunks concurrently.
// wg is a WaitGroup used to track the completion of goroutine processing for synchronization.
func processChunk(
	index int,
	chunks []string,
	chunk string,
	translator domain.Translator,
	translatedChunks []string,
	translationErrors []error,
	apiTokens chan struct{},
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	defer func() { <-apiTokens }() // 释放 worker
	apiTokens <- struct{}{}
	prevContext := getPreviousContext(index, chunks)
	result, err := translator.Translate(prevContext, chunk)
	mutex.Lock()
	translatedChunks[index] = result
	translationErrors[index] = err
	mutex.Unlock()
}

// getPreviousContext returns the last N words from the previous chunk in the list, or an empty string if index is 0.
func getPreviousContext(index int, chunks []string) string {
	if index == 0 {
		return ""
	}
	return utils.GetLastNWords(chunks[index-1], contextWordCount)
}

// initResults initializes two slices: one for strings and one for errors, each with a specified size.
func initResults(size int) ([]string, []error) {
	return make([]string, size), make([]error, size)
}

// reportErrors logs any non-nil errors in the provided slice, associating each error with its corresponding chunk index.
func reportErrors(errors []error) {
	for i, err := range errors {
		if err != nil {
			log.Printf("Error translating chunk %d: %v\n", i, err)
		}
	}
}
