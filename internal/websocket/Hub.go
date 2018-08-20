package websocket

import (
	"encoding/json"
	"log"
	"takeaway/takeaway-server/internal/vote"
)

// HubInstance is a singleton instance of the Hub struct.
var HubInstance = newHub()

// Hub provides a collection for managing all websocket connections to the server,
// based off https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go.
type Hub struct {
	Broadcast  chan *vote.Poll
	clients    map[string][]*Client
	register   chan *Client
	unregister chan *Client
}

// Run starts the event loop for the given Hub object, note this method will block so should be ran as a seporate goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.ID] == nil {
				h.clients[client.ID] = make([]*Client, 0)
			}
			h.clients[client.ID] = append(h.clients[client.ID], client)
		case client := <-h.unregister:
			h.delete(client)
		case poll := <-h.Broadcast:
			data, err := json.Marshal(poll)
			if err != nil {
				log.Printf("Hub: could not marshal poll object due to, %s \n", err)
			}
			// send given data to all registered clients for the specified poll's ID.
			for _, c := range h.clients[poll.ID] {
				// attempt to place data onto given client's send channel, should this not be possible assume the given
				// client is not valid and unregister it.
				select {
				case c.send <- []byte(data):
				default:
					h.delete(c)
				}
			}
		}
	}
}
func (h *Hub) delete(client *Client) {
	for i, c := range h.clients[client.ID] {
		if c == client {
			h.clients[client.ID] = append(h.clients[client.ID][:i], h.clients[client.ID][i+1:]...)
		}
	}
	close(client.send)
}

// newHub create a brand new Hub objecct instance, returning a pointer to said instance.
func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan *vote.Poll),
		clients:    make(map[string][]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// NotifyChange is a utility method allowing an updated poll object to be broadcasted using the HubInstance.
func NotifyChange(p *vote.Poll) {
	HubInstance.Broadcast <- p
}
