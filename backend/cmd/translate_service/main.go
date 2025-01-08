package main

import (
	"fmt"
	pb "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/oOSomnus/transflate/internal/translate_service/server"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetPrefix("[Translate Service] ")
}

func main() {
	// viper config
	env := os.Getenv("TRANSFLATE_ENV")
	if env == "" {
		env = "local"
	}
	viper.SetConfigName(fmt.Sprintf("config.%s", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Failed to read config file: %s", err)
	}
	port := ":50052"
	log.Printf("Starting server on port %s", port)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTranslateServiceServer(grpcServer, &server.TranslateServiceServer{})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
