package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/redis/go-redis/v9"
	"github.com/sevlyar/go-daemon"

	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/haotrader"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
)

var (
	web *gin.Engine
	rdc *redis.Client
)

func initredis(r string) {
	rdc = redis.NewClient(&redis.Options{
		Addr: r,
	})
}

func main() {
	port := flag.String("port", "8080", "port")
	rd := flag.String("redis", "127.0.0.1:6379", "redis")
	d := flag.Bool("d", false, "deamon")
	flag.Parse()

	initredis(*rd)

	if *d {
		cntxt := &daemon.Context{
			PidFileName: "run.pid",
			PidFilePerm: 0644,
			LogFileName: "run.log",
			LogFilePerm: 0640,
			WorkDir:     "./",
			Umask:       027,
			// Args:        ,
		}

		d, err := cntxt.Reborn()
		if err != nil {
			log.Fatal("Unable to run: ", err)
		}
		if d != nil {
			return
		}
		defer cntxt.Release()
	}

	gin.SetMode(gin.DebugMode)
	startWeb(*port)
}

func startWeb(port string) {
	web = gin.New()
	web.LoadHTMLGlob("./*.html")
	web.StaticFS("/statics", http.Dir("./statics"))

	web.POST("/api/new_order", newOrder)
	web.POST("/api/cancel_order", cancelOrder)
	web.GET("/api/test_rand", testOrder)

	web.GET("/:symbol", func(c *gin.Context) {
		support := []string{"usdjpy", "eurusd"}
		symbol := strings.ToLower(c.Param("symbol"))

		if !arrutil.Contains(support, symbol) {
			c.Redirect(301, "/")
			return
		}

		c.HTML(200, "demo.html", gin.H{
			"symbol": symbol,
		})
	})

	web.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/usdjpy")
	})

	web.Run(":" + port)
}

func newOrder(c *gin.Context) {
	type args struct {
		OrderId   string `json:"order_id"`
		OrderType string `json:"order_type"`
		PriceType string `json:"price_type"`
		Price     string `json:"price"`
		Quantity  string `json:"quantity"`
		Amount    string `json:"amount"`
		MaxQty    string `json:"max_qty"`
		MaxAmount string `json:"max_amount"`
		Symbol    string `json:"symbol"`
	}

	var param args
	c.BindJSON(&param)

	orderId := uuid.NewString()
	param.OrderId = orderId

	amount := string2decimal(param.Amount)
	price := string2decimal(param.Price)
	quantity := string2decimal(param.Quantity)

	var pt trading_core.OrderType
	if param.PriceType == "market" {
		param.Price = "0"
		pt = trading_core.OrderTypeMarket
		if param.Amount != "" {
			pt = trading_core.OrderTypeMarketAmount
			//市价按成交金额卖出时，默认持有该资产1000个
			param.MaxQty = "10000"
			if amount.Cmp(decimal.NewFromFloat(100000000)) > 0 || amount.Cmp(decimal.Zero) <= 0 {
				c.JSON(200, gin.H{
					"ok":    false,
					"error": "金额必须大于0，且不能超过 100000000",
				})
				return
			}

		} else if param.Quantity != "" {
			pt = trading_core.OrderTypeMarketQuantity
			//市价按数量买入资产时，需要用户账户所有可用资产数量，测试默认100块
			param.MaxAmount = "10000"
			if quantity.Cmp(decimal.NewFromFloat(100000000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
				c.JSON(200, gin.H{
					"ok":    false,
					"error": "数量必须大于0，且不能超过 100000000",
				})
				return
			}
		}
	} else {
		pt = trading_core.OrderTypeLimit

		param.Amount = "0"
		param.MaxAmount = "0"
		param.MaxQty = "0"

		if price.Cmp(decimal.NewFromFloat(2000)) > 0 || price.Cmp(decimal.Zero) < 0 {
			c.JSON(200, gin.H{
				"ok":    false,
				"error": "价格必须大于等于0，且不能超过 2000",
			})
			return
		}
		if quantity.Cmp(decimal.NewFromFloat(1000)) > 0 || quantity.Cmp(decimal.Zero) <= 0 {
			c.JSON(200, gin.H{
				"ok":    false,
				"error": "数量必须大于0，且不能超过 1000",
			})
			return
		}
	}

	data := haotrader.Order{
		OrderId:   param.OrderId,
		OrderType: pt.String(),
		Side:      strings.ToLower(param.OrderType),
		Price:     param.Price,
		Qty:       param.Quantity,
		MaxQty:    param.MaxQty,
		Amount:    param.Amount,
		MaxAmount: param.MaxAmount,
		At:        time.Now().UnixNano(),
	}
	if data.Side == "ask" {
		data.OrderId = fmt.Sprintf("a-%s", orderId)
	} else {
		data.OrderId = fmt.Sprintf("b-%s", orderId)
	}

	push_redis(param.Symbol, data.Json())
	c.JSON(200, gin.H{
		"ok":   true,
		"data": gin.H{
			// "ask_len": btcusdt.AskLen(),
			// "bid_len": btcusdt.BidLen(),
		},
	})
}

func push_redis(symbol string, data []byte) {
	ctx := context.Background()
	rdc.RPush(ctx, types.FormatNewOrder.Format(symbol), data)
}

func testOrder(c *gin.Context) {
	symbol := c.Query("symbol")
	op := strings.ToLower(c.Query("op_type"))
	latest := c.Query("latest_price")
	if op != "ask" {
		op = "bid"
	}

	latest_price := string2decimal(latest)
	if latest_price.Cmp(string2decimal("50")) <= 0 {
		latest_price = string2decimal("50")
	}
	d2 := string2decimal("30")

	func() {
		cnt := 10
		for i := 0; i < cnt; i++ {
			orderId := uuid.NewString()
			price := ""
			if op == "ask" {
				orderId = fmt.Sprintf("a-%s", orderId)
				price = randDecimal(latest_price.IntPart(), latest_price.Add(d2).IntPart()).String()
			} else {
				orderId = fmt.Sprintf("b-%s", orderId)
				price = randDecimal(latest_price.Sub(d2).IntPart(), latest_price.IntPart()).String()
			}

			data := haotrader.Order{
				OrderId:   orderId,
				OrderType: "limit",
				Side:      op,
				Price:     price,
				Qty:       "1",
				At:        time.Now().UnixNano(),
			}
			push_redis(symbol, data.Json())
		}
	}()

	c.JSON(200, gin.H{
		"ok":   true,
		"data": gin.H{},
	})
}

func cancelOrder(c *gin.Context) {
	type args struct {
		OrderId string `json:"order_id"`
	}

	var param args
	c.BindJSON(&param)

	if param.OrderId == "" {
		c.Abort()
		return
	}
	if strings.HasPrefix(param.OrderId, "a-") {
		// btcusdt.CancelOrder(trading_engine.OrderSideSell, param.OrderId)
	} else {
		// btcusdt.CancelOrder(trading_engine.OrderSideBuy, param.OrderId)
	}

	c.JSON(200, gin.H{
		"ok": true,
	})
}

func string2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}

func randDecimal(min, max int64) decimal.Decimal {
	d := decimal.New(rand.Int63n(max-min)+min, 0)
	rnd := rand.Float64()
	return d.Add(decimal.NewFromFloat(rnd))
}
