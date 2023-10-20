package order

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/utils"
)

type order_cancel_request_args struct {
	OrderId string `json:"order_id" binding:"required"`
}

func Cancel(ctx *gin.Context) {
	var req order_cancel_request_args
	if err := ctx.BindJSON(&req); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	if err := orders.SubmitOrderCancel(req.OrderId); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	utils.ResponseOkJson(ctx, "")
}
