package matching

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

func (e *Engine) orderBook(queue *OrderQueue, size int) [][2]string {
	queue.mx.Lock()
	defer queue.mx.Unlock()

	max := len(queue.orderBook)
	if size <= 0 || size > max {
		size = max
	}

	book := make([][2]string, size)
	copy(book, queue.orderBook[0:size])
	return book
}

func (e *Engine) orderBookTicker(que *OrderQueue) {

	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)

	for {
		select {
		case <-e.ctx.Done():
			return
		case <-ticker.C:
			func() {
				e.mx.Lock()
				defer e.mx.Unlock()

				que.mx.Lock()
				defer que.mx.Unlock()

				que.orderBook = [][2]string{}

				bookMap := make(map[string]string)

				if que.pq.Len() > 0 {

					for i := 0; i < que.pq.Len(); i++ {

						if len(bookMap) > e.opts.orderBookMaxLen {
							break
						}

						item := (*que.pq)[i]

						price := item.GetPrice().StringFixedBank(e.opts.priceDecimals)
						if _, ok := bookMap[price]; !ok {
							bookMap[price] = item.GetQuantity().StringFixedBank(e.opts.quantityDecimals)
						} else {
							old_qunantity, err := decimal.NewFromString(bookMap[price])
							if err != nil {
								e.logger.Sugar().Errorf("[matching] orderbook parse error: %v value: %v", err, bookMap[price])
								continue
							}
							bookMap[price] = old_qunantity.Add(item.GetQuantity()).StringFixedBank(e.opts.quantityDecimals)
						}
					}

					//按价格排序map
					que.orderBook = sortMap2Slice(bookMap, que.Top().GetOrderSide())
				}
			}()
		}

	}
}

func sortMap2Slice(m map[string]string, ask_bid types.OrderSide) [][2]string {
	res := [][2]string{}
	keys := []string{}
	for k, _ := range m {
		keys = append(keys, k)
	}

	if ask_bid == types.OrderSideSell {
		keys = quickSort(keys, "asc")
	} else {
		keys = quickSort(keys, "desc")
	}

	for _, k := range keys {
		res = append(res, [2]string{k, m[k]})
	}
	return res
}

func quickSort(nums []string, asc_desc string) []string {
	if len(nums) <= 1 {
		return nums
	}

	spilt := nums[0]
	left := []string{}
	right := []string{}
	mid := []string{}

	for _, v := range nums {
		vv, _ := decimal.NewFromString(v)
		sp, _ := decimal.NewFromString(spilt)
		if vv.Cmp(sp) == -1 {
			left = append(left, v)
		} else if vv.Cmp(sp) == 1 {
			right = append(right, v)
		} else {
			mid = append(mid, v)
		}
	}

	left = quickSort(left, asc_desc)
	right = quickSort(right, asc_desc)

	if asc_desc == "asc" {
		return append(append(left, mid...), right...)
	} else {
		return append(append(right, mid...), left...)
	}
}
