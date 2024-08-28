package router

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"notify/internal/config"
	"notify/internal/connections"
	"notify/internal/protos/notification"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewRouter() {
	r := mux.NewRouter()
	a := connections.NewService().W
	r.HandleFunc("/ws", a.HandleWebSocket)
	go Grpc()
	fmt.Println("server started on port 8083")
	log.Fatal(http.ListenAndServe("notification_service:8083", r))
}

func Grpc() {
	c := config.Configuration()
	ls, err := net.Listen(c.User.Host, c.User.Port)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	server := connections.NewService()
	notification.RegisterNotificationServer(s, server)
	reflection.Register(s)
	// a := connections.NewConsumer()
	// go func() {
	// 	a.Consumer()
	// }()
	fmt.Printf("server started on the port %s", c.User.Port)

	if err := s.Serve(ls); err != nil {
		log.Fatal(err)
	}
}