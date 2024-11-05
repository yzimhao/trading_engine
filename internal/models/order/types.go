package order

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"golang.org/x/exp/rand"
)

type Order struct {
	Symbol       string                   `json:"symbol"`
	OrderId      string                   `json:"order_id"`
	OrderSide    matching_types.OrderSide `json:"order_side"`
	OrderType    matching_types.OrderType `json:"order_type"` //价格策略，市价单，限价单
	UserId       string                   `json:"user_id"`
	Price        *string                  `json:"price"`
	Quantity     string                   `json:"quantity"`
	FeeRate      string                   `json:"fee_rate"`
	Amount       *string                  `json:"amount"`
	FreezeQty    string                   `json:"freeze_qty"`
	FreezeAmount string                   `json:"freeze_amount"`
	Status       models_types.OrderStatus `json:"status"`
	NanoTime     int64                    `json:"nano_time"`
}

func GenerateOrderId(side matching_types.OrderSide) string {
	if side == matching_types.OrderSideBuy {
		return generateOrderId("B")
	} else {
		return generateOrderId("A")
	}
}

func GenerateTradeId(ask, bid string) string {
	date := time.Now().Format("060102")
	raw := fmt.Sprintf("%s%s", ask, bid)

	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", raw)))
	hashed := fmt.Sprintf("%x", hash.Sum(nil))
	return fmt.Sprintf("T%s%s", date, hashed[0:17])
}

func generateOrderId(prefix string) string {
	prefix = strings.ToUpper(prefix)
	t := time.Now().Format("060102150405")
	rn := rand.Intn(9999)
	return fmt.Sprintf("%s%s%04d", prefix, t, rn)
}
