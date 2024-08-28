package notification17

import (
	"log"
	notificationss "user/internal/protos/notification"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Hotel() notificationss.NotificationClient {
	conn, err := grpc.NewClient("notification_service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("notification error",err)
	}
	client := notificationss.NewNotificationClient(conn)
	return client
}
