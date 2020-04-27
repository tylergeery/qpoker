package connection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"qpoker/auth"
	"qpoker/cards/games/holdem"
	"qpoker/models"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// ClientAdminStart is an admin start of game
	ClientAdminStart = "start"

	// ClientChipResponse is admin response to player's chip request
	ClientChipResponse = "chip_response"

	// ClientChipRequest is a player's request for chips
	ClientChipRequest = "chip_request"
)

// ClientEvent represents a player gameplay action
type ClientEvent struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// IsAdminEvent tells whether the event is an admin action
func (e ClientEvent) IsAdminEvent() bool {
	return e.Type == ActionAdmin
}

// IsMsgEvent tells whether the event is an admin action
func (e ClientEvent) IsMsgEvent() bool {
	return e.Type == ActionMsg
}

// ToAdminEvent parses the game action from the GameEvent
func (e ClientEvent) ToAdminEvent(c *Client) AdminEvent {
	adminEvent := AdminEvent{
		Action:   e.Data["action"].(string),
		GameID:   c.GameID,
		PlayerID: c.PlayerID,
		Value:    e.Data["value"],
	}

	switch adminEvent.Action {
	default:
		return adminEvent
	}
}

// ToMsgEvent parses the game action from the GameEvent
func (e ClientEvent) ToMsgEvent(c *Client) MsgEvent {
	return MsgEvent{}
}

// ToGameEvent parses the game action from the GameEvent
func (e ClientEvent) ToGameEvent(c *Client) GameEvent {
	fmt.Printf("Got game event: %s\n", e.Data)
	return GameEvent{
		GameID:   c.GameID,
		PlayerID: c.PlayerID,
		Action: holdem.Action{
			Name:   e.Data["name"].(string),
			Amount: e.Data["amount"].(int64),
		},
	}
}

// Client holds connection information
type Client struct {
	conn     *websocket.Conn
	connOpen bool
	PlayerID int64
	GameID   int64

	GameChannel    chan GameEvent
	AdminChannel   chan AdminEvent
	MessageChannel chan MsgEvent
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

	fmt.Printf("Client auth token: %s\n", string(authEvent.Token))
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
			return
		}

		var event ClientEvent
		err = json.Unmarshal([]byte(msg), &event)
		if err != nil {
			fmt.Printf("Could not unmarshal json action: %s\n", err)
			continue
		}

		if event.IsAdminEvent() {
			c.AdminChannel <- event.ToAdminEvent(c)
			continue
		}

		if event.IsMsgEvent() {
			c.MessageChannel <- event.ToMsgEvent(c)
			continue
		}

		action := event.ToGameEvent(c)
		c.GameChannel <- action
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
		fmt.Printf("Client write error: %s %s\n", err, string(msg))
		return err
	}

	return nil
}
