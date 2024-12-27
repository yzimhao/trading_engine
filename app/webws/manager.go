package webws

import (
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type WsManager struct {
	Broadcast  chan Message
	logger     *zap.Logger
	recv       chan []byte
	register   chan *client
	unregister chan *client
	members    map[*client]bool
	mx         sync.Mutex
}

func NewWsManager(logger *zap.Logger) *WsManager {
	m := WsManager{
		Broadcast:  make(chan Message),
		logger:     logger,
		recv:       make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		members:    make(map[*client]bool),
	}

	go m.run()

	return &m
}

func (m *WsManager) Listen(writer http.ResponseWriter, req *http.Request, responseHeader http.Header) {
	conn, err := upgrader.Upgrade(writer, req, nil)
	if err != nil {
		m.logger.Sugar().Errorf("webws upgrader.Upgrade %v", err)
		return
	}

	cli := newClient(m, conn)
	cli.m.register <- cli
	go cli.writePump()
	go cli.readPump()
}

func (m *WsManager) Members() []*client {
	m.mx.Lock()
	defer m.mx.Unlock()

	var members []*client
	for c := range m.members {
		members = append(members, c)
	}
	return members
}

func (m *WsManager) ClientHasAttr(client *client, tag string) bool {
	return client.hasAttr(tag)
}

func (m *WsManager) GetClientAttrs(client *client) map[string]bool {
	client.mx.Lock()
	defer client.mx.Unlock()

	return client.attrs
}

func (m *WsManager) run() {
	for {
		select {
		case cli := <-m.register:

			func(client *client) {
				m.mx.Lock()
				defer m.mx.Unlock()

				m.logger.Sugar().Debugf("[wss] register client: %v", cli)
				m.members[client] = true
			}(cli)

		case cli := <-m.unregister:
			func(client *client) {
				m.mx.Lock()
				defer m.mx.Unlock()

				if _, ok := m.members[client]; ok {
					delete(m.members, client)
					close(client.send)

					client.attrs = nil
					client.lastMessageHash = nil
				}
			}(cli)

		case msg := <-m.Broadcast:
			go func(message Message) {
				m.mx.Lock()
				defer m.mx.Unlock()

				m.logger.Sugar().Debugf("[wss] broadcast message: %v", message)

				for client := range m.members {
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

					client.lastMessageHash[message.To] = sign
					client.send <- message.Body()
				}
			}(msg)

		case data := <-m.recv:
			m.logger.Sugar().Debugf("[wss] recivce message: %s", data)
		}
	}
}
