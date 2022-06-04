package wss

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type msgBody struct {
	Tag  string      `json:"tag"`
	Data interface{} `json:"data"`
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Send(msg []byte) {
	h.broadcast <- msg
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			var body msgBody
			err := json.Unmarshal(message, &body)
			if err == nil {
				msgHash := md5String(message)
				for client := range h.clients {

					if _, ok := client.lastMsgHash[body.Tag]; ok {
						if client.lastMsgHash[body.Tag] == msgHash {
							continue
						}
					}
					client.lastMsgHash[body.Tag] = msgHash

					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

func md5String(str []byte) string {
	hasher := md5.New()
	hasher.Write(str)
	return hex.EncodeToString(hasher.Sum(nil))
}
