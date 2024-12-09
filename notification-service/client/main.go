package main

import (
	"context"
	"grpc-microservices/api/proto"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewNotificationServiceClient(conn)
	userId := "test-user-123"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := client.SubscribeToNotifications(ctx, &proto.SubscribeRequest{
		UserId: userId,
	})
	if err != nil {
		log.Fatalf("could not subscribe: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Горутина для получения уведомлений
	go func() {
		defer wg.Done()
		for {
			update, err := stream.Recv()
			if err != nil {
				log.Printf("Stream closed: %v", err)
				return
			}
			log.Printf("Received notification: %v", update.Notification)
		}
	}()

	// Горутина для отправки тестовых уведомлений
	go func() {
		defer wg.Done()
		// Отправляем несколько тестовых уведомлений
		for i := 1; i <= 5; i++ {
			_, err := client.SendNotification(ctx, &proto.SendNotificationRequest{
				UserId:  userId,
				Title:   "Test Notification " + time.Now().Format(time.RFC3339),
				Content: "This is test notification #" + string(rune(i+'0')),
			})
			if err != nil {
				log.Printf("Failed to send notification: %v", err)
				return
			}
			time.Sleep(2 * time.Second) // Пауза между уведомлениями
		}
		cancel() // Завершаем после отправки всех уведомлений
	}()

	wg.Wait()
}
