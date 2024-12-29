package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/duolacloud/crud-core/cache"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/common"
	"github.com/yzimhao/trading_engine/v2/internal/modules/matching"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MarketController struct {
	logger       *zap.Logger
	cache        cache.Cache
	klineRepo    persistence.KlineRepository
	tradeLogRepo persistence.TradeLogRepository
	tradeVariety persistence.TradeVarietyRepository
}

type inMarketContext struct {
	fx.In
	Logger       *zap.Logger
	Cache        cache.Cache
	KlineRepo    persistence.KlineRepository
	TradeLogRepo persistence.TradeLogRepository
	TradeVariety persistence.TradeVarietyRepository
}

func NewMarketController(in inMarketContext) *MarketController {
	return &MarketController{
		logger:       in.Logger,
		cache:        in.Cache,
		klineRepo:    in.KlineRepo,
		tradeLogRepo: in.TradeLogRepo,
		tradeVariety: in.TradeVariety,
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
	symbol := c.DefaultQuery("symbol", "")
	limit := c.DefaultQuery("limit", "1000")

	tradeVariety, err := ctrl.tradeVariety.FindBySymbol(c, symbol)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	data, err := ctrl.tradeLogRepo.Find(c, symbol, limitInt)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

	var response []map[string]any
	for _, v := range data {
		response = append(response, map[string]any{
			"id":       v.Id,
			"price":    common.FormatStrNumber(v.Price, tradeVariety.PriceDecimals),
			"qty":      common.FormatStrNumber(v.Quantity, tradeVariety.QtyDecimals),
			"amount":   common.FormatStrNumber(v.Amount, 6), //TODO 金额现实位数控制
			"trade_at": v.CreatedAt.UnixNano(),
		})
	}

	common.ResponseOK(c, response)
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

	tradeVariety, err := ctrl.tradeVariety.FindBySymbol(c, symbol)
	if err != nil {
		common.ResponseError(c, err)
		return
	}

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

	// [
	//     [
	//       1499040000000,      // k线开盘时间
	//       "0.01634790",       // 开盘价
	//       "0.80000000",       // 最高价
	//       "0.01575800",       // 最低价
	//       "0.01577100",       // 收盘价(当前K线未结束的即为最新价)
	//       "148976.11427815",  // 成交量
	//       1499644799999,      // k线收盘时间
	//       "2434.19055334",    // 成交额
	//       308,                // 成交笔数
	//       "1756.87402397",    // 主动买入成交量
	//       "28.46694368",      // 主动买入成交额
	//       "17928899.62484339" // 请忽略该参数
	//     ]
	//   ]

	response := make([][6]any, 0)
	for _, v := range data {
		response = append(response, [6]any{
			v.OpenAt.UnixMilli(),
			common.FormatStrNumber(v.Open, tradeVariety.PriceDecimals),
			common.FormatStrNumber(v.High, tradeVariety.PriceDecimals),
			common.FormatStrNumber(v.Low, tradeVariety.PriceDecimals),
			common.FormatStrNumber(v.Close, tradeVariety.PriceDecimals),
			common.FormatStrNumber(v.Volume, tradeVariety.QtyDecimals),
		})
	}

	common.ResponseOK(c, response)
}
