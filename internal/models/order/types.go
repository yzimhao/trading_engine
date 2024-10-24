package order

import (
	"fmt"
	"strings"
	"time"

	"math/rand"

	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

type CreateOrder struct {
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

func (o *CreateOrder) GenerateOrderId() {
	if o.OrderSide == matching_types.OrderSideBuy {
		o.OrderId = generateOrderId("B")
	} else {
		o.OrderId = generateOrderId("A")
	}
}

func generateOrderId(prefix string) string {
	prefix = strings.ToUpper(prefix)
	t := time.Now().Format("060102150405")
	rn := rand.Intn(9999)
	return fmt.Sprintf("%s%s%04d", prefix, t, rn)
}
