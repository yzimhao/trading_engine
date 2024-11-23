package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"github.com/yzimhao/trading_engine/v2/internal/persistence"
)

type UserController struct {
	assetRepo persistence.AssetRepository
}

type inUserContext struct {
	fx.In
	AssetRepo persistence.AssetRepository
}

func NewUserController(in inUserContext) *UserController {
	return &UserController{
		assetRepo: in.AssetRepo,
	}
}

func (u *UserController) DemoLogin(ctx *gin.Context) {
	//随机一个用户ID
	// 充值usd和jpy两种货币，返回这个用户的基本信息

	// userId := time.Now().Unix()
	// u.assetRepo.Despoit()

	ctx.JSON(http.StatusOK, gin.H{"message": "demo login"})
}
