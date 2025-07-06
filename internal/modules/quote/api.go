package quote

import (
	"fmt"
	"strings"

	"github.com/duolacloud/crud-core/cache"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/matching"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	"go.uber.org/zap"
)

type quoteApi struct {
	router *provider.Router
	logger *zap.Logger
	cache  cache.Cache
}

func newQuoteApi(
	router *provider.Router,
	logger *zap.Logger,
	cache cache.Cache,
) {
	q := quoteApi{
		router: router,
		logger: logger,
		cache:  cache,
	}
	q.registerRouter()
}

func (q *quoteApi) registerRouter() {
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

// TODO
func (q *quoteApi) depth(c *gin.Context) {
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
func (q *quoteApi) trades(c *gin.Context) {}

// TODO
func (q *quoteApi) historicalTrades(c *gin.Context) {}

// TODO
func (q *quoteApi) klines(c *gin.Context) {}

// TODO
func (q *quoteApi) avgPrice(c *gin.Context) {}

// TODO
func (q *quoteApi) ticker24hr(c *gin.Context) {}

// TODO
func (q *quoteApi) tickerPrice(c *gin.Context) {}

// TODO
func (q *quoteApi) tickerBookTicker(c *gin.Context) {}
