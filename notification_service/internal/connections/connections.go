package connections

import (
	"notify/internal/api/handler"
	"notify/internal/services"
	"sync"

	"github.com/gorilla/websocket"
)

func NewWebSocket() *handler.WebSocket{
	return &handler.WebSocket{
		Map:   make(map[string]*websocket.Conn),
		Mutex: &sync.Mutex{},                    
	}
}

func NewService()*services.Service{
	a:=NewWebSocket()
	return &services.Service{W: a}
}