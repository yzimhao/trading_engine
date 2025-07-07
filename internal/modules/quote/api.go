package quote

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/duolacloud/crud-core/cache"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/matching"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"

	"go.uber.org/zap"
)

type QuoteApi struct {
	router      *provider.Router
	logger      *zap.Logger
	cache       cache.Cache
	klineRepo   persistence.KlineRepository
	productRepo persistence.ProductRepository
}

func newQuoteApi(
	router *provider.Router,
	logger *zap.Logger,
	cache cache.Cache,
	kline persistence.KlineRepository,
	product persistence.ProductRepository,
) *QuoteApi {
	q := QuoteApi{
		router:      router,
		logger:      logger,
		cache:       cache,
		klineRepo:   kline,
		productRepo: product,
	}
	return &q
}

func (q *QuoteApi) Run() {
	q.registerRouter()
}

func (q *QuoteApi) registerRouter() {
	//深度信息
	q.router.APIv1.GET("depth", q.depth)
	//近期成交
	q.router.APIv1.GET("trades", q.trades)
	//查询历史成交
	q.router.APIv1.GET("historicalTrades", q.historicalTrades)
	//k线
	q.router.APIv1.GET("klines", q.klines)
	//当前平均价格
	q.router.APIv1.GET("avgPrice", q.avgPrice)
	//24hr价格变动情况
	q.router.APIv1.GET("ticker/24hr", q.ticker24hr)
	//最新价格接口
	q.router.APIv1.GET("ticker/price", q.tickerPrice)
	//最优挂单接口
	q.router.APIv1.GET("ticker/bookTicker", q.tickerBookTicker)
}

// @Summary depth
// @Description get depth
// @ID v1.depth
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/depth [get]
func (q *QuoteApi) depth(c *gin.Context) {
	symbol := strings.ToLower(c.Query("symbol"))
	var orderbook map[string]any
	err := q.cache.Get(c, fmt.Sprintf(matching.CacheKeyOrderbook, symbol), &orderbook)
	if err != nil {
		q.logger.Sugar().Errorf("depth: ", err)
		q.router.ResponseError(c, types.ErrInternalError)
		return
	}
	q.router.ResponseOk(c, orderbook)
}

// TODO
func (q *QuoteApi) trades(c *gin.Context) {}

// TODO
func (q *QuoteApi) historicalTrades(c *gin.Context) {}

// @Summary K线数据
// @Description 获取K线数据
// @ID v1.klines
// @Tags market
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param period query string false "period" Enums(M1, M3, M5, M15, M30, H1, H2, H4, H6, H8, H12, D1, D3, W1, MN)
// @Param start query int false "start"
// @Param end query int false "end"
// @Param limit query int false "limit"
// @Success 200 {string} any
// @Router /api/v1/klines [get]
func (q *QuoteApi) klines(c *gin.Context) {
	symbol := strings.ToLower(c.DefaultQuery("symbol", ""))
	period := strings.ToLower(c.DefaultQuery("period", "m1"))
	start := c.DefaultQuery("start", "0")
	end := c.DefaultQuery("end", "0")
	limit := c.DefaultQuery("limit", "1000")

	product, err := q.productRepo.Get(symbol)
	if err != nil {
		q.logger.Sugar().Errorf("v1.klines err: %s", err)
		q.router.ResponseError(c, types.ErrInternalError)
		return
	}

	peroidType, err := kline_types.ParsePeriod(period)
	if err != nil {
		q.logger.Sugar().Errorf("v1.klines err: %s", err)
		q.router.ResponseError(c, types.ErrInternalError)
		return
	}

	startInt, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		q.logger.Sugar().Errorf("v1.klines err: %s", err)
		q.router.ResponseError(c, types.ErrInternalError)
		return
	}

	endInt, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		q.logger.Sugar().Errorf("v1.klines err: %s", err)
		q.router.ResponseError(c, types.ErrInternalError)
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		q.logger.Sugar().Errorf("v1.klines err: %s", err)
		q.router.ResponseError(c, types.ErrInternalError)
		return
	}

	data, err := q.klineRepo.Find(c, symbol, peroidType, startInt, endInt, limitInt)
	if err != nil {
		q.logger.Sugar().Errorf("v1.klines err: %s", err)
		q.router.ResponseError(c, types.ErrInternalError)
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
			v.Open.Truncate(product.PriceDecimals).String(),
			v.High.Truncate(product.PriceDecimals).String(),
			v.Low.Truncate(product.PriceDecimals).String(),
			v.Close.Truncate(product.PriceDecimals).String(),
			v.Volume.Truncate(product.QtyDecimals).String(),
		})
	}

	q.router.ResponseOk(c, response)
}

// TODO
func (q *QuoteApi) avgPrice(c *gin.Context) {}

// TODO
func (q *QuoteApi) ticker24hr(c *gin.Context) {}

// TODO
func (q *QuoteApi) tickerPrice(c *gin.Context) {}

// TODO
func (q *QuoteApi) tickerBookTicker(c *gin.Context) {}
