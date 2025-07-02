package order

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

func generateOrderId(side matching_types.OrderSide) string {
	if side == matching_types.OrderSideBuy {
		return generateId("B")
	} else {
		return generateId("A")
	}
}

func generateId(prefix string) string {
	prefix = strings.ToUpper(prefix)
	t := time.Now().Format("060102150405")
	rn := rand.Intn(9999)
	return fmt.Sprintf("%s%s%04d", prefix, t, rn)
}
