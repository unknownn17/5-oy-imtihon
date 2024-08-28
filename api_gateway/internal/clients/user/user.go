package userservice

import (
	"api/internal/protos/user"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UserClinet() user.UserClient {
	conn, err := grpc.NewClient("user_service:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := user.NewUserClient(conn)
	return client
}
