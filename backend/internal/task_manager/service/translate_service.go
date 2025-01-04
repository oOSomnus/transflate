package service

import (
	"context"
	"log"
	"sync"
	"time"

	pbt "github.com/oOSomnus/transflate/api/generated/translate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	transGrpcConn *grpc.ClientConn
	transGrpcOnce sync.Once
	transGrpcErr  error
)

func getTransGrpcConn() (*grpc.ClientConn, error) {
	transGrpcOnce.Do(
		func() {
			transGrpcConn, transGrpcErr = grpc.NewClient(
				"localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if transGrpcErr != nil {
				log.Printf("TransGrpcErr: %v", transGrpcErr)
			}
		},
	)
	return transGrpcConn, transGrpcErr
}

func CloseTransGrpcConn() error {
	err := transGrpcConn.Close()
	if err != nil {
		return err
	}
	return nil
}

func TranslateText(text string) (*pbt.TranslateResult, error) {
	conn, err := getTransGrpcConn()
	if err != nil {
		return nil, err
	}

	client := pbt.NewTranslateServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := client.ProcessTranslation(ctx, &pbt.TranslateRequest{Text: text})
	if err != nil {
		return nil, err
	}

	return response, nil
}
