package base

import (
	"fmt"

	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
)

var (
	symbol *Symbols
)

type Symbols struct {
	all_varieties []varieties.Varieties
}

func NewSymbols() *Symbols {
	if symbol != nil {
		return symbol
	}
	symbol = &Symbols{}
	symbol.init()
	return symbol
}

func (t *Symbols) init() {
	t.all_varieties = varieties.AllVarieties()
}

func (t *Symbols) Get(symbol string) (*varieties.Varieties, error) {
	for _, item := range t.all_varieties {
		if item.Symbol == symbol {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("不存在的交易对symbol")
}

func (t *Symbols) All() []varieties.Varieties {
	t.init()
	return t.all_varieties
}
