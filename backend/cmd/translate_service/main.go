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
	viper.AddConfigPath("..")
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTranslateServiceServer(grpcServer, &server.TranslateServiceServer{})
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
