package trading_engine

import (
	"time"

	"github.com/shopspring/decimal"
)

func (t *TradePair) GetAskDepth(limit int) [][2]string {
	return t.askQueue.GetDepth(limit)
}

func (t *TradePair) GetBidDepth(limit int) [][2]string {
	return t.bidQueue.GetDepth(limit)
}

func (t *TradePair) depthTicker(que *OrderQueue) {

	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)

	for {
		<-ticker.C
		func() {
			que.Lock()
			defer que.Unlock()
			que.depth = [][2]string{}
			depthMap := make(map[string]string)

			if que.pq.Len() > 0 {

				for i := 0; i < que.pq.Len(); i++ {
					item := (*que.pq)[i]

					price := FormatDecimal2String(item.GetPrice(), t.priceDigit)

					if _, ok := depthMap[price]; !ok {
						depthMap[price] = FormatDecimal2String(item.GetQuantity(), t.quantityDigit)
					} else {
						old_qunantity, _ := decimal.NewFromString(depthMap[price])
						depthMap[price] = FormatDecimal2String(old_qunantity.Add(item.GetQuantity()), t.quantityDigit)
					}
				}

				//按价格排序map
				que.depth = sortMap2Slice(depthMap, que.Top().GetOrderSide())
			}
		}()
	}
}
