package connections

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	notification17 "user/internal/client/notification"
	"user/internal/config"
	"user/internal/database/adjust"
	interface17 "user/internal/interface"
	"user/internal/interface/service"
	kafkaconsumer "user/internal/kafka/consumer"
	adjustservice "user/internal/service/adjust"
	grpcmethods "user/internal/service/methods"

	_ "github.com/lib/pq"
)

func NewDatabase() interface17.User {
	c := config.Configuration()
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", c.Database.User, c.Database.Password, c.Database.Host, c.Database.DBname))
	if err != nil {
		log.Println(err)
	}
	if err := db.Ping(); err != nil {
		log.Println(err)
	}
	n := notification17.Hotel()
	return &adjust.Database{Db: db, N: n}
}

func NewService() *service.Service {
	a := NewDatabase()
	return &service.Service{D: a}
}

func NewAdjust() interface17.Adjust {
	a := NewService()
	return &adjustservice.Adjust{S: a}
}

func NewAdjustService() *service.Adjust {
	a := NewAdjust()
	return &service.Adjust{A: a}
}

func NewGrpc() *grpcmethods.Service {
	a := NewAdjustService()
	return &grpcmethods.Service{S: a}
}

func NewConsumer() *kafkaconsumer.Consumer17 {
	a := NewGrpc()
	ctx := context.Background()
	return &kafkaconsumer.Consumer17{C: a, Ctx: ctx}
}
