package bus

import (
	"qpoker/cards/games/holdem"
	"qpoker/websocket/connection"
	"qpoker/websocket/events"
	"qpoker/websocket/utils"
)

func isAdminEvent(e connection.ClientEvent) bool {
	return e.Type == events.ActionAdmin
}

func isMsgEvent(e connection.ClientEvent) bool {
	return e.Type == events.ActionMsg
}

func isVideoEvent(e connection.ClientEvent) bool {
	return e.Type == events.ActionVideo
}

// ToAdminEvent parses the game action from the GameEvent
func ToAdminEvent(e connection.ClientEvent, c *connection.Client) events.AdminEvent {
	adminEvent := events.AdminEvent{
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
func ToMsgEvent(e connection.ClientEvent, c *connection.Client) events.MsgEvent {
	return events.MsgEvent{
		Value:    e.Data["message"].(string),
		GameID:   c.GameID,
		PlayerID: c.PlayerID,
		Username: e.Data["username"].(string),
	}
}

// ToGameEvent parses the game action from the GameEvent
func ToGameEvent(e connection.ClientEvent, c *connection.Client) events.GameEvent {
	return events.GameEvent{
		GameID:   c.GameID,
		PlayerID: c.PlayerID,
		Action: holdem.Action{
			Name:   e.Data["name"].(string),
			Amount: utils.InterfaceInt64(e.Data["amount"]),
		},
	}
}

// ToVideoEvent parses the video action from the ClientEvent
func ToVideoEvent(e connection.ClientEvent, c *connection.Client) events.VideoEvent {
	videoEvent := events.VideoEvent{
		Type:         e.Data["type"].(string),
		FromPlayerID: c.PlayerID,
		ToPlayerID:   int64(e.Data["to_player_id"].(float64)),
		GameID:       c.GameID,
	}

	if offer, ok := e.Data["offer"]; ok {
		videoEvent.Offer = offer
	}

	if candidate, ok := e.Data["candidate"]; ok {
		videoEvent.Candidate = candidate
	}

	return videoEvent
}

// HandleClientEvent handles parsing/transmitting client events
func (e *EventBus) HandleClientEvent(event connection.ClientEvent, c *connection.Client) {
	switch {
	case isAdminEvent(event):
		e.AdminChannel <- ToAdminEvent(event, c)
	case isMsgEvent(event):
		e.MessageChannel <- ToMsgEvent(event, c)
	case isVideoEvent(event):
		e.VideoChannel <- ToVideoEvent(event, c)
	default:
		e.GameChannel <- ToGameEvent(event, c)
	}
}
