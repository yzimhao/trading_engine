package orders

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"userOrders",
	fx.Invoke(newUserOrdersModule),
)

type userOrderModule struct {
	router *provider.Router
	logger *zap.Logger
}

func newUserOrdersModule(
	router *provider.Router,
	logger *zap.Logger,
) {
	uo := userOrderModule{
		router: router,
		logger: logger,
	}
	uo.registerRouter()
}

func (u *userOrderModule) registerRouter() {
	uo := u.router.APIv1.Group("/user/order")
	//TODO 权限认证
	uo.GET("/history", u.orderHistory)
	uo.GET("/unfinished", u.unfinishedList)
	uo.GET("/trade/history", u.tradeHistory)
}

// TODO implement
func (u *userOrderModule) orderHistory(c *gin.Context) {}

// TODO implement
func (u *userOrderModule) tradeHistory(c *gin.Context) {}

// TODO implement
func (u *userOrderModule) unfinishedList(c *gin.Context) {}
