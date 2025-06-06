package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type KLine struct {
	Symbol  string
	OpenAt  time.Time        //开盘时间
	CloseAt time.Time        // 收盘时间
	Open    *decimal.Decimal //开盘价
	High    *decimal.Decimal // 最高价
	Low     *decimal.Decimal //最低价
	Close   *decimal.Decimal //收盘价(当前K线未结束的即为最新价)
	Volume  *decimal.Decimal //成交量
	Amount  *decimal.Decimal //成交额
	Period  PeriodType
}
