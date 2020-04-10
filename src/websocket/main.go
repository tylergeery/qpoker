package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var eventBus *EventBus

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	eventBus = NewEventBus()

	go eventBus.ListenForEvents()

	http.HandleFunc("/", handleSocketConnection)
	http.ListenAndServe(":8080", nil)
}

func handleSocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	fmt.Println("Websocket request received")
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Websocket upgrade error: %s", err.Error())
		return
	}

	defer conn.Close()

	// create client
	client := NewClient(conn)
	client.HandleShutdown()
	err = client.Authenticate()
	if err != nil {
		fmt.Printf("Client authentication error: %s\n", err)
		return
	}

	eventBus.PlayerChannel <- PlayerEvent{client, ActionPlayerRegister}

	// client will spin here until disconnected
	client.ReadMessages()

	eventBus.PlayerChannel <- PlayerEvent{client, ActionPlayerLeave}

	fmt.Println("Websocket request terminated")
}
