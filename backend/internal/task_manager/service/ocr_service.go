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

var (
	grpcConn     *grpc.ClientConn
	grpcConnOnce sync.Once
	grpcConnErr  error
)

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

func CloseOcrGRPCConn() error {
	if grpcConn != nil {
		return grpcConn.Close()
	}
	return nil
}

// ProcessOCR processes an OCR request by sending the given file content and language to the OCR service via gRPC.
// Parameters: fileContent ([]byte) - The content of the file to be processed.
// lang (string) - The language code for OCR processing.
// Returns: *pb.StringListResponse containing extracted text data, or an error if the process fails.
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
