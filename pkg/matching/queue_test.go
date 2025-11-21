package matching_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/yzimhao/trading_engine/v2/pkg/matching"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

var askQueue *matching.OrderQueue
var bidQueue *matching.OrderQueue

func init() {
	askQueue = matching.NewQueue()
	bidQueue = matching.NewQueue()
}

func d(v float64) decimal.Decimal {
	return decimal.NewFromFloat(v)
}

func TestMarketSubOrderType(t *testing.T) {
	Convey("市价订单", t, func() {
		Convey("按金额卖出", func() {
			item := matching.NewAskMarketAmountItem("1", d(1), d(100), 1)
			So(item.GetSubOrderType(), ShouldEqual, types.SubOrderTypeMarketByAmount)
		})

		Convey("按数量卖出", func() {
			item := matching.NewAskMarketQtyItem("1", d(2), 2)
			So(item.GetSubOrderType(), ShouldEqual, types.SubOrderTypeMarketByQty)
		})
	})
}

func TestAskQueue(t *testing.T) {

	Convey("卖盘挂单队列", t, func() {

		Convey("队列元素的顺序", func() {
			askQueue.Push(matching.NewAskLimitItem("1", d(1.8), d(1), 11111111))
			askQueue.Push(matching.NewAskLimitItem("2", d(0.99), d(10), 11111111))
			askQueue.Push(matching.NewAskLimitItem("3", d(1.1), d(12), 11111111))
			askQueue.Push(matching.NewAskLimitItem("5", d(1.1), d(12), 11111110))
			So(askQueue.Len(), ShouldEqual, 4)

			top := askQueue.Top()
			So(top.GetUniqueId(), ShouldEqual, "2")
			So(top.GetPrice(), ShouldEqual, d(0.99))
			So(top.GetQuantity(), ShouldEqual, d(10))

		})

		Convey("重新插入低价卖单，重新获取队列顶元素", func() {
			askQueue.Push(matching.NewAskLimitItem("4", d(0.01), d(10), 11111111))
			top := askQueue.Top()
			So(top.GetUniqueId(), ShouldEqual, "4")
		})

		Convey("更新队列顶元素", func() {
			top := askQueue.Top()
			top.SetQuantity(d(10.01))
			So(top.GetQuantity(), ShouldEqual, d(10.01))
		})

		Convey("移除队列一个指定的订单号", func() {
			So(askQueue.Len(), ShouldEqual, 5)
			remove := askQueue.Remove("4")
			So(remove.GetUniqueId(), ShouldEqual, "4")
			So(askQueue.Len(), ShouldEqual, 4)
		})
	})

}

func TestBidQueue(t *testing.T) {

	Convey("买盘挂单队列", t, func() {

		Convey("队列元素的顺序", func() {
			bidQueue.Push(matching.NewBidLimitItem("1", d(1.8), d(1), 11111111))
			bidQueue.Push(matching.NewBidLimitItem("2", d(0.99), d(10), 11111111))
			bidQueue.Push(matching.NewBidLimitItem("3", d(1.1), d(12), 11111111))
			bidQueue.Push(matching.NewBidLimitItem("5", d(1.1), d(12), 11111110))
			So(bidQueue.Len(), ShouldEqual, 4)

			top := bidQueue.Top()
			So(top.GetUniqueId(), ShouldEqual, "1")
			So(top.GetPrice(), ShouldEqual, d(1.8))
			So(top.GetQuantity(), ShouldEqual, d(1))

		})

		Convey("重新插入高价买单，重新获取队列顶元素", func() {
			bidQueue.Push(matching.NewBidLimitItem("4", d(2), d(10), 11111111))
			top := bidQueue.Top()
			So(top.GetUniqueId(), ShouldEqual, "4")
		})

		Convey("更新队列顶元素", func() {
			top := bidQueue.Top()
			top.SetQuantity(d(10.01))
			So(top.GetQuantity(), ShouldEqual, d(10.01))
		})

		Convey("移除队列一个指定的订单号", func() {
			So(bidQueue.Len(), ShouldEqual, 5)
			remove := bidQueue.Remove("4")
			So(remove.GetUniqueId(), ShouldEqual, "4")
			So(bidQueue.Len(), ShouldEqual, 4)
		})
	})

}

func BenchmarkAskQueue(b *testing.B) {
	askQueue := matching.NewQueue()

	for i := 0; i < b.N; i++ {
		id := uuid.New().String()
		price := decimal.NewFromFloat(float64(rand.Intn(1000)) / 100)
		quantity := decimal.NewFromFloat(float64(rand.Intn(10000)) / 100)
		askQueue.Push(matching.NewAskLimitItem(id, price, quantity, time.Now().Unix()))
	}
}
