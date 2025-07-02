package assets

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"userAssets",
	fx.Invoke(newUserAssetsModule),
)

type userAssetsModule struct {
	router *provider.Router
	logger *zap.Logger
}

func newUserAssetsModule(
	router *provider.Router,
	logger *zap.Logger,
) {
	asset := userAssetsModule{
		router: router,
		logger: logger,
	}
	asset.registerRouter()
}

func (a *userAssetsModule) registerRouter() {
	ua := a.router.APIv1.Group("/user/asset")
	//TODO 权限认证
	ua.POST("/despoit", a.despoit)
	ua.POST("/withdraw", a.withdraw)
	ua.GET("/query", a.query)
	ua.GET("/:symbol/history", a.queryAssetHistory)
	ua.POST("/transfer/:symbol", a.assetTransfer)
}

func (a *userAssetsModule) despoit(c *gin.Context) {
	//TODO
}

func (a *userAssetsModule) withdraw(c *gin.Context) {
	//TODO
}

func (a *userAssetsModule) query(c *gin.Context) {
	//TODO
}

func (a *userAssetsModule) queryAssetHistory(c *gin.Context) {
	//TODO
}

func (a *userAssetsModule) assetTransfer(c *gin.Context) {
	//TODO
}
