package server

import (
	"context"
	pb "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/oOSomnus/transflate/internal/translate_service/usecase"
	"log"
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
func (s *TranslateServiceServer) ProcessTranslation(ctx context.Context, req *pb.TranslateRequest) (
	*pb.TranslateResult, error,
) {
	longString := req.Text
	finalTranslation, err := usecase.TranslateText(longString)
	if err != nil {
		log.Println("translation error", err)
		return nil, err
	}
	return &pb.TranslateResult{Lines: finalTranslation}, nil
}
