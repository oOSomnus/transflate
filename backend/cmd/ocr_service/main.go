package main

import (
	pb "github.com/oOSomnus/transflate/api/generated/ocr"
	"github.com/oOSomnus/transflate/internal/ocr_service/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[OCR Service] ")
}

func main() {
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
