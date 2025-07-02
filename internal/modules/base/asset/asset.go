package asset

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"go.uber.org/zap"
)

type assetModule struct {
	logger    *zap.Logger
	router    *provider.Router
	assetRepo persistence.AssetRepository
}

func newAssetModule(logger *zap.Logger, router *provider.Router, repo persistence.AssetRepository) {
	asset := assetModule{
		logger:    logger,
		router:    router,
		assetRepo: repo,
	}
	asset.registerRouter()
}

func (a *assetModule) registerRouter() {
	assetGroup := a.router.APIv1.Group("/asset")
	assetGroup.GET("/", a.query)
	assetGroup.GET("/:symbol", a.detail)

}

func (a *assetModule) query(c *gin.Context) {
	//TODO implement
	a.router.ResponseOk(c, nil)
}

func (a *assetModule) detail(c *gin.Context) {
	//TODO implement
	a.router.ResponseOk(c, nil)
}
