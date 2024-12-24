// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webws

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	//TODO 通过配置文件修改该值
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type client struct {
	m *WsManager
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	//客户端链接属性标记，max 250个
	attrs map[string]bool

	//服务端推送消息的hash，用来去重,每一种消息类型单独去重
	lastMessageHash map[string]string

	mx sync.Mutex
}

func newClient(m *WsManager, conn *websocket.Conn) *client {
	c := client{
		m:               m,
		conn:            conn,
		send:            make(chan []byte),
		attrs:           make(map[string]bool),
		lastMessageHash: make(map[string]string),
	}
	return &c
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {
	defer func() {
		c.m.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// app.Logger.Errorf("[wss ] IsUnexpectedCloseError: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.m.recv <- message

		c.handleRecvData(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *client) handleRecvData(body []byte) {
	var msg RecviceTag
	err := json.Unmarshal(body, &msg)
	if err != nil {
		return
	}

	//新增订阅属性处理
	for _, attr := range msg.Subscribe {
		if strings.HasPrefix(attr, "_") {
			//带有"_"标记的tag只能是内部程序设置的，不能通过前端发送过来指定
			continue
		}
		if strings.HasPrefix(attr, "token.") {
			// a := strings.Split(attr, ".")
			// _token := a[1]
			// user_id := token.Get(_token)
			// if user_id != "" {
			// 	c.setAttr(types.MsgUser.Format(map[string]string{
			// 		"user_id": user_id,
			// 	}))
			// }
		} else {
			c.setAttr(attr)
		}
	}
	//取消订阅属性处理
	for _, attr := range msg.Unsubscribe {
		c.delAttr(attr)
	}

	// app.Logger.Debugf("[wss] recv: %v attrs: %v", msg, c.attrs)
}

func (c *client) setAttr(tag string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.attrs[tag] = true
}

func (c *client) hasAttr(tag string) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	if _, ok := c.attrs[tag]; ok {
		return true
	}
	return false
}

func (c *client) delAttr(tag string) bool {
	c.mx.Lock()
	defer c.mx.Unlock()

	delete(c.attrs, tag)
	return true
}
