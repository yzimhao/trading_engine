package types

import "time"

type KLine struct {
	Symbol  string
	OpenAt  time.Time //开盘时间
	CloseAt time.Time // 收盘时间
	Open    *string   //开盘价
	High    *string   // 最高价
	Low     *string   //最低价
	Close   *string   //收盘价(当前K线未结束的即为最新价)
	Volume  *string   //成交量
	Amount  *string   //成交额
	Period  PeriodType
}
