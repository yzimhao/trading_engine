package orders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"userOrders",
	fx.Invoke(newUserOrdersModule),
)

type userOrderModule struct {
	router    *provider.Router
	logger    *zap.Logger
	orderRepo persistence.OrderRepository
	auth      *middlewares.AuthMiddleware
}

func newUserOrdersModule(
	router *provider.Router,
	logger *zap.Logger,
	auth *middlewares.AuthMiddleware,
	orderRepo persistence.OrderRepository,
) {
	uo := userOrderModule{
		router:    router,
		logger:    logger,
		orderRepo: orderRepo,
		auth:      auth,
	}
	uo.registerRouter()
}

func (u *userOrderModule) registerRouter() {
	uo := u.router.APIv1.Group("/user/order")
	// 权限认证
	uo.Use(u.auth.Auth())

	// 未完成订单接口
	uo.GET("/unfinished", u.unfinishedList)
	// 历史订单
	uo.GET("/history", u.orderHistory)
	// 成交记录
	uo.GET("/trade/history", u.tradeHistory)
}

// @Summary 历史订单
// @Description history list
// @ID v1.user.order.history
// @Tags 用户中心
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param start query int64 true "start"
// @Param end query int64 true "end"
// @Param limit query int true "limit"
// @Success 200 {string} any
// @Router /api/v1/order/history [get]
func (u *userOrderModule) orderHistory(c *gin.Context) {
	userId := u.router.ParseUserID(c)

	symbol := c.DefaultQuery("symbol", "")
	start, _ := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	end, _ := strconv.ParseInt(c.DefaultQuery("end", "0"), 10, 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, err := u.orderRepo.HistoryList(c, userId, symbol, start, end, limit)
	if err != nil {
		u.router.ResponseError(c, types.ErrDatabaseError)
		return
	}
	u.router.ResponseOk(c, orders)
}

// @Summary 成交历史
// @Description trade history list
// @ID v1.order.trade_history
// @Tags 用户中心
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/user/order/trade/history [get]
func (u *userOrderModule) tradeHistory(c *gin.Context) {
	//TODO
}

// @Summary 未成交的订单
// @Description unfinished list
// @ID v1.user.order.unfinished
// @Tags 用户中心
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Success 200 {string} any
// @Router /api/v1/user/order/unfinished [get]
func (u *userOrderModule) unfinishedList(c *gin.Context) {
	symbol := c.Query("symbol")
	// userId := common.GetUserId(c)

	// tradeVariety, err := ctrl.tradeVariety.FindBySymbol(c, symbol)
	// if err != nil {
	// 	common.ResponseError(c, err)
	// 	return
	// }

	//TODO 这个未完成订单列表 需要重新写一个方法
	orders, err := u.orderRepo.LoadUnfinishedOrders(c, symbol)
	if err != nil {
		u.router.ResponseError(c, types.ErrDatabaseError)
		return
	}
	u.router.ResponseOk(c, orders)
}
