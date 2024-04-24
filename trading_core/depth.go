package trading_core

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/utils"
)

func (t *TradePair) GetAskDepth(size int) [][2]string {
	return t.depth(t.askQueue, size)
}

func (t *TradePair) GetBidDepth(size int) [][2]string {
	return t.depth(t.bidQueue, size)
}

func (t *TradePair) depth(queue *OrderQueue, size int) [][2]string {
	queue.Lock()
	defer queue.Unlock()

	max := len(queue.depth)
	if size <= 0 || size > max {
		size = max
	}

	return queue.depth[0:size]
}

func (t *TradePair) depthTicker(que *OrderQueue) {

	ticker := time.NewTicker(time.Duration(100) * time.Millisecond)

	for {
		<-ticker.C
		func() {
			t.w.Lock()
			defer t.w.Unlock()

			que.Lock()
			defer que.Unlock()
			que.depth = [][2]string{}
			depthMap := make(map[string]string)

			if que.pq.Len() > 0 {

				for i := 0; i < que.pq.Len(); i++ {

					if len(depthMap) > t.depthMaxLen {
						break
					}

					item := (*que.pq)[i]

					price := utils.D2Str(item.GetPrice(), t.priceDigit)

					if _, ok := depthMap[price]; !ok {
						depthMap[price] = utils.D2Str(item.GetQuantity(), t.quantityDigit)
					} else {
						old_qunantity, _ := decimal.NewFromString(depthMap[price])
						depthMap[price] = utils.D2Str(old_qunantity.Add(item.GetQuantity()), t.quantityDigit)
					}

				}

				//按价格排序map
				que.depth = sortMap2Slice(depthMap, que.Top().GetOrderSide())
			}
		}()
	}
}
