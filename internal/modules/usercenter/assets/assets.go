package assets

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
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
	auth           *middlewares.AuthMiddleware
}

func newUserAssetsModule(
	router *provider.Router,
	logger *zap.Logger,
	auth *middlewares.AuthMiddleware,
	userAssetsRepo persistence.UserAssetRepository,
) {
	asset := userAssetsModule{
		router:         router,
		logger:         logger,
		userAssetsRepo: userAssetsRepo,
		auth:           auth,
	}
	asset.registerRouter()
}

func (a *userAssetsModule) registerRouter() {
	ua := a.router.APIv1.Group("/user/asset")
	//权限认证
	ua.Use(a.auth.Auth())

	//内部接口，充值接口
	ua.POST("/despoit", a.despoit)
	//内部接口，提现接口
	ua.POST("/withdraw", a.withdraw)

	//用户资产查询接口
	ua.GET("/query", a.query)
	//用户某个资产的历史记录
	ua.GET("/:symbol/history", a.queryAssetHistory)
	//用户资产转移接口
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

	symbols := strings.ToLower(c.Query("symbols"))
	symbolsSlice := strings.Split(symbols, ",")

	assets, err := a.userAssetsRepo.QueryUserAssets(userId, symbolsSlice...)
	if err != nil {
		a.logger.Sugar().Errorf("userAssets query error: %v", err)
		a.router.ResponseError(c, types.ErrInternalError)
		return
	}

	var response []any
	for _, item := range assets {
		response = append(response, gin.H{
			"symbol":         item.Symbol,
			"total_balance":  item.TotalBalance.StringFixedBank(4),
			"freeze_balance": item.FreezeBalance.StringFixedBank(4),
			"avail_balance":  item.AvailBalance.StringFixedBank(4),
			"updated_at":     item.UpdatedAt,
		})
	}
	a.router.ResponseOk(c, response)
}

func (a *userAssetsModule) queryAssetHistory(c *gin.Context) {
	//TODO
}

func (a *userAssetsModule) assetTransfer(c *gin.Context) {
	//TODO
}
