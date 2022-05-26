package trading_engine

import (
	"container/heap"

	"github.com/shopspring/decimal"
)

type QueueItem interface {
	Less(item QueueItem) bool
	GetIndex() int
	SetIndex(index int)
	GetUniqueId() string
	GetPrice() decimal.Decimal
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
	return &OrderQueue{
		pq: &pq,
		m:  make(map[string]*QueueItem),
	}
}

type OrderQueue struct {
	pq *PriorityQueue
	m  map[string]*QueueItem
}

func (o *OrderQueue) Len() int {
	return o.pq.Len()
}

func (o *OrderQueue) Push(item QueueItem) (exist bool) {
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
	old, ok := o.m[id]
	if !ok {
		return nil
	}

	item := heap.Remove(o.pq, (*old).GetIndex())
	delete(o.m, id)
	return item.(QueueItem)
}
