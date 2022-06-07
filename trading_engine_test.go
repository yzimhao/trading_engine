package trading_engine

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var btcusdt = NewTradePair("btcusdt", 2, 0)

func init() {

}

func TestTradeFunc(t *testing.T) {

}

func TestTradeFunc_LimitOrder(t *testing.T) {
	btcusdt.PushNewOrder(NewBidLimitItem("uid1", decimal.NewFromFloat(1.1), decimal.NewFromFloat(1.2), 1112))
	assert.Equal(t, 0, btcusdt.askQueue.Len())
	assert.Equal(t, 1, btcusdt.bidQueue.Len())
	assert.Equal(t, "uid1", btcusdt.bidQueue.Top().GetUniqueId())

	time.Sleep(time.Duration(100) * time.Millisecond)
	a := btcusdt.GetBidDepth(0)
	stra, _ := json.Marshal(a)
	assert.Equal(t, "[[\"1.10\",\"1\"]]", string(stra))

	btcusdt.PushNewOrder(NewAskLimitItem("uid2", decimal.NewFromFloat(1.1), decimal.NewFromFloat(1.2), 1112))

	assert.Equal(t, 1, btcusdt.askQueue.Len())
	assert.Equal(t, "uid2", btcusdt.askQueue.Top().GetUniqueId())

	tradeLog := <-btcusdt.ChTradeResult
	assert.Equal(t, "uid2", tradeLog.AskOrderId)
	assert.Equal(t, "uid1", tradeLog.BidOrderId)
	assert.Equal(t, "1.1", tradeLog.TradePrice.String())
	assert.Equal(t, "1.2", tradeLog.TradeQuantity.String())

	btcusdt.PushNewOrder(NewBidLimitItem("uid3", decimal.NewFromFloat(1.01), decimal.NewFromFloat(3.0), 1112))
	btcusdt.PushNewOrder(NewAskLimitItem("uid4", decimal.NewFromFloat(0.9), decimal.NewFromFloat(1.0), 1113))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, "1.01", tradeLog.TradePrice.String())
	assert.Equal(t, "1", tradeLog.TradeQuantity.String())
	assert.Equal(t, "2", btcusdt.bidQueue.Top().GetQuantity().String())

	btcusdt.PushNewOrder(NewBidLimitItem("uid5", decimal.NewFromFloat(1.02), decimal.NewFromFloat(3.0), 1114))
	btcusdt.PushNewOrder(NewAskLimitItem("uid6", decimal.NewFromFloat(0.1), decimal.NewFromFloat(5.0), 1115))
	tradeLog = <-btcusdt.ChTradeResult
	assert.Equal(t, 0, btcusdt.bidQueue.Len())
	assert.Equal(t, 0, btcusdt.askQueue.Len())
}

func TestTradeFunc_MarketOrder(t *testing.T) {

}
