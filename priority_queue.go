package trading_engine

import (
	"container/heap"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type QueueItem interface {
	SetIndex(index int)
	SetQuantity(quantity decimal.Decimal)
	SetAmount(amount decimal.Decimal)
	Less(item QueueItem) bool
	GetIndex() int
	GetUniqueId() string
	GetPrice() decimal.Decimal
	GetQuantity() decimal.Decimal
	GetCreateTime() int64
	GetOrderSide() OrderSide
	GetPriceType() PriceType
	GetAmount() decimal.Decimal //订单金额，在市价订单的时候生效，限价单不需要这个字段
}

type PriorityQueue []QueueItem

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Less(pq[j])
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].SetIndex(i)
	pq[j].SetIndex(j)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.SetIndex(-1)
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	x.(QueueItem).SetIndex(n)
	*pq = append(*pq, x.(QueueItem))
}

func NewQueue(priceDigit, quantityDigit int) *OrderQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	queue := OrderQueue{
		pq:            &pq,
		m:             make(map[string]*QueueItem),
		priceDigit:    priceDigit,
		quantityDigit: quantityDigit,
	}

	go queue.setDepth()
	return &queue
}

type OrderQueue struct {
	pq *PriorityQueue
	m  map[string]*QueueItem
	sync.Mutex

	priceDigit    int
	quantityDigit int
	depth         [][2]string
}

func (o *OrderQueue) GetDepth(limit int) [][2]string {
	o.Lock()
	defer o.Unlock()

	max := len(o.depth)
	if limit <= 0 || limit > max {
		limit = max
	}

	dp := o.depth[0:limit]
	return dp
}

func (o *OrderQueue) setDepth() {

	sortMap := func(m map[string]string, ask_bid OrderSide) [][2]string {
		res := [][2]string{}
		keys := []string{}
		for k, _ := range m {
			keys = append(keys, k)
		}

		if ask_bid == OrderSideSell {
			keys = quickSort(keys, "asc")
		} else {
			keys = quickSort(keys, "desc")
		}

		for _, k := range keys {
			res = append(res, [2]string{k, m[k]})
		}
		return res
	}

	for {
		o.Lock()

		o.depth = [][2]string{}
		depthMap := make(map[string]string)

		if o.pq.Len() > 0 {

			for i := 0; i < o.pq.Len(); i++ {
				item := (*o.pq)[i]

				price := FormatDecimal2String(item.GetPrice(), o.priceDigit)

				if _, ok := depthMap[price]; !ok {
					depthMap[price] = FormatDecimal2String(item.GetQuantity(), o.quantityDigit)
				} else {
					old_qunantity, _ := decimal.NewFromString(depthMap[price])
					depthMap[price] = FormatDecimal2String(old_qunantity.Add(item.GetQuantity()), o.quantityDigit)
				}
			}

			//按价格排序map
			o.depth = sortMap(depthMap, o.Top().GetOrderSide())
		}
		o.Unlock()
		time.Sleep(time.Millisecond * 20)
	}
}

func (o *OrderQueue) Len() int {
	return o.pq.Len()
}

func (o *OrderQueue) Push(item QueueItem) (exist bool) {
	o.Lock()
	defer o.Unlock()

	//todo 触发撮合订单

	id := item.GetUniqueId()
	if _, ok := o.m[id]; ok {
		return true
	}

	heap.Push(o.pq, item)
	o.m[id] = &item
	return false
}

// func (o *OrderQueue) Pop() QueueItem {
// 	item := heap.Pop(o.pq)
// 	id := item.(QueueItem).GetUniqueId()
// 	delete(o.m, id)
// 	return item.(QueueItem)
// }

func (o *OrderQueue) Get(index int) QueueItem {
	n := o.pq.Len()
	if n <= index {
		return nil
	}

	return (*o.pq)[index]
}

func (o *OrderQueue) Top() QueueItem {
	return o.Get(0)
}

func (o *OrderQueue) Remove(uniqId string) QueueItem {
	o.Lock()
	defer o.Unlock()

	old, ok := o.m[uniqId]
	if !ok {
		return nil
	}

	item := heap.Remove(o.pq, (*old).GetIndex())
	delete(o.m, uniqId)
	return item.(QueueItem)
}
