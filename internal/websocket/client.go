package websocket

import (
	"encoding/json"
	"log"

	"avi_bd/internal/dto"

	"github.com/gorilla/websocket"
)

// Client represents a WebSocket client connection
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan interface{}
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan interface{}, 256),
	}
}

// ReadLoop reads messages from the WebSocket connection
func (c *Client) ReadLoop() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(NewDeadline())
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(NewDeadline())
		return nil
	})

	for {
		var event dto.WebSocketEvent
		err := c.conn.ReadJSON(&event)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}
		// Process incoming events if needed
		_ = event
	}
}

// WriteLoop writes messages to the WebSocket connection
func (c *Client) WriteLoop() {
	ticker := NewTicker()
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case event, ok := <-c.send:
			c.conn.SetWriteDeadline(NewDeadline())
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(event)
			if err != nil {
				log.Printf("JSON marshal error: %v", err)
				continue
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(NewDeadline())
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send sends an event to the client
func (c *Client) Send(event *dto.WebSocketEvent) {
	select {
	case c.send <- event:
	default:
		// Skip if send channel is full
	}
}
