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
	TradeTime     time.Time
}

var ChTradeResult = make(chan TradeResult, 100)

func MatchingEngine(askQ *OrderQueue, bidQ *OrderQueue) {
	logrus.Infof("MatchingEngine start")
	go func() {
		for {
			ok := func() bool {
				if askQ == nil || bidQ == nil {
					logrus.Warningf("askQ or bidQ is nil")
					return false
				}

				if askQ.Len() == 0 || bidQ.Len() == 0 {
					logrus.Warningf("askQ or bidQ is empty. askLen:%d, bidLen:%d", askQ.Len(), bidQ.Len())
					return false
				}

				askTop := askQ.Top()
				bidTop := bidQ.Top()

				if bidTop.GetPrice().Cmp(askTop.GetPrice()) >= 0 {
					tradelog := TradeResult{}

					if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == 0 {
						tradelog.TradeQuantity = askTop.GetQuantity()

						askQ.Remove(askTop.GetUniqueId())
						bidQ.Remove(bidTop.GetUniqueId())
					} else if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == 1 {
						tradelog.TradeQuantity = askTop.GetQuantity()

						askQ.Remove(askTop.GetUniqueId())
						bidTop.SetQuantity(bidTop.GetQuantity().Sub(askTop.GetQuantity()))
					} else if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == -1 {
						tradelog.TradeQuantity = bidTop.GetQuantity()

						bidQ.Remove(bidTop.GetUniqueId())
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

					//通知交易结果
					logrus.Infof("tradelog: %+v", tradelog)
					ChTradeResult <- tradelog
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
