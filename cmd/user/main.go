package main

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"grpc-microservices/api/proto"
	"log"
	"net"
)

type server struct {
	proto.UnimplementedUserServiceServer
	users              map[string]*proto.User
	notificationClient proto.NotificationServiceClient
}

func NewServer() *server {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to notification service: %v", err)
	}

	return &server{
		users:              make(map[string]*proto.User),
		notificationClient: proto.NewNotificationServiceClient(conn),
	}
}

func (s *server) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	userID := uuid.New().String()

	user := &proto.User{
		Id:    userID,
		Name:  req.Name,
		Email: req.Email,
	}

	s.users[userID] = user

	log.Printf("Created user: %v", user)

	_, err := s.notificationClient.SendNotification(ctx, &proto.SendNotificationRequest{
		UserId:  userID,
		Title:   "Welcome to the platform!",
		Content: "Hello " + user.Name + "! Your account has been successfully created.",
	})
	if err != nil {
		log.Printf("Warning: failed to send welcome notification: %v", err)
	}

	return &proto.CreateUserResponse{
		User: user,
	}, nil
}

func (s *server) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	user, exists := s.users[req.Id]
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", req.Id)
	}

	log.Printf("Retrieved user: %v", user)

	return &proto.GetUserResponse{
		User: user,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, NewServer())

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
