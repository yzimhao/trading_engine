package base

import (
	"fmt"

	"github.com/yzimhao/trading_engine/cmd/haobase/base/symbols"
)

var (
	tsymbol *TSymbols
)

func Init() {
	symbols.Init()
}

type TSymbols struct {
	all_symbols []symbols.TradingVarieties
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
	t.all_symbols = symbols.AllTradingVarieties()
}

func (t *TSymbols) Get(symbol string) (*symbols.TradingVarieties, error) {
	for _, item := range t.all_symbols {
		if item.Symbol == symbol {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("不存在的交易对symbol")
}

func (t *TSymbols) All() []symbols.TradingVarieties {
	return t.all_symbols
}
