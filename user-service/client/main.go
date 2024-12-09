package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-microservices/api/proto"
	"log"
	"log/slog"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	createResp, err := client.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "John Doe",
		Email: "john@example.com",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	slog.Info("createResp: %v", createResp.GetUser())

	getUserResp, err := client.GetUser(ctx, &proto.GetUserRequest{
		Id: createResp.GetUser().GetId(),
	})

	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	slog.Info("Retrieved user: %v", getUserResp.GetUser())

	_, err = client.GetUser(ctx, &proto.GetUserRequest{
		Id: "non-existent-id",
	})
	if err != nil {
		slog.Error("Error getting user: %v", err)
	}
}
