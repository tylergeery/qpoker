package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"qpoker/auth"
	"qpoker/models"
	"time"

	"github.com/gorilla/websocket"
)

// Client holds connection information
type Client struct {
	conn     *websocket.Conn
	connOpen bool
	PlayerID int64
	GameID   int64
}

// NewClient returns a new client
func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn: conn}
}

// Active returns whether the client is still active
func (c *Client) Active() bool {
	return c.connOpen
}

// HandleShutdown prepares for when the client goes away
func (c *Client) HandleShutdown() {
	c.connOpen = true
	c.conn.SetCloseHandler(func(code int, text string) error {
		c.connOpen = false
		fmt.Printf("Websocket connection closed: %s\n", text)

		return nil
	})
	c.conn.SetReadDeadline(time.Now().Add(10 * time.Minute)) // TODO: game duration
}

// Authenticate ensures that the player first sends a valid token
func (c *Client) Authenticate() error {
	var authEvent AuthEvent

	msg, err := c.getMessage()
	if err != nil {
		return err
	}

	json.Unmarshal([]byte(msg), &authEvent)

	playerID, err := auth.GetPlayerIDFromAccessToken(authEvent.Token)
	if err != nil {
		fmt.Printf("Client authentication token read error: %s\n", err)
		return err
	}

	game, err := models.GetGameBy("id", authEvent.GameID)
	if err != nil {
		fmt.Printf("Client authentication game fetch error: %s\n", err)
		return err
	}

	c.GameID = game.ID
	c.PlayerID = playerID

	return nil
}

// ReadMessages listens for messages until client disappears
func (c *Client) ReadMessages() {
	for c.connOpen {
		// TODO: play with read timeout
		msg, err := c.getMessage()
		if err != nil {
			continue
		}

		// TODO: turn msg into event
		fmt.Println(msg)
	}
}

func (c *Client) getMessage() (s string, err error) {
	_, input, err := c.conn.ReadMessage()
	if err != nil {
		fmt.Printf("Client read error: %s\n", err)
		return
	}

	s = string(bytes.TrimRight(input, "\x00"))

	return
}

// SendMessage sends a message to websocket client
func (c *Client) SendMessage(msg []byte) error {
	err := c.conn.WriteMessage(websocket.TextMessage, msg)
	if err != nil {
		fmt.Printf("Client write error: %s\n", string(msg))
		return err
	}

	return nil
}
