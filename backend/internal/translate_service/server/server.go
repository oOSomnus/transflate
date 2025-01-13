package server

import (
	"context"
	pb "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/oOSomnus/transflate/internal/translate_service/usecase"
	"log"
)

// TranslateServiceServer implements the server API for the TranslateService service.
// It embeds UnimplementedTranslateServiceServer for forward compatibility.
type TranslateServiceServer struct {
	pb.UnimplementedTranslateServiceServer
}

// ProcessTranslation handles incoming translation requests and returns the translated result or an error if translation fails.
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
