package www

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/utils"
)

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
