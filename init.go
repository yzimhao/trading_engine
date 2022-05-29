package trading_engine

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type PriceType int
type OrderSide int

const (
	PriceTypeLimit  PriceType = 0
	PriceTypeMarket PriceType = 1

	OrderSideBuy  OrderSide = 0
	OrderSideSell OrderSide = 1
)

var (
	priceDigits    = 2
	quantityDigits = 0

	priceFormat    = "%." + fmt.Sprintf("%d", priceDigits) + "f"
	quantityFormat = "%." + fmt.Sprintf("%d", quantityDigits) + "f"
)

func FormatPrice2Str(price decimal.Decimal) string {
	return formatDecimal(priceFormat, price)
}

func FormatQuantity2Str(quantity decimal.Decimal) string {
	return formatDecimal(quantityFormat, quantity)
}

func formatDecimal(format string, d decimal.Decimal) string {
	f, _ := d.Float64()
	return fmt.Sprintf(format, f)
}

func quickSort(nums []string, asc_desc string) []string {
	if len(nums) <= 1 {
		return nums
	}

	spilt := nums[0]
	left := []string{}
	right := []string{}
	mid := []string{}

	for _, v := range nums {
		vv, _ := decimal.NewFromString(v)
		sp, _ := decimal.NewFromString(spilt)
		if vv.Cmp(sp) == -1 {
			left = append(left, v)
		} else if vv.Cmp(sp) == 1 {
			right = append(right, v)
		} else {
			mid = append(mid, v)
		}
	}

	left = quickSort(left, asc_desc)
	right = quickSort(right, asc_desc)

	if asc_desc == "asc" {
		return append(append(left, mid...), right...)
	} else {
		return append(append(right, mid...), left...)
	}

	//return append(append(left, mid...), right...)
}
