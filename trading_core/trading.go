package trading_core

import "github.com/shopspring/decimal"

func (t *TradePair) Price2String(price decimal.Decimal) string {
	return FormatDecimal2String(price, t.priceDigit)
}

func (t *TradePair) Qty2String(qty decimal.Decimal) string {
	return FormatDecimal2String(qty, t.quantityDigit)
}
