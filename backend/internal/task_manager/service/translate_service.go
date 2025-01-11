package service

import (
	"context"
	pbt "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
	"time"
)

// transGrpcConn represents a gRPC client connection, shared across the application.
// transGrpcOnce ensures the initialization of transGrpcConn is done only once.
// transGrpcErr stores any error encountered during the initialization of transGrpcConn.
var (
	transGrpcConn *grpc.ClientConn
	transGrpcOnce sync.Once
	transGrpcErr  error
)

// getTransGrpcConn initializes and returns a gRPC client connection for the translation service.
// It ensures that the connection setup is executed only once using sync.Once.
// Returns the client connection and any error encountered during initialization.
func getTransGrpcConn() (*grpc.ClientConn, error) {
	transGrpcOnce.Do(
		func() {
			//utils.LoadEnv()
			grpcServiceHost := viper.GetString("translate.host")
			transGrpcConn, transGrpcErr = grpc.NewClient(
				grpcServiceHost+":50052", grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if transGrpcErr != nil {
				log.Printf("TransGrpcErr: %v", transGrpcErr)
			}
		},
	)
	return transGrpcConn, transGrpcErr
}

// CloseTransGrpcConn closes the translation gRPC connection and returns an error if the closure fails.
func CloseTransGrpcConn() error {
	err := transGrpcConn.Close()
	if err != nil {
		return err
	}
	return nil
}

// TranslateText translates the given text by sending it to the translation service via gRPC and returns the result.
func TranslateText(text string) (*pbt.TranslateResult, error) {
	conn, err := getTransGrpcConn()
	if err != nil {
		return nil, err
	}

	client := pbt.NewTranslateServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	response, err := client.ProcessTranslation(ctx, &pbt.TranslateRequest{Text: text})
	if err != nil {
		return nil, err
	}

	return response, nil
}
