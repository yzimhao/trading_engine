package ws

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"go.uber.org/zap"
)

var (
	manager     *WsManager
	test_symbol = "usdjpy"
)

func init() {

	go func() {
		logger, _ := zap.NewDevelopment()
		v := provider.NewViper(logger)
		r := provider.NewRedis(v, logger)
		produce := provider.NewProduce(r, logger)
		consume := provider.NewConsume(r, logger)
		manager = NewWsManager(logger, produce, consume)
		router := gin.New()
		router.Any("/ws", func(ctx *gin.Context) {
			manager.Listen(ctx.Writer, ctx.Request, ctx.Request.Header)
		})
		if err := router.Run(":8090"); err != nil {
			log.Fatalf("router.Run failed: %v", err)
		}
	}()
}

func clientConn() *websocket.Conn {
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
		ws := clientConn()
		err := ws.WriteMessage(websocket.TextMessage, []byte("hello"))
		So(err, ShouldBeNil)

		time.Sleep(time.Second * time.Duration(1))
		So(len(manager.membersMap), ShouldEqual, 1)
		if err := ws.Close(); err != nil {
			t.Errorf("ws.Close failed: %v", err)
		}
		time.Sleep(time.Second * time.Duration(2))
		So(len(manager.membersMap), ShouldEqual, 0)
	})

	Convey("客户端添加属性", t, func() {
		ws := clientConn()
		defer func() {
			if err := ws.Close(); err != nil {
				t.Errorf("ws.Close failed: %v", err)
			}
		}()

		tags := []string{}
		for _, v := range AllWebSocketMsg {
			tags = append(tags, v.Format(map[string]string{
				"symbol": test_symbol,
				"period": "h1",
			}))
		}

		subM := RecviceTag{
			Subscribe: tags,
		}

		body, _ := json.Marshal(subM)

		t.Logf("%+v, body: %s", subM, body)
		if err := ws.WriteMessage(websocket.TextMessage, body); err != nil {
			t.Fatalf("%v", err)
		}

		time.Sleep(time.Second * time.Duration(1))

		So(len(manager.membersMap), ShouldEqual, 1)
		for c := range manager.membersMap {
			for _, tag := range tags {
				So(c.hasAttr(tag), ShouldBeTrue)
			}
			t.Logf("c.attrs: %#v", c)
		}

	})
}
