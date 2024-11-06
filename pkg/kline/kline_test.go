package kline_test

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	"github.com/yzimhao/trading_engine/v2/pkg/kline"
	"github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
)

type klineTest struct {
	suite.Suite
	ctx     context.Context
	symbol  string
	service kline.KLinePeriod
}

func init() {
	_ = gotenv.Load()
}

func TestKLine(t *testing.T) {
	suite.Run(t, new(klineTest))
}

func (suite *klineTest) SetupTest() {
	suite.ctx = context.Background()

	v := di.NewViper()
	logger := zap.NewNop()
	redis := di.NewRedis(v, logger)
	suite.symbol = "BTCUSDT"
	suite.service = kline.NewKLine(redis, logger, suite.symbol)
}

func (suite *klineTest) TearDownTest() {}

func (suite *klineTest) Test_GetKLine() {
	testCases := []struct {
		name  string
		setup func()
	}{
		{
			name: "kline第一个成交记录",
			setup: func() {
				tradeTime, err := time.Parse("2006-01-02 15:04:05", "2024-05-01 10:00:00")
				suite.Require().NoError(err)
				tradeResult := matching_types.TradeResult{
					Symbol:        suite.symbol,
					AskOrderId:    "ask1",
					BidOrderId:    "bid1",
					TradePrice:    decimal.NewFromFloat(1.00),
					TradeQuantity: decimal.NewFromFloat(10),
					TradeTime:     tradeTime.UnixNano(),
				}
				kline, err := suite.service.GetData(suite.ctx, types.PERIOD_M1, tradeResult)
				suite.Require().NoError(err)
				suite.Equal(*kline.Open, "1.00")
				suite.Equal(*kline.High, "1.00")
				suite.Equal(*kline.Low, "1.00")
				suite.Equal(*kline.Close, "1.00")
				suite.Equal(*kline.Volume, "10")
				suite.Equal(*kline.Amount, "10")
				suite.service.CleanCache(suite.ctx, kline.OpenAt, kline.CloseAt)

			},
		},

		{
			name: "kline成交记录不按先后顺序到达",
			setup: func() {
				tradeTime, err := time.Parse("2006-01-02 15:04:05", "2024-05-01 10:00:15")
				suite.Require().NoError(err)
				tradeResult := matching_types.TradeResult{
					Symbol:        suite.symbol,
					AskOrderId:    "ask1",
					BidOrderId:    "bid1",
					TradePrice:    decimal.NewFromFloat(1.00),
					TradeQuantity: decimal.NewFromFloat(10),
					TradeTime:     tradeTime.UnixNano(),
				}

				tradeTime1, err := time.Parse("2006-01-02 15:04:05", "2024-05-01 10:00:05")
				suite.Require().NoError(err)
				tradeResult1 := matching_types.TradeResult{
					Symbol:        suite.symbol,
					AskOrderId:    "ask1",
					BidOrderId:    "bid1",
					TradePrice:    decimal.NewFromFloat(2.00),
					TradeQuantity: decimal.NewFromFloat(10),
					TradeTime:     tradeTime1.UnixNano(),
				}

				tradeTime2, err := time.Parse("2006-01-02 15:04:05", "2024-05-01 10:00:30")
				suite.Require().NoError(err)
				tradeResult2 := matching_types.TradeResult{
					Symbol:        suite.symbol,
					AskOrderId:    "ask1",
					BidOrderId:    "bid1",
					TradePrice:    decimal.NewFromFloat(0.95),
					TradeQuantity: decimal.NewFromFloat(10),
					TradeTime:     tradeTime2.UnixNano(),
				}

				_, err = suite.service.GetData(suite.ctx, types.PERIOD_M1, tradeResult)
				suite.Require().NoError(err)
				_, err = suite.service.GetData(suite.ctx, types.PERIOD_M1, tradeResult1)
				suite.Require().NoError(err)
				kline, err := suite.service.GetData(suite.ctx, types.PERIOD_M1, tradeResult2)
				suite.Require().NoError(err)
				suite.Equal(*kline.Open, "2.00")
				suite.Equal(*kline.High, "2.00")
				suite.Equal(*kline.Low, "0.95")
				suite.Equal(*kline.Close, "0.95")
				suite.Equal(*kline.Volume, "30")
				suite.Equal(*kline.Amount, "39.5") //0.95 *10 + 2 *10 + 1* 10
				suite.service.CleanCache(suite.ctx, kline.OpenAt, kline.CloseAt)

			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setup()
		})
	}
}
