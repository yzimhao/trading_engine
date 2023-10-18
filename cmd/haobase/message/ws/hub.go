// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ws

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Broadcast chan MsgBody

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	recv chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	sync.Mutex
}

func NewHub() *Hub {
	h := Hub{
		Broadcast:  make(chan MsgBody),
		recv:       make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

	go h.run()
	M = &h
	return &h
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.clients[client] = true
			h.Unlock()

		case client := <-h.unregister:
			h.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				client.attrs = nil
				client.lastSendMsgHash = nil
			}
			h.Unlock()

		case message := <-h.Broadcast:
			go func() {
				h.Lock()
				defer h.Unlock()

				// logrus.Debugf("[wss] broadcast: %+v", message)
				for client := range h.clients {
					//广播消息之前，检查client是否拥有对应属性
					if !client.hasAttr(message.To) {
						continue
					}

					//去重相同两条连续的重复消息
					newHash := message.BodyHash()
					if lastHash, ok := client.lastSendMsgHash[message.Response.Type]; ok {
						if newHash == lastHash {
							continue
						}
					}

					client.lastSendMsgHash[message.To] = message.BodyHash()

					select {
					case client.send <- message.GetBody():
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}()

		case <-h.recv:

		}
	}
}

// serveWs handles websocket requests from the peer.
func (h *Hub) ServeWs(ctx *gin.Context) {

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorf("[wss] upgrader.Upgrade %s", err)
		return
	}

	client := &Client{hub: h, conn: conn, send: make(chan []byte, 256)}
	//注册
	client.hub.register <- client
	client.lastSendMsgHash = make(map[string]string)
	client.attrs = make(map[string]bool)

	//登录信息处理
	// uraw, _ := user.AuthMiddleware.GetClaimsFromJWT(ctx)
	// user_id := int64(0)
	// if _, ok := uraw["user_id"]; ok {
	// 	user_id = int64(uraw["user_id"].(float64))
	// }
	// client.setAttr(fmt.Sprintf("%d", user_id))
	// logrus.Debugf("[wss] login info %#v  attr: %v", uraw, client.attrs)

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

}
