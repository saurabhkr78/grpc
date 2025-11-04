package main

import (
	"context"
	"log"
	"net"

	pb "grpc/proto"

	"google.golang.org/grpc"
)

// server is used to implement hello_grpc.HelloServiceServer.
type server struct {
	pb.UnimplementedHelloServiceServer
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterHelloServiceServer(grpcServer, &server{})
	log.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Println("Received  request from client")
	log.Printf("Hello %s! Welcome to gRPC server.\n", req.Name)
	return &pb.HelloResponse{Message: "Hello " + req.Name + "! Welcome to gRPC services."}, nil

}
