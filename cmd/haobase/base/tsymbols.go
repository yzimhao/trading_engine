package base

import (
	"fmt"

	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
)

var (
	tsymbol *TSymbols
)

type TSymbols struct {
	all_symbols []varieties.TradingVarieties
}

func NewTSymbols() *TSymbols {
	if tsymbol != nil {
		return tsymbol
	}
	tsymbol = &TSymbols{}
	tsymbol.init()
	return tsymbol
}

func (t *TSymbols) init() {
	t.all_symbols = varieties.AllTradingVarieties()
}

func (t *TSymbols) Get(symbol string) (*varieties.TradingVarieties, error) {
	for _, item := range t.all_symbols {
		if item.Symbol == symbol {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("不存在的交易对symbol")
}

func (t *TSymbols) All() []varieties.TradingVarieties {
	return t.all_symbols
}
