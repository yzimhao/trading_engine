package asset

import (
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"go.uber.org/zap"
)

type AssetModule struct {
	logger    *zap.Logger
	router    *provider.Router
	assetRepo persistence.AssetRepository
}

func NewAssetModule(logger *zap.Logger, router *provider.Router, repo persistence.AssetRepository) *AssetModule {
	asset := AssetModule{
		logger:    logger,
		router:    router,
		assetRepo: repo,
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
