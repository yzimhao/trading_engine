package internal_api

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/utils"
)

type req_deposit_withdraw_args struct {
	OrderId string `json:"order_id" binding:"required"`
	UserId  string `json:"user_id" binding:"required"`
	Symbol  string `json:"symbol" binding:"required"`
	Amount  string `json:"amount" binding:"required"`
}

func Deposit(ctx *gin.Context) {
	var req req_deposit_withdraw_args
	if err := ctx.BindJSON(&req); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	if utils.D(req.Amount).Cmp(utils.D("0")) <= 0 {
		utils.ResponseFailJson(ctx, "金额非法")
		return
	}

	if _, err := base.NewSymbols().Get(req.Symbol); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	if assets.QueryAssetsLogBusIdIsExist(req.UserId, req.OrderId) {
		utils.ResponseFailJson(ctx, "重复充值")
		return
	}

	_, err := assets.SysDeposit(req.UserId, req.Symbol, req.Amount, req.OrderId)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	utils.ResponseOkJson(ctx, gin.H{})
	return
}

func Withdraw(ctx *gin.Context) {
	var req req_deposit_withdraw_args
	if err := ctx.BindJSON(&req); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	if utils.D(req.Amount).Cmp(utils.D("0")) <= 0 {
		utils.ResponseFailJson(ctx, "金额非法")
		return
	}

	if _, err := base.NewSymbols().Get(req.Symbol); err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}

	if assets.QueryAssetsLogBusIdIsExist(req.UserId, req.OrderId) {
		utils.ResponseFailJson(ctx, "重复提现")
		return
	}

	_, err := assets.SysWithdraw(req.UserId, req.Symbol, req.Amount, req.OrderId)
	if err != nil {
		utils.ResponseFailJson(ctx, err.Error())
		return
	}
	utils.ResponseOkJson(ctx, gin.H{})
	return
}
