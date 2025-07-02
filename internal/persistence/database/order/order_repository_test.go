package order_test

// import (
// 	"context"
// 	"testing"

// 	"github.com/duolacloud/crud-core/datasource"
// 	"github.com/google/uuid"
// 	"github.com/spf13/viper"
// 	"github.com/stretchr/testify/suite"
// 	"github.com/subosito/gotenv"
// 	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
// 	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
// 	models_variety "github.com/yzimhao/trading_engine/v2/internal/models/variety"
// 	"github.com/yzimhao/trading_engine/v2/internal/persistence"
// 	gorm_asset "github.com/yzimhao/trading_engine/v2/internal/persistence/database"
// 	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
// 	gorm_order "github.com/yzimhao/trading_engine/v2/internal/persistence/database/order"
// 	mock_asset "github.com/yzimhao/trading_engine/v2/mocks/persistence/asset"
// 	mock_trade_variety "github.com/yzimhao/trading_engine/v2/mocks/persistence/trade_variety"
// 	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
// 	"go.uber.org/mock/gomock"
// 	"go.uber.org/zap"
// 	_gorm "gorm.io/gorm"
// )

// var (
// 	testUser         = "1"
// 	testSymbol       = "BTCUSDT"
// 	testBaseSymbol   = "USDT"
// 	testTargetSymbol = "BTC"
// 	initBalance      = models_types.Numeric("10000")
// )

// type orderRepoTest struct {
// 	suite.Suite
// 	ctx                  context.Context
// 	ctrl                 *gomock.Controller
// 	repo                 persistence.OrderRepository
// 	v                    *viper.Viper
// 	gorm                 *_gorm.DB
// 	logger               *zap.Logger
// 	mockTradeVarietyRepo *mock_trade_variety.MockTradeVarietyRepository
// 	mockAssetRepo        *mock_asset.MockAssetRepository
// 	assetRepo            persistence.AssetRepository
// }

// func (suite *orderRepoTest) SetupTest() {
// 	_ = gotenv.Load("../../../../.env")

// 	suite.ctx = context.Background()
// 	suite.ctrl = gomock.NewController(suite.T())
// 	suite.v = provider.NewViper()
// 	suite.gorm = provider.NewGorm(suite.v)
// 	suite.logger = zap.NewExample()
// 	redis := provider.NewRedis(suite.v, suite.logger)
// 	cache, _ := provider.NewCache(suite.v, redis)
// 	logger := zap.NewNop()
// 	mockTradeVarietyRepo := mock_trade_variety.NewMockTradeVarietyRepository(suite.ctrl)
// 	mockAssetRepo := mock_asset.NewMockAssetRepository(suite.ctrl)
// 	suite.mockTradeVarietyRepo = mockTradeVarietyRepo
// 	suite.mockAssetRepo = mockAssetRepo
// 	suite.assetRepo = gorm_asset.NewAssetRepo(datasource.NewDataSource(suite.gorm), cache, logger)
// 	suite.repo = gorm_order.NewOrderRepo(suite.gorm, logger, suite.mockTradeVarietyRepo, suite.assetRepo)
// }

// func TestOrderRepo(t *testing.T) {
// 	suite.Run(t, new(orderRepoTest))
// }

// func (suite *orderRepoTest) TearDownTest() {
// 	// migrations.MigrateDown(suite.gorm, suite.v, suite.logger)
// 	suite.ctrl.Finish()
// }

// func (suite *orderRepoTest) initMockData() {
// 	tables := []any{
// 		&entities.UserAsset{},
// 		&entities.UserAssetFreeze{},
// 		&entities.UserAssetLog{},
// 	}

// 	for _, table := range tables {
// 		err := suite.gorm.Migrator().CreateTable(table)
// 		suite.Require().NoError(err)
// 	}

// 	suite.mockTradeVarietyRepo.EXPECT().FindBySymbol(suite.ctx, "BTCUSDT").Return(&models_variety.TradeVariety{
// 		ID:            1,
// 		Symbol:        "BTCUSDT",
// 		Name:          "BTCUSDT",
// 		TargetId:      2,
// 		BaseId:        1,
// 		PriceDecimals: 2,
// 		QtyDecimals:   4,
// 		FeeRate:       "0.01",
// 		TargetVariety: &models_variety.Variety{
// 			ID:     2,
// 			Symbol: "BTC",
// 		},
// 		BaseVariety: &models_variety.Variety{
// 			ID:     1,
// 			Symbol: "USDT",
// 		},
// 	}, nil).AnyTimes()

// 	suite.assetRepo.Despoit(uuid.New().String(), testUser, testTargetSymbol, initBalance)
// 	suite.assetRepo.Despoit(uuid.New().String(), testUser, testBaseSymbol, initBalance)

// }

// func (suite *orderRepoTest) cleanMockData() {

// 	tables := []any{
// 		&entities.Asset{},
// 		&entities.UserAssetFreeze{},
// 		&entities.UserAssetLog{},
// 		&entities.Order{Symbol: testSymbol},
// 		&entities.UnfinishedOrder{},
// 	}

