package main

// Event is the event broadcasted to all clients
type Event struct {
	State string `json:"state"`
}
