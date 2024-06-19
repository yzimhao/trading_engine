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

func (t *TradeSymbol) all_symbols() ([]varieties.TradingVarieties, error) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()

	var data []varieties.TradingVarieties

	raw, err := redis.Bytes(rdc.Do("get", redisdb.BaseTradeSymbolAll.Format(redisdb.Replace{})))

	if err != nil {
		return []varieties.TradingVarieties{}, err
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		data = varieties.AllTradingVarieties()
		if raw, err = json.Marshal(data); err != nil {
			return []varieties.TradingVarieties{}, err
		}

		if _, err := rdc.Do("set", redisdb.BaseTradeSymbolAll.Format(redisdb.Replace{}), raw); err != nil {
			return []varieties.TradingVarieties{}, err
		}

	}

	return data, nil
}

func (t *TradeSymbol) Get(symbol string) (*varieties.TradingVarieties, error) {
	symbols, err := t.all_symbols()
	if err != nil {
		return nil, err
	}

	for _, item := range symbols {
		if item.Symbol == symbol {
			return &item, nil
		}
	}
	return nil, fmt.Errorf("不存在的交易对symbol")
}

func (t *TradeSymbol) All() []varieties.TradingVarieties {
	a, err := t.all_symbols()
	if err != nil {
		app.Logger.Warnf("获取交易对失败: %s", err)
	}
	return a
}
