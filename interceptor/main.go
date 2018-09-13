package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/asa-taka/grpc-tiny-examples/interceptor/proto"
)

const (
	port = 10000
)

// gRPC Service Implementation
// ---------------------------

type greetingServer struct {
	name string
	mu   sync.Mutex
}

func newServer() *greetingServer {
	return &greetingServer{name: "gentle-server"}
}

func (s *greetingServer) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	res := &pb.HelloResponse{
		Message: fmt.Sprintf("Hello %s, I am %s", in.Name, s.name),
	}
	return res, nil
}

// Interceptor
// -----------

// A simple implementation example of `grpc.UnaryServerInterceptor`
// See the signature of https://godoc.org/google.golang.org/grpc#UnaryServerInterceptor
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	res, err := handler(ctx, req)
	log.Println("My interceptor called!")
	log.Printf("%s: %v -> %v", info.FullMethod, req, res)
	return res, err
}

// Main Program
// ------------

func main() {

	// Set interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	// Register service implementations
	pb.RegisterGreetingServer(grpcServer, newServer())

	reflection.Register(grpcServer) // for grpc_cli

	// ...then, followings are basic gRPC server setup...

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server starts on localhost:%d", port)
	grpcServer.Serve(lis)
}
