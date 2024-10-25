package order_test

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	gorm_order "github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/order"
	mock_asset "github.com/yzimhao/trading_engine/v2/mocks/persistence/asset"
	mock_trade_variety "github.com/yzimhao/trading_engine/v2/mocks/persistence/trade_variety"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	_gorm "gorm.io/gorm"
)

type orderRepoTest struct {
	suite.Suite
	ctx                  context.Context
	ctrl                 *gomock.Controller
	repo                 persistence.OrderRepository
	v                    *viper.Viper
	gorm                 *_gorm.DB
	logger               *zap.Logger
	mockTradeVarietyRepo mock_trade_variety.MockTradeVarietyRepository
	mockAssetRepo        mock_asset.MockAssetRepository
}

func (suite *orderRepoTest) SetupTest() {
	_ = gotenv.Load("../../../.env")

	suite.ctx = context.Background()
	suite.ctrl = gomock.NewController(suite.T())
	suite.v = di.NewViper()
	suite.gorm = di.NewGorm(suite.v)
	suite.logger = zap.NewNop()
	// redis := di.NewRedis(suite.v, suite.logger)
	// cache, _ := di.NewCache(suite.v, redis)
	logger := zap.NewNop()
	mockTradeVarietyRepo := mock_trade_variety.NewMockTradeVarietyRepository(suite.ctrl)
	mockAssetRepo := mock_asset.NewMockAssetRepository(suite.ctrl)
	suite.repo = gorm_order.NewOrderRepo(suite.gorm, logger, mockTradeVarietyRepo, mockAssetRepo)
}

func TestOrderRepo(t *testing.T) {
	suite.Run(t, new(orderRepoTest))
}

func (suite *orderRepoTest) TearDownTest() {
	// migrations.MigrateDown(suite.gorm, suite.v, suite.logger)
	suite.ctrl.Finish()
}

func (suite *orderRepoTest) TestCreateOrder() {
	suite.mockTradeVarietyRepo.EXPECT().FindBySymbol(suite.ctx, "BTCUSDT").Return(&models_variety.TradeVariety{
		ID:            1,
		Symbol:        "BTCUSDT",
		Name:          "BTCUSDT",
		TargetId:      2,
		BaseId:        1,
		PriceDecimals: 2,
		QtyDecimals:   4,
		FeeRate:       "0.01",
		TargetVariety: &models_variety.Variety{
			ID:   2,
			Name: "BTC",
		},
		BaseVariety: &models_variety.Variety{
			ID:   1,
			Name: "USDT",
		},
	}, nil).AnyTimes()

}
