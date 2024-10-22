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
	Data          types.KLine
	OpenLastTime  int64
	CloseLastTime int64
}

type KLinePeriod interface {
	GetData(ctx context.Context, periodType types.PeriodType, tradeResult matching_types.TradeResult) (*types.KLine, error)
	CleanCache(ctx context.Context, openAt, closeAt time.Time) error
}

type kLine struct {
	redis  *redis.Client
	symbol string
	logger *zap.Logger
}

var _ KLinePeriod = (*kLine)(nil)

func NewKLine(cli *redis.Client, logger *zap.Logger, symbol string) KLinePeriod {
	return &kLine{
		redis:  cli,
		symbol: symbol,
		logger: logger,
	}
}

func cacheKey(symbol string, openAt, closeAt time.Time) string {
	return fmt.Sprintf("kline:%s:%d:%d", symbol, openAt.Unix(), closeAt.Unix())
}

func (k *kLine) CleanCache(ctx context.Context, openAt, closeAt time.Time) error {
	key := cacheKey(k.symbol, openAt, closeAt)
	return k.redis.Del(ctx, key).Err()
}

func (k *kLine) GetData(ctx context.Context, periodType types.PeriodType, tradeResult matching_types.TradeResult) (*types.KLine, error) {
	tradeTime := time.Unix(int64(tradeResult.TradeTime/1e9), 0)
	openAt, closeAt := types.ParsePeriodTime(tradeTime, periodType)

	key := cacheKey(k.symbol, openAt, closeAt)

	//lock
	lockKey := fmt.Sprintf("lock:%s", key)
	lock, err := k.redis.SetNX(ctx, lockKey, 1, 10*time.Second).Result()
	if err != nil {
		k.logger.Error("[kline] set lock for kline calculation failed", zap.Error(err))
		return nil, err
	}
	if !lock {
		k.logger.Warn("[kline] failed to acquire lock for kline calculation")
		return nil, fmt.Errorf("failed to acquire lock for kline calculation")
	}
	defer k.redis.Del(ctx, lockKey)

	cache, err := k.redis.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		k.logger.Error("[kline] get cache data from redis failed", zap.Error(err))
		return nil, err
	}

	var cacheData kline
	if err == redis.Nil {
		cacheData = kline{}
	} else {
		if err := json.Unmarshal([]byte(cache), &cacheData); err != nil {
			k.logger.Error("[kline] unmarshal kline cachedata from redis failed", zap.Error(err))
			return nil, err
		}
	}
	data := kline{
		Data: types.KLine{
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
		},
		OpenLastTime:  cacheData.OpenLastTime,
		CloseLastTime: cacheData.CloseLastTime,
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
	ttl := data.Data.CloseAt.Unix() - time.Now().Unix() + 3600*24
	if ttl < 0 {
		ttl = 3600 * 24
	}

	err = k.redis.Expire(ctx, key, time.Duration(ttl)*time.Second).Err()
	if err != nil {
		k.logger.Error("[kline] update kline key ttl failed", zap.Error(err))
		return nil, err
	}

	return &data.Data, nil
}

func (k *kLine) getOpen(cacheData *kline, tradeResult *matching_types.TradeResult) *string {

	if cacheData.Data.Open == nil {
		cacheData.Data.Open = &tradeResult.TradePrice
		cacheData.OpenLastTime = tradeResult.TradeTime
	} else {

		if tradeResult.TradeTime < cacheData.OpenLastTime {
			cacheData.Data.Open = &tradeResult.TradePrice
			cacheData.OpenLastTime = tradeResult.TradeTime
		}
	}

	return cacheData.Data.Open
}

func (k *kLine) getHigh(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.Data.High == nil {
		cacheData.Data.High = &tradeResult.TradePrice
	} else {

		if k.formatD(tradeResult.TradePrice).Cmp(k.formatD(*cacheData.Data.High)) > 0 {
			cacheData.Data.High = &tradeResult.TradePrice
		}
	}

	return cacheData.Data.High
}

func (k *kLine) getLow(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.Data.Low == nil {
		cacheData.Data.Low = &tradeResult.TradePrice
	} else {

		if k.formatD(tradeResult.TradePrice).Cmp(k.formatD(*cacheData.Data.Low)) < 0 {
			cacheData.Data.Low = &tradeResult.TradePrice
		}
	}

	return cacheData.Data.Low
}

func (k *kLine) getClose(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.Data.Close == nil {
		cacheData.Data.Close = &tradeResult.TradePrice
		cacheData.CloseLastTime = tradeResult.TradeTime
	} else {

		if tradeResult.TradeTime > cacheData.CloseLastTime {
			cacheData.Data.Close = &tradeResult.TradePrice
			cacheData.CloseLastTime = tradeResult.TradeTime
		}
	}
	return cacheData.Data.Close
}

func (k *kLine) getVolume(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	if cacheData.Data.Volume == nil {
		cacheData.Data.Volume = &tradeResult.TradeQuantity
	} else {
		volume := k.formatD(*cacheData.Data.Volume).Add(k.formatD(tradeResult.TradeQuantity)).String()
		cacheData.Data.Volume = &volume
	}
	return cacheData.Data.Volume
}

func (k *kLine) getAmount(cacheData *kline, tradeResult *matching_types.TradeResult) *string {
	amount := k.formatD(tradeResult.TradePrice).Mul(k.formatD(tradeResult.TradeQuantity)).String()

	if cacheData.Data.Amount == nil {
		cacheData.Data.Amount = &amount
	} else {
		amount = k.formatD(*cacheData.Data.Amount).Add(k.formatD(amount)).String()
		cacheData.Data.Amount = &amount
	}

	return cacheData.Data.Amount
}

func (k *kLine) formatD(d1 string) decimal.Decimal {
	d, err := decimal.NewFromString(d1)
	if err != nil {
		k.logger.Sugar().Errorf("[kline] new decimal from string failed d1: %s error: %v", d1, zap.Error(err))
		return decimal.Zero
	}
	return d
}
