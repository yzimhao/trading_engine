package asset

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"go.uber.org/zap"
)

type AssetModule struct {
	logger *zap.Logger
	router *provider.Router
}

func NewAssetModule(logger *zap.Logger, router *provider.Router) *AssetModule {
	asset := AssetModule{
		logger: logger,
		router: router,
	}
	asset.registerRouter()
	return &asset
}

func (a *AssetModule) registerRouter() {
	assetGroup := a.router.APIv1.Group("/asset")
	assetGroup.GET("/", a.query)
	assetGroup.GET("/:symbol", a.detail)

}

func (a *AssetModule) query(c *gin.Context) {
	//TODO implement
	a.router.ResponseOk(c, nil)
}

func (a *AssetModule) detail(c *gin.Context) {
	//TODO implement
	a.router.ResponseOk(c, nil)
}
