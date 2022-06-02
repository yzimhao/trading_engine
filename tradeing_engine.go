package trading_engine

import (
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type TradeResult struct {
	AskOrderId    string
	BidOrderId    string
	TradeQuantity decimal.Decimal
	TradePrice    decimal.Decimal
	TradeAmount   decimal.Decimal
	TradeTime     time.Time
}

type TradePair struct {
	Symbol        string
	ChTradeResult chan TradeResult

	PriceDigit    int
	QuantityDigit int

	askQueue *OrderQueue
	bidQueue *OrderQueue
}

func NewTradePair(symbol string, priceDigit, quantityDigit int) *TradePair {
	t := &TradePair{
		Symbol:        symbol,
		ChTradeResult: make(chan TradeResult, 100),

		PriceDigit:    priceDigit,
		QuantityDigit: quantityDigit,

		askQueue: NewQueue(priceDigit, quantityDigit),
		bidQueue: NewQueue(priceDigit, quantityDigit),
	}
	t.matching()
	return t
}

func (t *TradePair) PushNewOrder(side OrderSide, order QueueItem) {
	//todo 先将新订单和对手单比较，能否成交，不能成交的才放入队列
	if side == OrderSideSell {
		t.askQueue.Push(order)
	} else {
		t.bidQueue.Push(order)
	}
}

func (t *TradePair) CancelOrder(side OrderSide, uniq string) {
	//todo 最好根据订单编号知道是买单还是卖单，方便直接查找到相应的队列，从中删除
	if side == OrderSideSell {
		t.askQueue.Remove(uniq)
	} else {
		t.bidQueue.Remove(uniq)
	}
	//todo 删除成功后需要发送通知
}

func (t *TradePair) GetAskDepth(limit int) [][2]string {
	return t.askQueue.GetDepth(limit)
}

func (t *TradePair) GetBidDepth(limit int) [][2]string {
	return t.bidQueue.GetDepth(limit)
}

func (t *TradePair) AskLen() int {
	return t.askQueue.Len()
}

func (t *TradePair) BidLen() int {
	return t.bidQueue.Len()
}

func (t *TradePair) matching() {
	go func() {
		for {
			ok := func() bool {
				if t.askQueue == nil || t.bidQueue == nil {
					logrus.Warningf("%s askQueue or bidQueue is nil", t.Symbol)
					return false
				}

				if t.askQueue.Len() == 0 || t.bidQueue.Len() == 0 {
					logrus.Warningf("%s askQueue or bidQueue is empty. askLen:%d, bidLen:%d", t.Symbol, t.askQueue.Len(), t.bidQueue.Len())
					return false
				}

				askTop := t.askQueue.Top()
				bidTop := t.bidQueue.Top()

				if bidTop.GetPrice().Cmp(askTop.GetPrice()) >= 0 {
					tradelog := TradeResult{}

					if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == 0 {
						tradelog.TradeQuantity = askTop.GetQuantity()

						t.askQueue.Remove(askTop.GetUniqueId())
						t.bidQueue.Remove(bidTop.GetUniqueId())
					} else if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == 1 {
						tradelog.TradeQuantity = askTop.GetQuantity()

						t.askQueue.Remove(askTop.GetUniqueId())
						bidTop.SetQuantity(bidTop.GetQuantity().Sub(askTop.GetQuantity()))
					} else if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == -1 {
						tradelog.TradeQuantity = bidTop.GetQuantity()

						t.bidQueue.Remove(bidTop.GetUniqueId())
						askTop.SetQuantity(askTop.GetQuantity().Sub(bidTop.GetQuantity()))
					}

					if askTop.GetCreateTime() >= bidTop.GetCreateTime() {
						tradelog.TradePrice = bidTop.GetPrice()
					} else {
						tradelog.TradePrice = askTop.GetPrice()
					}

					tradelog.TradeTime = time.Now()
					tradelog.AskOrderId = askTop.GetUniqueId()
					tradelog.BidOrderId = bidTop.GetUniqueId()
					tradelog.TradeAmount = tradelog.TradeQuantity.Mul(tradelog.TradePrice)

					//通知交易结果
					logrus.Infof("%s tradelog: %+v", t.Symbol, tradelog)
					t.ChTradeResult <- tradelog
					return true
				} else {
					return false
				}

			}()

			if !ok {
				time.Sleep(time.Duration(200) * time.Millisecond)
			}
		}
	}()
}
