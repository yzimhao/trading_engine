package trading_core

import (
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/utils"
)

func (t *TradePair) Price2String(price decimal.Decimal) string {
	return utils.D2Str(price, t.priceDigit)
}

func (t *TradePair) Qty2String(qty decimal.Decimal) string {
	return utils.D2Str(qty, t.quantityDigit)
}
