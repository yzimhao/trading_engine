package haotrader

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/types"
)

func InsertAsk(rdc *redis.Client, symbol string, order_type string, n int, base_price string, qty string) {
	if n <= 0 {
		n = 1
	}

	digit := viper.GetInt(fmt.Sprintf("symbol.%s.price_digit", symbol))
	pdigit := fmt.Sprintf("%%.%df", digit)

	price, _ := strconv.ParseFloat(base_price, 64)
	for i := 0; i < n; i++ {

		id, _ := uuid.NewRandom()
		obj := Order{
			OrderId:   fmt.Sprintf("a-%s-%d", id.String(), i),
			Side:      "ask",
			OrderType: "limit",
			Price:     fmt.Sprintf(pdigit, price+rand.Float64()),
			Qty:       qty,
			At:        time.Now().Unix(),
		}
		s, _ := json.Marshal(obj)
		ctx1 := context.Background()
		key := types.FormatNewOrder.Format(symbol)
		rdc.RPush(ctx1, key, s).Err()
	}
}

func InsertBid(rdc *redis.Client, symbol string, order_type string, n int, base_price string, qty string) {
	if n <= 0 {
		n = 1
	}

	digit := viper.GetInt(fmt.Sprintf("symbol.%s.price_digit", symbol))
	pdigit := fmt.Sprintf("%%.%df", digit)

	price, _ := strconv.ParseFloat(base_price, 64)
	for i := 0; i < n; i++ {
		id, _ := uuid.NewRandom()
		obj := Order{
			OrderId:   fmt.Sprintf("b-%s-%d", id.String(), i),
			Side:      "bid",
			OrderType: "limit",
			Price:     fmt.Sprintf(pdigit, price+rand.Float64()),
			Qty:       qty,
			At:        time.Now().Unix(),
		}
		s, _ := json.Marshal(obj)
		ctx1 := context.Background()
		key := types.FormatNewOrder.Format(symbol)
		rdc.RPush(ctx1, key, s).Err()
	}
}
