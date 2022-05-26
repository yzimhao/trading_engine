package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine"
)

var askQueue *trading_engine.OrderQueue
var bidQueue *trading_engine.OrderQueue

func main() {
	g := gin.New()
	gin.SetMode(gin.DebugMode)
	g.LoadHTMLGlob("./*.html")
	askQueue = trading_engine.NewQueue()
	bidQueue = trading_engine.NewQueue()

	g.GET("/depth", func(c *gin.Context) {
		a := askQueue.GetDepth()
		b := bidQueue.GetDepth()

		c.JSON(200, gin.H{
			"ask": a,
			"bid": b,
		})
	})

	type args struct {
		OrderId   string `json:"order_id"`
		OrderType string `json:"order_type"`
		Price     string `json:"price"`
		Quantity  string `json:"quantity"`
	}

	g.POST("/new_order", func(c *gin.Context) {
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
	})

	g.GET("/demo", func(c *gin.Context) {
		c.HTML(200, "demo.html", nil)
	})

	//websocket

	g.Run(":8080")
}

func string2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}
