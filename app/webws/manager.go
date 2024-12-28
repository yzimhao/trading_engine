package webws

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/duolacloud/broker-core"
	"go.uber.org/zap"
)

const (
	websocketMsg = "websocket_msg"
)

type WsManager struct {
	broadcast  chan Message
	logger     *zap.Logger
	recv       chan []byte
	register   chan *client
	unregister chan *client
	membersMap map[*client]bool
	mx         sync.Mutex
	broker     broker.Broker
	debug      bool
}

func NewWsManager(logger *zap.Logger, broker broker.Broker) *WsManager {
	m := WsManager{
		broadcast:  make(chan Message),
		logger:     logger,
		recv:       make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		membersMap: make(map[*client]bool),
		broker:     broker,
	}

	go m.run()
	go m.subscribe()

	return &m
}

func (m *WsManager) SetDebug(v bool) {
	m.debug = v
}

// example:
// 广播给所有订阅某个type的用户
// Broadcast("depth.usdjpy", "")
// Broadcast("trade.usdjpy", "")

func (m *WsManager) Broadcast(ctx context.Context, _type string, body any) error {
	return m.send(ctx, _type, _type, body)
}

// 单发消息给某一个用户
// SendTo("1001", "order.new.usdjpy", "")
// SendTo("1001", "order.cancel.usdjpy", "")
func (m *WsManager) SendTo(ctx context.Context, uid string, _type string, body any) error {
	to := fmt.Sprintf("_user.%s", uid)
	return m.send(ctx, to, _type, body)
}

func (m *WsManager) send(ctx context.Context, to string, _type string, body any) error {
	msg := NewMessage(to, _type, body)
	_body, err := msg.Marshal()
	if err != nil {
		m.logger.Sugar().Errorf("[ws] msg.Marshal error %v", err)
		return err
	}

	err = m.broker.Publish(ctx, websocketMsg, &broker.Message{
		Body: _body,
	})

	if err != nil {
		return err
	}

	return nil
}

func (m *WsManager) Listen(writer http.ResponseWriter, req *http.Request, responseHeader http.Header) {
	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		m.logger.Sugar().Errorf("[ws] upgrader.Upgrade %v", err)
		return
	}

	cli := newClient(m, conn)
	cli.m.register <- cli
	go cli.writePump()
	go cli.readPump()
}

// TODO 多个websocket节点的时候，订阅模式要修改
func (m *WsManager) subscribe() {
	m.broker.Subscribe(websocketMsg, func(ctx context.Context, event broker.Event) error {
		if m.debug {
			m.logger.Sugar().Debugf("[ws] broker subscribe message: %s", event.Message().Body)
		}

		var msg Message
		err := json.Unmarshal(event.Message().Body, &msg)
		if err != nil {
			m.logger.Sugar().Errorf("[ws] websocket message unmarshal error: %v", err)
			return err
		}

		m.broadcast <- msg
		return nil
	})
}

func (m *WsManager) run() {
	for {
		select {
		case cli := <-m.register:

			func(client *client) {
				m.mx.Lock()
				defer m.mx.Unlock()

				m.logger.Sugar().Debugf("[ws] register client: %v", cli)
				m.membersMap[client] = true
			}(cli)

		case cli := <-m.unregister:
			func(client *client) {
				m.mx.Lock()
				defer m.mx.Unlock()

				if _, ok := m.membersMap[client]; ok {
					delete(m.membersMap, client)
					close(client.send)

					client.attrs = nil
					client.lastMessageHash = nil
				}
			}(cli)

		case msg := <-m.broadcast:
			go func(message Message) {
				m.mx.Lock()
				defer m.mx.Unlock()

				if m.debug {
					m.logger.Sugar().Debugf("[ws] broadcast message: %v", message)
				}

				for client := range m.membersMap {
					if !client.hasAttr(message.To) {
						continue
					}

					//去重相同两条连续的重复消息
					sign := message.Sign()
					if lastHash, ok := client.lastMessageHash[message.Response.Type]; ok {
						if sign == lastHash {
							continue
						}
					}

					if m.debug {
						m.logger.Sugar().Infof("[ws] send to %s body: %v", message.To, message.Response)
					}

					client.lastMessageHash[message.To] = sign
					client.send <- message.ResponseBytes()
				}
			}(msg)

		case data := <-m.recv:
			m.logger.Sugar().Debugf("[ws] recivce message: %s", data)
		}
	}
}
