package service

import (
	"context"
	"time"

	pbt "github.com/oOSomnus/transflate/api/generated/translate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TranslateText(text string) (*pbt.TranslateResult, error) {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbt.NewTranslateServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	response, err := client.ProcessTranslation(ctx, &pbt.TranslateRequest{Text: text})
	if err != nil {
		return nil, err
	}

	return response, nil
}
