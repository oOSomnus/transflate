package service

import (
	"context"
	pbt "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type TranslateService interface {
	TranslateText(text string) (*pbt.TranslateResult, error)
	CloseTransGrpcConn() error
}

type TranslateServiceImpl struct {
	translateClient pbt.TranslateServiceClient
	transGrpcConn   *grpc.ClientConn
}

func NewTranslateService() (*TranslateServiceImpl, error) {
	service := &TranslateServiceImpl{}
	err := service.getTransGrpcConn()
	if err != nil {
		log.Printf("Error initializing translation gRPC connection: %v", err)
		return nil, err
	}
	client := pbt.NewTranslateServiceClient(service.transGrpcConn)
	service.translateClient = client
	return service, nil

}

// getTransGrpcConn initializes and returns a gRPC client connection for the translation service.
// It ensures that the connection setup is executed only once using sync.Once.
// Returns the client connection and any error encountered during initialization.
func (t *TranslateServiceImpl) getTransGrpcConn() error {
	grpcServiceHost := viper.GetString("translate.host")
	conn, transErr := grpc.NewClient(
		grpcServiceHost+":50052", grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	t.transGrpcConn = conn
	if transErr != nil {
		log.Printf("TransGrpcErr: %v", transErr)
		return transErr
	}
	return nil
}

// CloseTransGrpcConn closes the translation gRPC connection and returns an error if the closure fails.
func (t *TranslateServiceImpl) CloseTransGrpcConn() error {
	if t.transGrpcConn == nil {
		log.Println("Translation gRPC connection is already closed")
		return nil
	}
	err := t.transGrpcConn.Close()
	if err != nil {
		log.Printf(
			"Error closing translation gRPC connection: %v",
			err,
		)
		return err
	}
	return nil
}

// TranslateText translates the given text by sending it to the translation service via gRPC and returns the result.
func (t *TranslateServiceImpl) TranslateText(text string) (*pbt.TranslateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	response, err := t.translateClient.ProcessTranslation(ctx, &pbt.TranslateRequest{Text: text})
	if err != nil {
		log.Printf("Error translating text: %v", err)
		return nil, err
	}

	return response, nil
}
