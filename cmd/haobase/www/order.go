package www

import (
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/utils"
)

type order_create_request_args struct {
	User string `json:"user"` //test时使用

	Symbol    string `json:"symbol" binding:"required" example:"eurusd"`
	Side      string `json:"side" binding:"required" example:"sell/buy"`
	OrderType string `json:"order_type" binding:"required" example:"limit/market"`
	Price     string `json:"price" example:"1.00"`
	Quantity  string `json:"qty" example:"12"`
	Amount    string `json:"amount" example:"100.00"`
}

func order_create(ctx *gin.Context) {
	var req order_create_request_args
	if err := ctx.BindJSON(&req); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	var info *orders.Order
	var err error

	side := trading_core.OrderSide(req.Side)

	if req.OrderType == trading_core.OrderTypeLimit.String() {
		info, err = orders.NewLimitOrder(req.User, req.Symbol, side, req.Price, req.Quantity)
	} else if req.OrderType == trading_core.OrderTypeMarket.String() {
		if utils.D(req.Amount).Cmp(decimal.Zero) > 0 {
			info, err = orders.NewMarketOrderByAmount(req.User, req.Symbol, side, req.Amount)
		} else if utils.D(req.Quantity).Cmp(decimal.Zero) > 0 {
			info, err = orders.NewMarketOrderByQty(req.User, req.Symbol, side, req.Quantity)
		}
	}
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	utils.ResponseOkJson(ctx, info.OrderId)
}

type order_cancel_request_args struct {
	OrderId string `json:"order_id" binding:"required"`
}

func order_cancel(ctx *gin.Context) {
	var req order_cancel_request_args
	if err := ctx.BindJSON(&req); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	//todo
	utils.ResponseOkJson(ctx, "")
}
