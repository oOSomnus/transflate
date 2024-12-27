package server

import (
	"context"
	"fmt"
	pb "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/oOSomnus/transflate/internal/TranslateService/handlers"
	"github.com/oOSomnus/transflate/pkg/utils"
	"strings"
	"sync"
)

type TranslateServiceServer struct {
	pb.UnimplementedTranslateServiceServer
}

/*
ProcessTranslation handles the translation of a long string into chunks, processes them in parallel, and combines the results.

Parameters:
  - ctx (context.Context): The context for the request, used for managing deadlines and cancellations.
  - req (*pb.TranslateRequest): The translation request containing the input data to process.

Returns:
  - (*pb.TranslateResult): The result containing the fully translated string after processing.
  - (error): An error if the translation process fails or encounters issues.
*/
func (s *TranslateServiceServer) ProcessTranslation(ctx context.Context, req *pb.TranslateRequest) (*pb.TranslateResult, error) {
	longString := "Your long string here. It can be very lengthy."
	maxTokens := 1000
	// partition string
	chunks := utils.SplitString(longString, maxTokens)

	// translate in parallel
	var wg sync.WaitGroup
	results := make([]string, len(chunks))
	errors := make([]error, len(chunks))

	for i, chunk := range chunks {
		wg.Add(1)
		go func(i int, chunk string) {
			defer wg.Done()
			result, err := handlers.TranslateChunk(chunk)
			results[i] = result
			errors[i] = err
		}(i, chunk)
	}

	wg.Wait()

	// error checking
	for i, err := range errors {
		if err != nil {
			fmt.Printf("Error translating chunk %d: %v\n", i, err)
		}
	}

	// combine results
	finalTranslation := strings.Join(results, " ")
	return &pb.TranslateResult{Lines: finalTranslation}, nil
}
