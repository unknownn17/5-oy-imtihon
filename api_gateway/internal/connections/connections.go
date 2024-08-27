package connections

import (
	"api/internal/api/handler"
	broadcast17 "api/internal/broadcast"
	"api/internal/clients/booking"
	hotelservice "api/internal/clients/hotel"
	userservice "api/internal/clients/user"
	redismethod "api/internal/redis/method"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewBroadcast() *broadcast17.Adjust {
	u := userservice.UserClinet()
	h:=hotelservice.Hotel()
	b:=booking.Hotel()
	r := Redis()
	ctx := context.Background()
	return &broadcast17.Adjust{U: u, Ctx: ctx, R: r,H: h,B: b}
}

func NewHandler() *handler.Handler {
	a := NewBroadcast()
	return &handler.Handler{B: a}
}

func Redis() *redismethod.Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	return &redismethod.Redis{R: client, Ctx: ctx}
}
