package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/twmb/franz-go/pkg/kgo"
)

type WebSocket struct {
	Map   map[string]*websocket.Conn
	Mutex *sync.Mutex
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (u *WebSocket) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket is working")
	userID := r.Header.Get("id")
	fmt.Println(userID)
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Current connections map:", u.Map)

	kafkaReader, err := kgo.NewClient(
		kgo.SeedBrokers("localhost:9092"),
		kgo.ConsumeTopics("notification"),
		kgo.ConsumerGroup("my_group"),
	)
	if err != nil {
		log.Println("Kafka client creation error:", err)
		return
	}
	defer kafkaReader.Close()
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	if _, exists := u.Map[userID]; exists {
		log.Printf("User %s is already connected", userID)
	}
	u.Map[userID] = conn
	log.Printf("User %s added to the map", userID)
	if err := u.AddUser(userID, conn); err != nil {
		log.Println("Error adding user:", err)
		return
	}

	for {
		fetches := kafkaReader.PollFetches(r.Context())
		if fetches.IsClientClosed() {
			break
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
				message := record.Value
				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Println("Error writing message to WebSocket:", err)
					return
				}
				log.Printf("Sent message to user %s: %s", userID, string(message))
			}
		}
	}

func (u *WebSocket) AddUser(userID string, conn *websocket.Conn) error {
	u.Mutex.Lock()
	defer u.Mutex.Unlock()
	if _, exists := u.Map[userID]; exists {
		log.Printf("User %s is already connected", userID)
		return errors.New("user already exists")
	}
	u.Map[userID] = conn
	log.Printf("User %s added to the map", userID)
	return nil
}

// func (u *WebSocket) RemoveUser(userID string) {
// 	u.Mutex.Lock()
// 	defer u.Mutex.Unlock()
// 	delete(u.Map, userID)
// 	log.Printf("User %s disconnected", userID)
// }
