package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/notification/ws"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"notification",
	fx.Provide(
		ws.NewWsManager,
	),
	fx.Invoke(run),
)

func run(router *provider.Router, wsm *ws.WsManager) {
	router.GET("/ws", func(c *gin.Context) {
		wsm.Listen(c.Writer, c.Request, c.Request.Header)
	})
}
