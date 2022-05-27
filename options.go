package trading_engine

import "fmt"

func SetPriceDigits(digits int) {
	priceDigits = digits
	priceFormat = "%." + fmt.Sprintf("%d", priceDigits) + "f"
}

func SetQuantityDigits(digits int) {
	quantityDigits = digits
	quantityFormat = "%." + fmt.Sprintf("%d", quantityDigits) + "f"
}
