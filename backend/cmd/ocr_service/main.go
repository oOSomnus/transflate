package main

import (
	"fmt"
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"github.com/oOSomnus/transflate/internal/ocr_service/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[OCR Service] ")
}

func main() {
	// viper config
	env := os.Getenv("TRANSFLATE_ENV")
	if env == "" {
		env = "local"
	}
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	// Start gRPC service
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOCRServiceServer(grpcServer, &server.OCRServiceServer{})

	log.Println("Starting gRPC server on :50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
