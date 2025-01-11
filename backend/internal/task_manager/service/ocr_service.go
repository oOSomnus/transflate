package service

import (
	"context"
	"fmt"
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"sync"
	"time"
)

// OCRClient is an interface for Optical Character Recognition operations and resource cleanup.
// ProcessOCR processes the OCR request on given file content with a specified language.
// Close releases any resources used by the OCRClient.
type OCRClient interface {
	ProcessOCR(fileContent []byte, lang string) (*pb.StringListResponse, error)
	Close() error
}

// OCRService is a struct that manages gRPC connections and operations for the OCR service.
type OCRService struct {
	clientConn     *grpc.ClientConn
	clientConnOnce sync.Once
	clientConnErr  error
}

// NewOCRService initializes and returns a new instance of OCRService.
func NewOCRService() *OCRService {
	return &OCRService{}
}

func (s *OCRService) getGRPCConn() (*grpc.ClientConn, error) {
	host := viper.GetString("ocr.host")
	if host == "" {
		return nil, fmt.Errorf("OCR service host is not configured")
	}

	s.clientConnOnce.Do(
		func() {
			s.clientConn, s.clientConnErr = grpc.NewClient(
				host+":50051",
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
		},
	)

	return s.clientConn, s.clientConnErr
}

// ProcessOCR processes the given PDF file content using OCR and specified language, returning a structured response.
func (s *OCRService) ProcessOCR(fileContent []byte, lang string) (*pb.StringListResponse, error) {
	conn, err := s.getGRPCConn()
	if err != nil {
		return nil, fmt.Errorf("failed to get gRPC connection: %w", err)
	}

	client := pb.NewOCRServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	return client.ProcessPDF(
		ctx, &pb.PDFRequest{
			PdfData:  fileContent,
			Language: lang,
		},
	)
}

// Close releases the underlying gRPC connection if it is active and returns any error encountered during closure.
func (s *OCRService) Close() error {
	if s.clientConn != nil {
		if err := s.clientConn.Close(); err != nil {
			return fmt.Errorf("failed to close gRPC connection: %w", err)
		}
	}
	return nil
}
