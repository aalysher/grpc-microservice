package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc-microservices/proto"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNotificationServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	userId := "test-user-123"

	notification1, err := client.SendNotification(ctx, &pb.SendNotificationRequest{
		UserId:  userId,
		Title:   "Welcome!",
		Content: "Welcome to our platform!",
	})
	if err != nil {
		log.Fatalf("could not send notification: %v", err)
	}
	log.Printf("Sent notification: %v", notification1.GetNotification())

	notification2, err := client.SendNotification(ctx, &pb.SendNotificationRequest{
		UserId:  userId,
		Title:   "New Feature",
		Content: "Check out our new feature!",
	})
	if err != nil {
		log.Fatalf("could not send notification: %v", err)
	}
	log.Printf("Sent notification: %v", notification2.GetNotification())

	notifications, err := client.GetUserNotifications(ctx, &pb.GetUserNotificationsRequest{
		UserId: userId,
	})
	if err != nil {
		log.Fatalf("could not get notifications: %v", err)
	}

	log.Printf("Retrieved %d notifications:", len(notifications.GetNotifications()))
	for i, notification := range notifications.GetNotifications() {
		log.Printf("%d. Title: %s, Content: %s, Created At: %v",
			i+1,
			notification.Title,
			notification.Content,
			time.Unix(notification.CreatedAt, 0),
		)
	}

	emptyNotifications, err := client.GetUserNotifications(ctx, &pb.GetUserNotificationsRequest{
		UserId: "non-existent-user",
	})
	if err != nil {
		log.Fatalf("could not get notifications: %v", err)
	}
	log.Printf("Notifications for non-existent user: %d", len(emptyNotifications.GetNotifications()))
}
