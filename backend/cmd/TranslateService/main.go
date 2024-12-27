package TranslateService

import (
	pb "github.com/oOSomnus/transflate/api/generated/translate"
	"github.com/oOSomnus/transflate/internal/TranslateService/server"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	listner, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTranslateServiceServer(grpcServer, &server.TranslateServiceServer{})
	if err := grpcServer.Serve(listner); err != nil {
		log.Fatal(err)
	}
}
