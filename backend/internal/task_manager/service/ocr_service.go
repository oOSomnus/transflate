package service

import (
	"context"
	"github.com/spf13/viper"
	"sync"
	"time"

	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpcConn holds the gRPC client connection instance, which is initialized only once.
// grpcConnOnce ensures that the gRPC client connection is established exactly once.
// grpcConnErr captures any error that occurs during the initialization of grpcConn.
var (
	grpcConn     *grpc.ClientConn
	grpcConnOnce sync.Once
	grpcConnErr  error
)

// getOcrGRPCConn establishes and returns a gRPC client connection to the OCR service, ensuring it's created only once.
func getOcrGRPCConn() (*grpc.ClientConn, error) {
	grpcConnOnce.Do(
		func() {
			//utils.LoadEnv()
			grpcServiceHost := viper.GetString("ocr.host")
			grpcConn, grpcConnErr = grpc.NewClient(
				grpcServiceHost+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
		},
	)
	return grpcConn, grpcConnErr
}

// CloseOcrGRPCConn closes the OCR gRPC connection if it exists and returns an error if the closure fails.
func CloseOcrGRPCConn() error {
	if grpcConn != nil {
		return grpcConn.Close()
	}
	return nil
}

// ProcessOCR processes a PDF file's content using OCR with the specified language and returns the extracted text and metadata.
// Parameters: fileContent ([]byte) - The binary content of the PDF file.
// lang (string) - The language code for OCR processing.
// Returns: Extracted text as a StringListResponse and an error if processing fails.
func ProcessOCR(fileContent []byte, lang string) (*pb.StringListResponse, error) {
	conn, err := getOcrGRPCConn()
	if err != nil {
		return nil, err
	}

	client := pb.NewOCRServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return client.ProcessPDF(
		ctx, &pb.PDFRequest{
			PdfData:  fileContent,
			Language: lang,
		},
	)
}
