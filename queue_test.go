package trading_engine

import (
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestOrderInterface(t *testing.T) {
	// od := &Order{
	// 	orderId:    "1",
	// 	price:      decimal.NewFromFloat(1.1),
	// 	quantity:   decimal.NewFromFloat(1.103981),
	// 	createTime: 1,
	// }

	// assert.Equal(t, "1", od.GetOrderId())
	// assert.Equal(t, decimal.NewFromFloat(1.1), od.GetPrice())
	// assert.Equal(t, decimal.NewFromFloat(1.103981), od.GetQuantity())

}

func TestAskQueue(t *testing.T) {
	askQueue := NewQueue()
	askQueue.Push(NewAskItem("1", decimal.NewFromFloat(1.8), decimal.NewFromFloat(1), 11111111))
	askQueue.Push(NewAskItem("2", decimal.NewFromFloat(0.99), decimal.NewFromFloat(10), 11111111))
	askQueue.Push(NewAskItem("3", decimal.NewFromFloat(1.1), decimal.NewFromFloat(12), 11111111))

	assert.Equal(t, 3, askQueue.Len())

	//取出堆顶的一个元素
	top := askQueue.Top()
	assert.Equal(t, "2", top.GetUniqueId())

	//重新插入一个低价订单，重新获取堆顶的item
	askQueue.Push(NewAskItem("4", decimal.NewFromFloat(0.01), decimal.NewFromFloat(10), 11111111))
	top = askQueue.Top()
	assert.Equal(t, "4", top.GetUniqueId())

	//取出队列最后一个插入的元素
	// last := askQueue.Pop()
	// assert.Equal(t, "4", last.GetUniqueId())

	//从队列里移除一个指定的订单号
	remove := askQueue.Remove("4")
	assert.Equal(t, "4", remove.GetUniqueId())
	assert.Equal(t, 3, askQueue.Len())

}

func TestBidQueue(t *testing.T) {
	bidQueue := NewQueue()
	bidQueue.Push(NewBidItem("1", decimal.NewFromFloat(1.8), decimal.NewFromFloat(1), 11111111))
	bidQueue.Push(NewBidItem("2", decimal.NewFromFloat(1.1), decimal.NewFromFloat(1), 11111111))
	bidQueue.Push(NewBidItem("3", decimal.NewFromFloat(2), decimal.NewFromFloat(1), 11111111))

	assert.Equal(t, 3, bidQueue.Len())

	//取出堆顶的一个元素
	top := bidQueue.Top()
	assert.Equal(t, "3", top.GetUniqueId())

	//重新插入一个高价订单，重新获取堆顶的item
	bidQueue.Push(NewBidItem("4", decimal.NewFromFloat(10.01), decimal.NewFromFloat(10), 11111111))
	top = bidQueue.Top()
	assert.Equal(t, "4", top.GetUniqueId())

	//从队列里移除一个指定的订单号
	remove := bidQueue.Remove("3")
	assert.Equal(t, "3", remove.GetUniqueId())
	assert.Equal(t, 3, bidQueue.Len())
}

func BenchmarkAskQueue(b *testing.B) {
	askQueue := NewQueue()
	rand.Seed(time.Now().Unix())

	for i := 0; i < b.N; i++ {
		id := uuid.New().String()
		price := decimal.NewFromFloat(float64(rand.Intn(1000)) / 100)
		quantity := decimal.NewFromFloat(float64(rand.Intn(10000)) / 100)
		askQueue.Push(NewAskItem(id, price, quantity, time.Now().Unix()))
	}
}
