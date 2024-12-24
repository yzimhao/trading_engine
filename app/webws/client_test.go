package webws

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	_socket     *Hub
	test_symbol = "usdjpy"
)

func init() {

	go func() {
		_socket = NewHub()
		r := gin.New()
		r.Any("/ws", func(ctx *gin.Context) {
			_socket.ServeWs(ctx)
		})
		r.Run(":8090")
	}()
}

func newClient() *websocket.Conn {
	// s := httptest.NewServer(http.HandlerFunc(_socket.ServeWs))
	// defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws://127.0.0.1:8090/ws"
	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatalf("%v", err)
		return nil
	}
	// defer ws.Close()
	return ws
}

func TestClient(t *testing.T) {

	Convey("hello testing", t, func() {
		ws := newClient()
		err := ws.WriteMessage(websocket.TextMessage, []byte("hello"))
		So(err, ShouldBeNil)

		time.Sleep(time.Second * time.Duration(1))
		So(len(_socket.clients), ShouldEqual, 1)
		ws.Close()
		time.Sleep(time.Second * time.Duration(2))
		So(len(_socket.clients), ShouldEqual, 0)
	})

	Convey("客户端添加属性", t, func() {
		ws := newClient()
		defer ws.Close()

		tags := []string{}
		for _, v := range AllWebSocketMsg {
			tags = append(tags, v.Format(map[string]string{
				"symbol": test_symbol,
				"period": "h1",
			}))
		}

		subM := subMessage{
			Subsc: tags,
		}

		body, _ := json.Marshal(subM)

		t.Logf("%+v, body: %s", subM, body)
		if err := ws.WriteMessage(websocket.TextMessage, body); err != nil {
			t.Fatalf("%v", err)
		}

		time.Sleep(time.Second * time.Duration(1))

		So(len(_socket.clients), ShouldEqual, 1)
		for c, _ := range _socket.clients {
			for _, tag := range tags {
				So(c.attrs, ShouldContainKey, tag)
			}
			t.Logf("c.attrs: %#v", c.attrs)
		}

	})
}
