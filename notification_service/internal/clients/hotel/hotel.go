package hotel

import (
	"log"
	"notify/internal/protos/hotel"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Hotel() hotel.HotelClient {
	conn, err := grpc.NewClient("hotel_service:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := hotel.NewHotelClient(conn)
	return client
}
