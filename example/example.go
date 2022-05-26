package main

import (
	"encoding/json"
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

	startWeb()
}

func startWeb() {
	web = gin.New()
	web.LoadHTMLGlob("./*.html")
	askQueue = trading_engine.NewQueue()
	bidQueue = trading_engine.NewQueue()
	sendMsg = make(chan []byte, 100)

	go pushDepth()

	web.GET("/api/depth", depth)
	web.POST("/api/new_order", newOrder)

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
		OrderId   string `json:"order_id"`
		OrderType string `json:"order_type"`
		Price     string `json:"price"`
		Quantity  string `json:"quantity"`
	}

	var param args
	c.BindJSON(&param)

	if param.Price == "" || param.Quantity == "" {
		c.Abort()
		return
	}

	// rand.Seed(time.Now().Unix())
	// rand_price := rand.Float64()

	if strings.ToLower(param.OrderType) == "ask" {
		askOrder := trading_engine.NewAskItem(uuid.NewString(), string2decimal(param.Price), string2decimal(param.Quantity), time.Now().Unix())
		askQueue.Push(askOrder)
	} else {
		bidOrder := trading_engine.NewBidItem(uuid.NewString(), string2decimal(param.Price), string2decimal(param.Quantity), time.Now().Unix())
		bidQueue.Push(bidOrder)
	}

	c.JSON(200, gin.H{
		"ok": true,
		"data": gin.H{
			"ask_len": askQueue.Len(),
			"bid_len": bidQueue.Len(),
		},
	})
}

func string2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}
