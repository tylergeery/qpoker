package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"qpoker/websocket/bus"
	"qpoker/websocket/connection"
	"qpoker/websocket/events"

	"github.com/gorilla/websocket"
)

var eventBus *bus.EventBus

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	eventBus = bus.StartEventBus()

	http.HandleFunc("/", handleSocketConnection)
	http.ListenAndServe(":8080", nil)
}

func handleSocketConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

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
	client := connection.NewClient(conn)
	client.HandleShutdown()
	err = client.Authenticate()
	if err != nil {
		fmt.Printf("Client authentication error: %s\n", err)
		return
	}

	fmt.Println("Sending player channel event")
	eventBus.PlayerChannel <- events.PlayerEvent{client, events.ActionPlayerRegister}

	// client will spin here until disconnected
	client.ReadMessages(eventBus.HandleClientEvent)

	eventBus.PlayerChannel <- events.PlayerEvent{client, events.ActionPlayerLeave}
	fmt.Println("Websocket request terminated")
}
