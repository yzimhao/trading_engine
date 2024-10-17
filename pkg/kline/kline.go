package kline

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
)

type kline struct {
	types.KLine
	OpenLastTime  int64
	CloseLastTime int64
}

type KLinePeriod interface {
	GetData(ctx context.Context, periodType types.PeriodType, tradeResult matching_types.TradeResult) (*types.KLine, error)
}

type kLine struct {
	redis  *redis.Client
	symbol string
	logger *zap.Logger
}

var _ KLinePeriod = &kLine{}

func NewKLinePeriod(cli *redis.Client, logger *zap.Logger, symbol string) KLinePeriod {
	return &kLine{
		redis:  cli,
		symbol: symbol,
		logger: logger,
	}
}

func cacheKey(symbol string, openAt, closeAt time.Time) string {
	return fmt.Sprintf("kline:%s:%d:%d", symbol, openAt.Unix(), closeAt.Unix())
}

func (k *kLine) GetData(ctx context.Context, periodType types.PeriodType, tradeResult matching_types.TradeResult) (*types.KLine, error) {
	tradeTime := time.Unix(int64(tradeResult.TradeTime/1e9), 0)
	openAt, closeAt := types.ParsePeriodTime(tradeTime, periodType)

	key := cacheKey(k.symbol, openAt, closeAt)
	cache, err := k.redis.Get(ctx, key).Result()
	if err != nil {
		k.logger.Error("[kline] get cache data from redis failed", zap.Error(err))
		return nil, err
	}

	var cacheData kline
	if err := json.Unmarshal([]byte(cache), &cacheData); err != nil {
		k.logger.Error("[kline] unmarshal kline cachedata from redis failed", zap.Error(err))
		return nil, err
	}

	data := types.KLine{
		Symbol:  k.symbol,
		OpenAt:  openAt,
		CloseAt: closeAt,
		Period:  periodType,
		Open:    k.getOpen(&cacheData, &tradeResult),
		High:    k.getHigh(&cacheData, &tradeResult),
		Low:     k.getLow(&cacheData, &tradeResult),
		Close:   k.getClose(&cacheData, &tradeResult),
		Volume:  k.getVolume(&cacheData, &tradeResult),
		Amount:  k.getAmount(&cacheData, &tradeResult),
	}

	dataJson, err := json.Marshal(data)
	if err != nil {
		k.logger.Error("[kline] marshal kline data to redis failed", zap.Error(err))
		return nil, err
	}

	err = k.redis.Set(ctx, key, string(dataJson), 0).Err()
	if err != nil {
		k.logger.Error("[kline] set kline data to redis failed", zap.Error(err))
		return nil, err
	}

	//update key ttl
	ttl := data.CloseAt.Unix() - time.Now().Unix() + 3600*24
	err = k.redis.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		k.logger.Error("[kline] update kline key ttl failed", zap.Error(err))
		return nil, err
	}

	return &data, nil
}

func (k *kLine) getOpen(cacheData *kline, tradeResult *matching_types.TradeResult) *string {

	if cacheData.Open == nil {
		cacheData.Open = &tradeResult.TradePrice
		cacheData.OpenLastTime = tradeResult.TradeTime
	} else {

		if tradeResult.TradeTime < cacheData.OpenLastTime {
			cacheData.Open = &tradeResult.TradePrice
			cacheData.OpenLastTime = tradeResult.TradeTime
		}
	}

	return cacheData.Open
}

func (k *kLine) getHigh(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.High == nil {
		cacheData.High = &tradeResult.TradePrice
	} else {

		if k.formatD(tradeResult.TradePrice).Cmp(k.formatD(*cacheData.High)) > 0 {
			cacheData.High = &tradeResult.TradePrice
		}
	}

	return cacheData.High
}

func (k *kLine) getLow(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.Low == nil {
		cacheData.Low = &tradeResult.TradePrice
	} else {

		if k.formatD(tradeResult.TradePrice).Cmp(k.formatD(*cacheData.Low)) < 0 {
			cacheData.Low = &tradeResult.TradePrice
		}
	}

	return cacheData.Low
}

func (k *kLine) getClose(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.Close == nil {
		cacheData.Close = &tradeResult.TradePrice
		cacheData.CloseLastTime = tradeResult.TradeTime
	} else {

		if tradeResult.TradeTime > cacheData.CloseLastTime {
			cacheData.Close = &tradeResult.TradePrice
			cacheData.CloseLastTime = tradeResult.TradeTime
		}
	}
	return cacheData.Close
}

func (k *kLine) getVolume(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.Volume == nil {
		cacheData.Volume = &tradeResult.TradeQuantity
	} else {
		volume := k.formatD(*cacheData.Volume).Add(k.formatD(tradeResult.TradeQuantity)).String()
		cacheData.Volume = &volume
	}
	return cacheData.Volume
}

func (k *kLine) getAmount(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	amount := k.formatD(tradeResult.TradePrice).Mul(k.formatD(tradeResult.TradeQuantity)).String()

	if cacheData.Amount == nil {
		cacheData.Amount = &amount
	} else {
		amount = k.formatD(*cacheData.Amount).Add(k.formatD(amount)).String()
		cacheData.Amount = &amount
	}

	return cacheData.Amount
}

func (k *kLine) formatD(d1 string) decimal.Decimal {
	d, err := decimal.NewFromString(d1)
	if err != nil {
		k.logger.Sugar().Errorf("[kline] new decimal from string failed d1: %s error: %v", d1, zap.Error(err))
		return decimal.Zero
	}
	return d
}
