package order_test

import (
	"context"
	"testing"

	"github.com/duolacloud/crud-core/datasource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	gorm_asset "github.com/yzimhao/trading_engine/v2/internal/persistence/gorm"
	gorm_order "github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/order"
	mock_asset "github.com/yzimhao/trading_engine/v2/mocks/persistence/asset"
	mock_trade_variety "github.com/yzimhao/trading_engine/v2/mocks/persistence/trade_variety"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
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
	mockTradeVarietyRepo *mock_trade_variety.MockTradeVarietyRepository
	mockAssetRepo        *mock_asset.MockAssetRepository
	assetRepo            persistence.AssetRepository
}

func (suite *orderRepoTest) SetupTest() {
	_ = gotenv.Load("../../../../.env")

	suite.ctx = context.Background()
	suite.ctrl = gomock.NewController(suite.T())
	suite.v = di.NewViper()
	suite.gorm = di.NewGorm(suite.v)
	suite.logger = zap.NewNop()
	redis := di.NewRedis(suite.v, suite.logger)
	cache, _ := di.NewCache(suite.v, redis)
	logger := zap.NewNop()
	mockTradeVarietyRepo := mock_trade_variety.NewMockTradeVarietyRepository(suite.ctrl)
	mockAssetRepo := mock_asset.NewMockAssetRepository(suite.ctrl)
	suite.mockTradeVarietyRepo = mockTradeVarietyRepo
	suite.mockAssetRepo = mockAssetRepo
	suite.assetRepo = gorm_asset.NewAssetRepo(datasource.NewDataSource(suite.gorm), cache, logger)
	suite.repo = gorm_order.NewOrderRepo(suite.gorm, logger, suite.mockTradeVarietyRepo, suite.assetRepo)
}

func TestOrderRepo(t *testing.T) {
	suite.Run(t, new(orderRepoTest))
}

func (suite *orderRepoTest) TearDownTest() {
	// migrations.MigrateDown(suite.gorm, suite.v, suite.logger)
	suite.ctrl.Finish()
}

func (suite *orderRepoTest) initMockData() {
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
			ID:     2,
			Symbol: "BTC",
		},
		BaseVariety: &models_variety.Variety{
			ID:     1,
			Symbol: "USDT",
		},
	}, nil).AnyTimes()

	// suite.assetRepo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", models_types.Amount("10000"))
	// suite.assetRepo.Despoit(suite.ctx, uuid.New().String(), "user1", "USDT", models_types.Amount("10000"))

}

func (suite *orderRepoTest) TestCreateLimitOrder() {
	suite.initMockData()
	order, err := suite.repo.CreateLimit(suite.ctx, "user1", "BTCUSDT", matching_types.OrderSideBuy, "10", "1")
	suite.Require().NoError(err)
	suite.Require().NotNil(order)

	//检查冻结的资产
	assetFreezes, err := suite.assetRepo.QueryFreeze(suite.ctx, map[string]any{
		"trans_id": map[string]any{
			"eq": order.OrderId,
		},
	})
	suite.Require().NoError(err)
	suite.Require().Equal(1, len(assetFreezes))
	suite.Require().Equal(models_types.Amount("10.1").Cmp(models_types.Amount(assetFreezes[0].FreezeAmount)), 0)
}
