package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
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

	trading_engine.Debug = false
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
	web.GET("/api/test_rand", testOrder)

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

func sendMessage(tag string, data interface{}) {
	msg := gin.H{
		"tag":  tag,
		"data": data,
	}
	msgByte, _ := json.Marshal(msg)
	sendMsg <- []byte(msgByte)
}

func watchTradeLog() {
	for {
		select {
		case log, ok := <-btcusdt.ChTradeResult:
			if ok {

				sendMessage("trade", gin.H{
					"TradePrice":    btcusdt.Price2String(log.TradePrice),
					"TradeAmount":   btcusdt.Price2String(log.TradeAmount),
					"TradeQuantity": btcusdt.Qty2String(log.TradeQuantity),
					"TradeTime":     log.TradeTime,
					"AskOrderId":    log.AskOrderId,
					"BidOrderId":    log.BidOrderId,
				})

				//latest price
				sendMessage("latest_price", gin.H{
					"latest_price": btcusdt.Price2String(log.TradePrice),
				})

			}
		case cancelOrderId := <-btcusdt.ChCancelResult:
			sendMessage("cancel_order", gin.H{
				"OrderId": cancelOrderId,
			})
		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
		}

	}
}

func pushDepth() {
	for {
		ask := btcusdt.GetAskDepth(10)
		bid := btcusdt.GetBidDepth(10)

		sendMessage("depth", gin.H{
			"ask": ask,
			"bid": bid,
		})

		time.Sleep(time.Duration(150) * time.Millisecond)
	}
}

func newOrder(c *gin.Context) {
	type args struct {
		OrderId    string `json:"order_id"`
		OrderType  string `json:"order_type"`
		PriceType  string `json:"price_type"`
		Price      string `json:"price"`
		Quantity   string `json:"quantity"`
		Amount     string `json:"amount"`
		CreateTime string `json:"create_time"`
	}

	var param args
	c.BindJSON(&param)

	orderId := uuid.NewString()
	param.OrderId = orderId
	param.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	var pt trading_engine.PriceType
	if param.PriceType == "market" {
		param.Price = "0"
		pt = trading_engine.PriceTypeMarket
		if param.Amount != "" {
			pt = trading_engine.PriceTypeMarketAmount
		} else if param.Quantity != "" {
			pt = trading_engine.PriceTypeMarketQuantity
		}
	} else {
		pt = trading_engine.PriceTypeLimit
		param.Amount = "0"
	}

	if strings.ToLower(param.OrderType) == "ask" {
		param.OrderId = fmt.Sprintf("a-%s", orderId)
		item := trading_engine.NewAskItem(pt, param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), string2decimal(param.Amount), time.Now().Unix())
		btcusdt.ChNewOrder <- item

	} else {
		param.OrderId = fmt.Sprintf("b-%s", orderId)
		item := trading_engine.NewBidItem(pt, param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), string2decimal(param.Amount), time.Now().Unix())
		btcusdt.ChNewOrder <- item
	}

	go sendMessage("new_order", param)

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"ask_len": btcusdt.AskLen(),
			"bid_len": btcusdt.BidLen(),
		},
	})
}

func testOrder(c *gin.Context) {
	op := strings.ToLower(c.Query("op_type"))
	if op != "ask" {
		op = "bid"
	}

	func() {
		cnt := 10
		for i := 0; i < cnt; i++ {
			orderId := uuid.NewString()
			if op == "ask" {
				orderId = fmt.Sprintf("a-%s", orderId)
				item := trading_engine.NewAskLimitItem(orderId, randDecimal(11, 20), randDecimal(20, 100), time.Now().Unix())
				btcusdt.ChNewOrder <- item
			} else {
				orderId = fmt.Sprintf("b-%s", orderId)
				item := trading_engine.NewBidLimitItem(orderId, randDecimal(1, 10), randDecimal(20, 100), time.Now().Unix())
				btcusdt.ChNewOrder <- item
			}

		}
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

	go sendMessage("cancel_order", param)

	c.JSON(200, gin.H{
		"ok": true,
	})
}

func string2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}

func randDecimal(min, max int64) decimal.Decimal {
	rand.Seed(time.Now().UnixNano())

	d := decimal.New(rand.Int63n(max-min)+min, 0)
	return d
}
