package booking

import (
	"log"
	"notify/internal/protos/booking"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Hotel() booking.BookHotelClient {
	conn, err := grpc.NewClient("booking_service:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	client := booking.NewBookHotelClient(conn)
	return client
}
