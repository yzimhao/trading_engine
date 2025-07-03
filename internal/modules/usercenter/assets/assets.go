package assets

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"userAssets",
	fx.Invoke(newUserAssetsModule),
)

type userAssetsModule struct {
	router         *provider.Router
	logger         *zap.Logger
	userAssetsRepo persistence.UserAssetRepository
}

func newUserAssetsModule(
	router *provider.Router,
	logger *zap.Logger,
	userAssetsRepo persistence.UserAssetRepository,
) {
	asset := userAssetsModule{
		router:         router,
		logger:         logger,
		userAssetsRepo: userAssetsRepo,
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

// @Summary 用户持仓资产
// @Description 获取用户持仓资产接口
// @ID v1.user.assets.query
// @Tags 用户中心
// @Accept json
// @Produce json
// @Param symbols query string true "symbols"
// @Success 200 {string} any
// @Router /api/v1/user/assets/query [post]
func (a *userAssetsModule) query(c *gin.Context) {
	userId := a.router.ParseUserID(c)

	symbols := c.Query("symbols")
	symbolsSlice := strings.Split(symbols, ",")

	assets, err := a.userAssetsRepo.QueryUserAssets(userId, symbolsSlice...)
	if err != nil {
		a.logger.Sugar().Errorf("userAssets query error: %v", err)
		a.router.ResponseError(c, types.ErrInternalError)
		return
	}
	a.router.ResponseOk(c, assets)
}

func (a *userAssetsModule) queryAssetHistory(c *gin.Context) {
	//TODO
}

func (a *userAssetsModule) assetTransfer(c *gin.Context) {
	//TODO
}
