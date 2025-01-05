package service

import (
	"context"
	"github.com/oOSomnus/transflate/pkg/utils"
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
			utils.LoadEnv()
			grpcServiceName := utils.GetEnv("OCR_CONTAINER_NAME")
			grpcConn, grpcConnErr = grpc.NewClient(
				grpcServiceName+":50051", grpc.WithTransportCredentials(insecure.NewCredentials()),
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
