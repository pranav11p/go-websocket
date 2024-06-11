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
	EventChangeRoom  = "change_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Name string `json:"name"`
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
		if client.chatroom == c.chatroom {
			client.egress <- outgoingEvent
		}
	}

	return nil
}

func ChangeRoomHandler(event Event, c *Client) error {
	var changeRoomEvent ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("error un-marshalling ChangeRoom event payload %v", err)
	}

	c.chatroom = changeRoomEvent.Name

	return nil
}
