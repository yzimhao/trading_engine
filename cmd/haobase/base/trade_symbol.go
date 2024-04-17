package base

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils/app"
)

var (
	tsymbol *TradeSymbol
)

type TradeSymbol struct {
}

func NewTradeSymbol() *TradeSymbol {
	if tsymbol != nil {
		return tsymbol
	}
	tsymbol = &TradeSymbol{}

	return tsymbol
}

func (t *TradeSymbol) all_symbols() []varieties.TradingVarieties {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	var data []varieties.TradingVarieties

	raw, _ := redis.Bytes(rdc.Do("get", redisdb.BaseTradeSymbolAll.Format(redisdb.Replace{})))
	if err := json.Unmarshal(raw, &data); err != nil {
		data = varieties.AllTradingVarieties()
		raw, _ = json.Marshal(data)
		rdc.Do("set", redisdb.BaseTradeSymbolAll.Format(redisdb.Replace{}), raw)
	}

	return data
}

func (t *TradeSymbol) Get(symbol string) (*varieties.TradingVarieties, error) {
	for _, item := range t.all_symbols() {
		if item.Symbol == symbol {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("不存在的交易对symbol")
}

func (t *TradeSymbol) All() []varieties.TradingVarieties {
	return t.all_symbols()
}
