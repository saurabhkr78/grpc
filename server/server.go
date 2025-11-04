package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "grpc/grpc/gen"
)

// Implementing the UserServiceServer interface
type UserServer struct {
	pb.UnimplementedUserServiceServer
}

// Unary RPC
func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Received GetUser request for user_id: %d", req.GetUserId())

	user := &pb.UserResponse{
		Id:       fmt.Sprintf("%d", req.GetUserId()),
		Name:     fmt.Sprintf("User_%d", req.GetUserId()),
		Email:    fmt.Sprintf("user_%d@example.com", req.GetUserId()),
		IsActive: true,
	}

	return &pb.GetUserResponse{User: user}, nil
}

// Server Streaming RPC
func (s *UserServer) GetUsersStream(req *pb.GetUsersRequest, stream pb.UserService_GetUsersStreamServer) error {
	for _, id := range req.UserIds {
		res := &pb.UserResponse{
			Id:       id,
			Name:     fmt.Sprintf("User_%s", id),
			Email:    fmt.Sprintf("user%s@example.com", id),
			IsActive: true,
		}

		if err := stream.Send(res); err != nil {
			return err
		}

		log.Printf("Sent user: %s", id)
		time.Sleep(time.Second)
	}
	return nil
}

// Client Streaming RPC
func (s *UserServer) CreateUsersStream(stream pb.UserService_CreateUsersStreamServer) error {
	var count int32
	var userIDs []string

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.CreateUsersResponse{
				UsersCreated: count,
				UserIds:      userIDs,
			})
		}
		if err != nil {
			return err
		}

		count++
		newID := fmt.Sprintf("%d", time.Now().UnixNano())
		userIDs = append(userIDs, newID)
		log.Printf("Created user: %s (%s)", req.Name, req.Email)
	}
}

// Bidirectional Streaming RPC
func (s *UserServer) Chat(stream pb.UserService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("[%s]: %s", msg.UserId, msg.Message)

		resp := &pb.ChatMessage{
			UserId:    "Server",
			Message:   "Received: " + msg.Message,
			Timestamp: time.Now().Unix(),
		}

		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

func main() {
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServer{})

	log.Printf("gRPC Server running on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
