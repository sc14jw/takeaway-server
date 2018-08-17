package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Time allowed to send a message to a client.
const writeWait = 10 * time.Second

// Client represents a singular websocket connection to the server from a client,
// based off https://github.com/gorilla/websocket/blob/master/examples/chat/client.go.
type Client struct {
	send chan []byte
	conn *websocket.Conn
	ID   string
}

func (c *Client) writeLoop() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// hub closed the connection
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Client - Could not create writer to send message: %v due to: %s", message, err)
				return
			}

			w.Write(message)

			// Write any queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				log.Printf("Client - Couldn't close writer due to: %s", err)
				return
			}
		}

	}
}
