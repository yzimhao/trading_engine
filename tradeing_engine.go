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
	AskDepth      [][2]string
	BidDepth      [][2]string
	ChTradeResult chan TradeResult

	askQueue *OrderQueue
	bidQueue *OrderQueue
}

func NewTradePair(symbol string, priceDigit, quantityDigit int) *TradePair {
	t := &TradePair{
		Symbol:        symbol,
		ChTradeResult: make(chan TradeResult, 100),

		askQueue: NewQueue(priceDigit, quantityDigit),
		bidQueue: NewQueue(priceDigit, quantityDigit),
	}
	t.matching()
	return t
}

func (t *TradePair) PushNewOrder(side OrderSide, order QueueItem) {
	if side == OrderSideSell {
		t.askQueue.Push(order)
	} else {
		t.bidQueue.Push(order)
	}
}

func (t *TradePair) GetAskDepth() [][2]string {
	return t.askQueue.GetDepth()
}

func (t *TradePair) GetBidDepth() [][2]string {
	return t.bidQueue.GetDepth()
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
