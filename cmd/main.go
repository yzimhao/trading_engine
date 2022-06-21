package main

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	te "github.com/yzimhao/trading_engine"
)

var (
	pair *te.TradePair
	rdb  *redis.Client
)

type Order struct {
	OrderId  string `json:"order_id"`
	Side     string `json:"side"` // buy、sell
	Price    string `json:"price"`
	Quantity string `json:"quantity"`
	Amount   string `json:"amount"`

	MaxHoldAmount string `json:"max_hold_amount"`
	MaxHoldQty    string `json:"max_hold_qty"`

	CreateTime string `json:"create_time"`
}

type NewOrderMsgBody struct {
	PriceType string `json:"price_type"` //limit、market-qty、market-amount
	Order     Order  `json:"order"`
}

type CancelOrderMsgBody struct {
	Side    string `json:"side"`
	OrderId string `json:"order_id"`
}

func main() {
	app := &cli.App{
		Name:  "trading_engine",
		Usage: "",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "config", Value: "config.yaml", Usage: "config file"},
			&cli.StringFlag{Name: "redis", Value: "localhost:6379", Usage: ""},
			&cli.StringFlag{Name: "redis-passwd", Value: "", Usage: ""},
			&cli.IntFlag{Name: "redis-db", Value: 0, Usage: ""},
			&cli.StringFlag{Name: "symbol", Value: "", Usage: "eg: btcusdt、ethusdt"},
			&cli.IntFlag{Name: "price-digit", Value: 2, Usage: "price digit"},
			&cli.IntFlag{Name: "qty-digit", Value: 0, Usage: "quantity digit"},
		},
		Action: func(c *cli.Context) error {
			rdb = redis.NewClient(&redis.Options{
				Addr:     c.String("redis"),
				Password: c.String("redis-passwd"),
				DB:       c.Int("redis-db"),
			})
			logrus.Info("start app")
			tradingEngineStart(c.String("symbol"), c.Int("price-digit"), c.Int("qty-digit"))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func tradingEngineStart(symbol string, pdig, qdig int) {
	pair = te.NewTradePair(strings.ToLower(symbol), pdig, qdig)

	wg := sync.WaitGroup{}
	//new order
	wg.Add(1)
	go newOrderFromRedis(&wg)

	//cancel order
	wg.Add(1)
	go cancelOrderFromRedis(&wg)

	//depth

	//publish msg
	wg.Add(1)
	go publishMsgToRedis(&wg)

	wg.Wait()
	logrus.Info("trading engine done.")
}

func str2decimal(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}

func str2Int64(a string) int64 {
	i, _ := strconv.ParseInt(a, 10, 64)
	return i
}
