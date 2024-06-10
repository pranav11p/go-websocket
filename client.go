package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager

	egress chan []byte
}

type ClientList map[*Client]bool

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan []byte),
	}
}

func (c *Client) readMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		log.Println("MessageType :", messageType)
		log.Println("Payload:", string(payload))
		for wsClient := range c.manager.clients {
			// log.Println("Inside add message:", wsClient)
			wsClient.egress <- payload
		}
	}
}

func (c *Client) writeMessage() {
	defer func() {
		c.manager.removeClient(c)
	}()
	log.Println("inside writeMessage")

	for {
		select {
		case message, ok := <-c.egress:
			if !ok {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println(err)
				}
			}

			if err := c.connection.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println(err)
			}
			log.Println("message sent")
		}
	}
}
