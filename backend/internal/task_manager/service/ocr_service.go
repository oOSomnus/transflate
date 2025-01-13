package service

import (
	"context"
	"fmt"
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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
	clientConn *grpc.ClientConn
	grpcClient pb.OCRServiceClient
}

// NewOCRService initializes and returns a new instance of OCRService.
func NewOCRService() (*OCRService, error) {
	service := &OCRService{}
	err := service.getGRPCConn()
	if err != nil {
		log.Println("Error initializing OCR gRPC connection")
		return nil, err
	}
	client := pb.NewOCRServiceClient(service.clientConn)
	service.grpcClient = client
	return service, nil
}

func (s *OCRService) getGRPCConn() error {
	host := viper.GetString("ocr.host")
	if host == "" {
		return fmt.Errorf("OCR service host is not configured")
	}
	clientConn, clientConnErr := grpc.NewClient(
		host+":50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if clientConnErr != nil {
		log.Println("Error initializing OCR gRPC connection")
		return clientConnErr
	}
	s.clientConn = clientConn
	return nil
}

// ProcessOCR processes the given PDF file content using OCR and specified language, returning a structured response.
func (s *OCRService) ProcessOCR(fileContent []byte, lang string) (*pb.StringListResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	client := s.grpcClient
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
