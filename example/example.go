package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"example/wss"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine"
)

var sendMsg chan []byte
var web *gin.Engine
var btcusdt *trading_engine.TradePair

func main() {

	port := flag.String("port", "8080", "port")
	flag.Parse()
	gin.SetMode(gin.DebugMode)

	btcusdt = trading_engine.NewTradePair("BTC_USDT", 2, 6)

	startWeb(*port)
}

func startWeb(port string) {
	web = gin.New()
	web.LoadHTMLGlob("./*.html")

	sendMsg = make(chan []byte, 100)

	go pushDepth()
	go watchTradeLog()

	web.GET("/api/depth", depth)
	web.POST("/api/new_order", newOrder)
	web.POST("/api/cancel_order", cancelOrder)

	web.GET("/demo", func(c *gin.Context) {
		c.HTML(200, "demo.html", nil)
	})

	//websocket
	{
		wss.HHub = wss.NewHub()
		go wss.HHub.Run()
		go func() {
			for {
				select {
				case data := <-sendMsg:
					wss.HHub.Send(data)
				default:
					time.Sleep(time.Duration(100) * time.Millisecond)
				}
			}
		}()

		web.GET("/ws", wss.ServeWs)
		web.GET("/pong", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	web.Run(":" + port)
}

func depth(c *gin.Context) {
	a := btcusdt.GetAskDepth(10)
	b := btcusdt.GetBidDepth(10)

	c.JSON(200, gin.H{
		"ask": a,
		"bid": b,
	})
}

func watchTradeLog() {
	for {
		if log, ok := <-btcusdt.ChTradeResult; ok {
			data := gin.H{
				"tag": "trade",
				"data": gin.H{
					"TradePrice":    trading_engine.FormatDecimal2String(log.TradePrice, btcusdt.PriceDigit),
					"TradeAmount":   trading_engine.FormatDecimal2String(log.TradeAmount, btcusdt.PriceDigit),
					"TradeQuantity": trading_engine.FormatDecimal2String(log.TradeQuantity, btcusdt.QuantityDigit),
					"TradeTime":     log.TradeTime,
					"AskOrderId":    log.AskOrderId,
					"BidOrderId":    log.BidOrderId,
				},
			}
			msg, _ := json.Marshal(data)
			sendMsg <- []byte(msg)
		}
	}
}

func pushDepth() {
	for {
		ask := btcusdt.GetAskDepth(10)
		bid := btcusdt.GetBidDepth(10)
		data := gin.H{
			"tag": "depth",
			"data": gin.H{
				"ask": ask,
				"bid": bid,
			},
		}
		msg, _ := json.Marshal(data)
		sendMsg <- []byte(msg)
		time.Sleep(time.Duration(500) * time.Millisecond)

	}
}

func newOrder(c *gin.Context) {
	type args struct {
		OrderId    string `json:"order_id"`
		OrderType  string `json:"order_type"`
		PriceType  string `json:"price_type"`
		Price      string `json:"price"`
		Quantity   string `json:"quantity"`
		CreateTime string `json:"create_time"`
	}

	var param args
	c.BindJSON(&param)

	if param.Price == "" || param.Quantity == "" {
		c.Abort()
		return
	}

	orderId := uuid.NewString()
	param.OrderId = orderId
	param.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	if strings.ToLower(param.OrderType) == "ask" {
		param.OrderId = fmt.Sprintf("a-%s", orderId)
		item := trading_engine.NewAskItem(param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), time.Now().Unix())
		btcusdt.PushNewOrder(trading_engine.OrderSideSell, item)

	} else {
		param.OrderId = fmt.Sprintf("b-%s", orderId)
		item := trading_engine.NewBidItem(param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), time.Now().Unix())
		btcusdt.PushNewOrder(trading_engine.OrderSideBuy, item)
	}

	go func() {
		msg := gin.H{
			"tag":  "new_order",
			"data": param,
		}
		msgByte, _ := json.Marshal(msg)
		sendMsg <- []byte(msgByte)
	}()

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"ask_len": btcusdt.AskLen(),
			"bid_len": btcusdt.BidLen(),
		},
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
		btcusdt.CancelOrder(trading_engine.OrderSideSell, param.OrderId)
	} else {
		btcusdt.CancelOrder(trading_engine.OrderSideBuy, param.OrderId)
	}

	go func() {
		msg := gin.H{
			"tag":  "cancel_order",
			"data": param,
		}
		msgByte, _ := json.Marshal(msg)
		sendMsg <- []byte(msgByte)
	}()

	c.JSON(200, gin.H{
		"ok": true,
	})
}

func string2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}
