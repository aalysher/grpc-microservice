package main

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	pb "grpc-microservices/proto"
	"log"
	"net"
	"time"
)

type server struct {
	pb.UnimplementedNotificationServiceServer
	notifications map[string][]*pb.Notification
}

func NewServer() *server {
	return &server{
		notifications: make(map[string][]*pb.Notification),
	}
}

func (s *server) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	notification := &pb.Notification{
		Id:        uuid.New().String(),
		UserId:    req.UserId,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now().Unix(),
	}

	s.notifications[req.UserId] = append(s.notifications[req.UserId], notification)
	log.Printf("Sent notification to user %s: %v", req.UserId, notification)

	return &pb.SendNotificationResponse{
		Notification: notification,
	}, nil
}

func (s *server) GetUserNotifications(ctx context.Context, req *pb.GetUserNotificationsRequest) (*pb.GetUserNotificationsResponse, error) {
	notifications, exists := s.notifications[req.UserId]
	if !exists {
		return &pb.GetUserNotificationsResponse{
			Notifications: []*pb.Notification{},
		}, nil
	}

	log.Printf("Retrieved %d notifications for user %s", len(notifications), req.UserId)

	return &pb.GetUserNotificationsResponse{
		Notifications: notifications,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50052") // Используем другой порт
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, NewServer())

	log.Printf("notification server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
