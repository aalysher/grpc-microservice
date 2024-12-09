package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"grpc-microservices/api/proto"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {
	proto.UnimplementedNotificationServiceServer
	notifications map[string][]*proto.Notification
	subscribers   map[string][]chan *proto.NotificationUpdate
	mu            sync.RWMutex
}

func NewServer() *server {
	return &server{
		notifications: make(map[string][]*proto.Notification),
		subscribers:   make(map[string][]chan *proto.NotificationUpdate),
	}
}

func (s *server) SendNotification(ctx context.Context, req *proto.SendNotificationRequest) (*proto.SendNotificationResponse, error) {
	notification := &proto.Notification{
		Id:        uuid.New().String(),
		UserId:    req.UserId,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: time.Now().Unix(),
	}

	s.mu.Lock()
	s.notifications[req.UserId] = append(s.notifications[req.UserId], notification)

	if subscribers, ok := s.subscribers[req.UserId]; ok {
		update := &proto.NotificationUpdate{
			Notification: notification,
			Timestamp:    time.Now().Unix(),
		}
		for _, subscriber := range subscribers {
			select {
			case subscriber <- update:
			default:

			}
		}
	}
	s.mu.Unlock()

	return &proto.SendNotificationResponse{
		Notification: notification,
	}, nil
}

func (s *server) GetUserNotifications(ctx context.Context, req *proto.GetUserNotificationsRequest) (*proto.GetUserNotificationsResponse, error) {
	s.mu.RLock()
	notifications, exists := s.notifications[req.UserId]
	s.mu.RUnlock()

	if !exists {
		log.Printf("No notifications found for user %s", req.UserId)
		return &proto.GetUserNotificationsResponse{
			Notifications: []*proto.Notification{},
		}, nil
	}

	log.Printf("Retrieved %d notifications for user %s", len(notifications), req.UserId)
	return &proto.GetUserNotificationsResponse{
		Notifications: notifications,
	}, nil
}

func (s *server) SubscribeToNotifications(req *proto.SubscribeRequest, stream proto.NotificationService_SubscribeToNotificationsServer) error {
	updateChan := make(chan *proto.NotificationUpdate, 10)

	s.mu.Lock()
	s.subscribers[req.UserId] = append(s.subscribers[req.UserId], updateChan)
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		subscribers := s.subscribers[req.UserId]
		for i, ch := range subscribers {
			if ch == updateChan {
				s.subscribers[req.UserId] = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
		s.mu.Unlock()
		close(updateChan)
	}()

	s.mu.Lock()
	existingNotifications := s.notifications[req.UserId]
	s.mu.Unlock()

	for _, notification := range existingNotifications {
		if err := stream.Send(&proto.NotificationUpdate{
			Notification: notification,
			Timestamp:    time.Now().Unix(),
		}); err != nil {
			return fmt.Errorf("error sending notification: %v", err)
		}
	}

	for {
		select {
		case update := <-updateChan:
			if err := stream.Send(update); err != nil {
				return fmt.Errorf("error sending notification: %v", err)
			}
		case <-stream.Context().Done():
			return nil
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50052") // Используем другой порт
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	proto.RegisterNotificationServiceServer(s, NewServer())

	log.Printf("notification server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
