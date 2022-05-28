package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"example/wss"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine"
)

var askQueue *trading_engine.OrderQueue
var bidQueue *trading_engine.OrderQueue
var sendMsg chan []byte
var web *gin.Engine

func main() {

	gin.SetMode(gin.DebugMode)
	trading_engine.SetPriceDigits(4)
	trading_engine.SetQuantityDigits(0)

	startWeb()
}

func startWeb() {
	web = gin.New()
	web.LoadHTMLGlob("./*.html")
	askQueue = trading_engine.NewQueue()
	bidQueue = trading_engine.NewQueue()
	sendMsg = make(chan []byte, 100)
	trading_engine.MatchingEngine(askQueue, bidQueue)

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
	web.Run(":8080")
}

func depth(c *gin.Context) {
	a := askQueue.GetDepth()
	b := bidQueue.GetDepth()

	c.JSON(200, gin.H{
		"ask": a,
		"bid": b,
	})
}

func watchTradeLog() {
	for {
		select {
		case log := <-trading_engine.ChTradeResult:
			data := gin.H{
				"tag":  "trade",
				"data": log,
			}
			msg, _ := json.Marshal(data)
			sendMsg <- []byte(msg)
		}
	}
}

func pushDepth() {
	for {
		ask := askQueue.GetDepth()
		bid := bidQueue.GetDepth()
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
	param.Price = trading_engine.FormatPrice2Str(string2decimal(param.Price))
	param.Quantity = trading_engine.FormatQuantity2Str(string2decimal(param.Quantity))
	param.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	if strings.ToLower(param.OrderType) == "ask" {
		param.OrderId = fmt.Sprintf("a-%s", orderId)
		askOrder := trading_engine.NewAskItem(param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), time.Now().Unix())
		askQueue.Push(askOrder)
	} else {
		param.OrderId = fmt.Sprintf("b-%s", orderId)
		bidOrder := trading_engine.NewBidItem(param.OrderId, string2decimal(param.Price), string2decimal(param.Quantity), time.Now().Unix())
		bidQueue.Push(bidOrder)
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
			"ask_len": askQueue.Len(),
			"bid_len": bidQueue.Len(),
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

	askQueue.Remove(param.OrderId)
	bidQueue.Remove(param.OrderId)

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
