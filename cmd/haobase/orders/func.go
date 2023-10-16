package orders

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
)

func generate_order_id_by_side(side trading_core.OrderSide) string {
	if side == trading_core.OrderSideSell {
		return generate_order_id("A")
	} else {
		return generate_order_id("B")
	}
}

func generate_order_id(prefix string) string {
	prefix = strings.ToUpper(prefix)
	s := time.Now().Format("060102150405")
	ns := time.Now().Nanosecond() / 1000
	rn := rand.Intn(99)
	return fmt.Sprintf("%s%s%06d%02d", prefix, s, ns, rn)
}

func push_new_order_to_redis(symbol string, data []byte) {
	topic := types.FormatNewOrder.Format(symbol)
	logrus.Infof("push %s new: %s", topic, data)
	ctx := context.Background()
	err := base.RDC().RPush(ctx, topic, data).Err()
	if err != nil {
		logrus.Errorf("push %s err: %s", topic, err.Error())
	}
}
