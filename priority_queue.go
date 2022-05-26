package trading_engine

import (
	"container/heap"
	"sync"
	"time"

	"github.com/shopspring/decimal"
)

type QueueItem interface {
	Less(item QueueItem) bool
	GetIndex() int
	SetIndex(index int)
	GetUniqueId() string
	GetPrice() decimal.Decimal
	GetQuantity() decimal.Decimal
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

func NewQueue() *OrderQueue {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	queue := OrderQueue{
		pq:    &pq,
		m:     make(map[string]*QueueItem),
		depth: make([][2]string, 0),
	}

	//flush depth
	go queue.flushDepth()
	return &queue
}

type OrderQueue struct {
	pq *PriorityQueue
	m  map[string]*QueueItem
	sync.Mutex

	depth [][2]string
}

func (o *OrderQueue) GetDepth() [][2]string {
	return o.depth
}

//刷新深度数据
func (o *OrderQueue) flushDepth() {

	sortMap := func(m map[string]string) [][2]string {
		var res [][2]string
		var keys []string
		for k, _ := range m {
			keys = append(keys, k)
		}

		keys = quickSort(keys)
		for _, k := range keys {
			res = append(res, [2]string{k, m[k]})
		}
		return res
	}

	for {
		o.Lock()

		o.depth = [][2]string{}
		depthMap := make(map[string]string)

		for i := 0; i < o.pq.Len(); i++ {
			item := (*o.pq)[i]
			price := item.GetPrice().String()

			qnt := item.GetQuantity()
			if _, ok := depthMap[price]; !ok {
				depthMap[price] = qnt.String()
			} else {
				old_qunantity, _ := decimal.NewFromString(depthMap[price])
				qnt = old_qunantity.Add(qnt)
				depthMap[price] = qnt.String()
			}
		}

		//按价格排序map
		o.depth = sortMap(depthMap)
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

func (o *OrderQueue) Remove(id string) QueueItem {
	o.Lock()
	defer o.Unlock()

	old, ok := o.m[id]
	if !ok {
		return nil
	}

	item := heap.Remove(o.pq, (*old).GetIndex())
	delete(o.m, id)
	return item.(QueueItem)
}