// 	for _, table := range tables {
// 		indexes, err := suite.gorm.Migrator().GetIndexes(table)
// 		suite.Require().NoError(err)
// 		for _, index := range indexes {
// 			err := suite.gorm.Migrator().DropIndex(table, index.Name())
// 			suite.Require().NoError(err)
// 		}
// 		err = suite.gorm.Migrator().DropTable(table)
// 		suite.Require().NoError(err)
// 	}
// }

// func (suite *orderRepoTest) TestCreateLimitOrderSideBuy() {
// 	suite.cleanMockData()
// 	suite.initMockData()

// 	order, err := suite.repo.CreateLimit(suite.ctx, testUser, testSymbol, matching_types.OrderSideBuy, "10", "1")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(order)

// 	//检查冻结的资产
// 	assetFreezes, err := suite.assetRepo.QueryFreeze(suite.ctx, map[string]any{
// 		"trans_id": map[string]any{
// 			"eq": order.OrderId,
// 		},
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(1, len(assetFreezes))
// 	suite.Require().Equal(models_types.Numeric("10.1").Cmp(models_types.Numeric(assetFreezes[0].FreezeAmount)), 0)
// }

// func (suite *orderRepoTest) TestCreateLimitOrderSideSell() {
// 	suite.cleanMockData()
// 	suite.initMockData()

// 	order, err := suite.repo.CreateLimit(suite.ctx, testUser, testSymbol, matching_types.OrderSideSell, "10", "1")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(order)

// 	//检查冻结的资产
// 	assetFreezes, err := suite.assetRepo.QueryFreeze(suite.ctx, map[string]any{
// 		"trans_id": map[string]any{
// 			"eq": order.OrderId,
// 		},
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(1, len(assetFreezes))
// 	suite.Require().Equal(models_types.Numeric("1").Cmp(models_types.Numeric(assetFreezes[0].FreezeAmount)), 0)
// }

// func (suite *orderRepoTest) TestCreateMarketOrderSideBuy_Qty() {
// 	suite.cleanMockData()
// 	suite.initMockData()

// 	order, err := suite.repo.CreateMarketByQty(suite.ctx, testUser, testSymbol, matching_types.OrderSideBuy, "1")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(order)

// 	//检查冻结的资产
// 	assetFreezes, err := suite.assetRepo.QueryFreeze(suite.ctx, map[string]any{
// 		"trans_id": map[string]any{
// 			"eq": order.OrderId,
// 		},
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(1, len(assetFreezes))
// 	suite.Require().Equal(initBalance.Cmp(models_types.Numeric(assetFreezes[0].FreezeAmount)), 0)
// 	suite.Require().Equal(assetFreezes[0].FreezeAmount, order.FreezeAmount)
// }

// func (suite *orderRepoTest) TestCreateMarketOrderSideBuy_Amount() {
// 	suite.cleanMockData()
// 	suite.initMockData()

// 	order, err := suite.repo.CreateMarketByAmount(suite.ctx, testUser, testSymbol, matching_types.OrderSideBuy, "1000")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(order)

// 	//检查冻结的资产
// 	assetFreezes, err := suite.assetRepo.QueryFreeze(suite.ctx, map[string]any{
// 		"trans_id": map[string]any{
// 			"eq": order.OrderId,
// 		},
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(1, len(assetFreezes))
// 	suite.Require().Equal(0, models_types.Numeric("1000").Cmp(models_types.Numeric(assetFreezes[0].FreezeAmount)))
// 	suite.Require().Equal(0, models_types.Numeric(assetFreezes[0].FreezeAmount).Cmp(models_types.Numeric(order.FreezeAmount)))
// }

// func (suite *orderRepoTest) TestCreateMarketOrderSideSell_Qty() {
// 	suite.cleanMockData()
// 	suite.initMockData()

// 	order, err := suite.repo.CreateMarketByQty(suite.ctx, testUser, testSymbol, matching_types.OrderSideSell, "1")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(order)

// 	//检查冻结的资产
// 	assetFreezes, err := suite.assetRepo.QueryFreeze(suite.ctx, map[string]any{
// 		"trans_id": map[string]any{
// 			"eq": order.OrderId,
// 		},
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(1, len(assetFreezes))
// 	suite.Require().Equal(0, models_types.Numeric("1").Cmp(models_types.Numeric(assetFreezes[0].FreezeAmount)))
// 	suite.Require().Equal(0, models_types.Numeric(assetFreezes[0].FreezeAmount).Cmp(models_types.Numeric(order.FreezeQty)))
// }

// func (suite *orderRepoTest) TestCreateMarketOrderSideSell_Amount() {
// 	suite.cleanMockData()
// 	suite.initMockData()

// 	order, err := suite.repo.CreateMarketByAmount(suite.ctx, testUser, testSymbol, matching_types.OrderSideSell, "1000")
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(order)

// 	//检查冻结的资产
// 	assetFreezes, err := suite.assetRepo.QueryFreeze(suite.ctx, map[string]any{
// 		"trans_id": map[string]any{
// 			"eq": order.OrderId,
// 		},
// 	})
// 	suite.Require().NoError(err)
// 	suite.Require().Equal(1, len(assetFreezes))
// 	suite.Require().Equal(0, initBalance.Cmp(models_types.Numeric(assetFreezes[0].FreezeAmount)))
// 	suite.Require().Equal(0, models_types.Numeric(assetFreezes[0].FreezeAmount).Cmp(models_types.Numeric(order.FreezeQty)))
// }
