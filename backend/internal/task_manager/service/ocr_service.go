package service

import (
	"context"
	"time"

	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProcessOCR(fileContent []byte, lang string) (*pb.StringListResponse, error) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

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
