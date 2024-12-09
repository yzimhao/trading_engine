package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/duolacloud/crud-core/cache"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/common"
	"github.com/yzimhao/trading_engine/v2/internal/modules/matching"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MarketController struct {
	logger    *zap.Logger
	cache     cache.Cache
	klineRepo persistence.KlineRepository
}

type inMarketContext struct {
	fx.In
	Logger    *zap.Logger
	Cache     cache.Cache
	KlineRepo persistence.KlineRepository
}

func NewMarketController(in inMarketContext) *MarketController {
	return &MarketController{
		logger:    in.Logger,
		cache:     in.Cache,
		klineRepo: in.KlineRepo,
	}
}

// @Summary depth
// @Description get depth
// @ID v1.market.depth
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/market/depth [get]
func (ctrl *MarketController) Depth(c *gin.Context) {
	symbol := c.Query("symbol")

	var orderbook map[string]any
	err := ctrl.cache.Get(c, fmt.Sprintf(matching.CacheKeyOrderbook, symbol), &orderbook)
	if err != nil {
		common.ResponseError(c, errors.New("orderbook not found"))
		return
	}
	common.ResponseOK(c, orderbook)
}

// @Summary trades
// @Description 获取近期成交记录
// @ID v1.market.trades
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/market/trades [get]
func (ctrl *MarketController) Trades(c *gin.Context) {
	common.ResponseOK(c, gin.H{})
}

// @Summary klines
// @Description 获取K线数据
// @ID v1.market.klines
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param period query string false "period" Enums(M1, M3, M5, M15, M30, H1, H2, H4, H6, H8, H12, D1, D3, W1, MN)
// @Param start query int false "start"
// @Param end query int false "end"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/market/klines [get]
func (ctrl *MarketController) Klines(c *gin.Context) {
	symbol := strings.ToLower(c.DefaultQuery("symbol", ""))
	period := strings.ToLower(c.DefaultQuery("period", "m1"))
	start := c.DefaultQuery("start", "0")
	end := c.DefaultQuery("end", "0")
	limit := c.DefaultQuery("limit", "1000")

	peroidType, err := kline_types.ParsePeriod(period)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	data, err := ctrl.klineRepo.Find(c, symbol, peroidType, startInt, endInt, limitInt)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, data)
}
