package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "grpc/grpc/gen"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// --- Unary RPC ---
	log.Println("Unary RPC: GetUser()")
	user, err := client.GetUser(context.Background(), &pb.GetUserRequest{UserId: 1})
	if err != nil {
		log.Fatalf("Error in GetUser: %v", err)
	}
	log.Printf("Response: %+v\n", user.User)

	// --- Server Streaming RPC ---
	log.Println("Server Streaming: GetUsersStream()")
	stream1, _ := client.GetUsersStream(context.Background(), &pb.GetUsersRequest{UserIds: []string{"101", "102", "103"}})
	for {
		res, err := stream1.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}
		log.Printf("Received user: %+v\n", res)
	}

	// --- Client Streaming RPC ---
	log.Println("Client Streaming: CreateUsersStream()")
	stream2, _ := client.CreateUsersStream(context.Background())
	for i := 1; i <= 30; i++ {
		stream2.Send(&pb.CreateUserRequest{
			Name:  "User_" + time.Now().Format("150405"),
			Email: "temp@example.com",
		})
		time.Sleep(500 * time.Millisecond)
	}
	resp, _ := stream2.CloseAndRecv()
	log.Printf("Users Created: %d | IDs: %v", resp.UsersCreated, resp.UserIds)

	// --- Bidirectional Streaming RPC ---
	log.Println("Bidirectional Streaming: Chat()")
	stream3, _ := client.Chat(context.Background())

	go func() {
		for i := 1; i <= 3; i++ {
			msg := &pb.ChatMessage{
				UserId:    "client_1",
				Message:   "Hello Server!",
				Timestamp: time.Now().Unix(),
			}
			stream3.Send(msg)
			time.Sleep(time.Second)
		}
		stream3.CloseSend()
	}()

	for {
		in, err := stream3.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error in chat stream: %v", err)
		}
		log.Printf("Server: %s", in.Message)
	}
}
