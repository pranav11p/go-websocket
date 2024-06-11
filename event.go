package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

func SendMessageHandler(event Event, c *Client) error {

	// Marshal the payload received from FE to the format we want
	var chatEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload received from server: %v", err)
	}

	var broadcastMessage NewMessageEvent
	broadcastMessage.Message = chatEvent.Message
	broadcastMessage.From = chatEvent.From
	broadcastMessage.Sent = time.Now()

	data, err := json.Marshal(broadcastMessage)
	if err != nil {
		return fmt.Errorf("error marshalling message: %v", err)
	}

	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	for client := range c.manager.clients {
		client.egress <- outgoingEvent
	}

	return nil
}
