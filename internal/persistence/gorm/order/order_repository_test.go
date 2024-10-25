package order_test

import (
	"context"
	"testing"

	"github.com/duolacloud/crud-core/datasource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm"
	gorm_order "github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/order"
	mock_trade_variety "github.com/yzimhao/trading_engine/v2/mocks/persistence/trade_variety"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	_gorm "gorm.io/gorm"
)

type orderRepoTest struct {
	suite.Suite
	ctx    context.Context
	ctrl   *gomock.Controller
	repo   persistence.OrderRepository
	v      *viper.Viper
	gorm   *_gorm.DB
	logger *zap.Logger
}

func (suite *orderRepoTest) SetupTest() {
	_ = gotenv.Load("../../../.env")

	suite.ctx = context.Background()
	suite.ctrl = gomock.NewController(suite.T())
	suite.v = di.NewViper()
	suite.gorm = di.NewGorm(suite.v)
	suite.logger = zap.NewNop()
	redis := di.NewRedis(suite.v, suite.logger)
	cache, _ := di.NewCache(suite.v, redis)
	logger := zap.NewNop()
	mockTradeVarietyRepo := mock_trade_variety.NewMockTradeVarietyRepository(suite.ctrl)
	assetRepo := gorm.NewAssetRepo(datasource.NewDataSource(suite.gorm), cache, logger)
	suite.repo = gorm_order.NewOrderRepo(suite.gorm, logger, mockTradeVarietyRepo, assetRepo)
}

func TestOrderRepo(t *testing.T) {
	suite.Run(t, new(orderRepoTest))
}

func (suite *orderRepoTest) TearDownTest() {
	// migrations.MigrateDown(suite.gorm, suite.v, suite.logger)
}

func (suite *orderRepoTest) TestCreateOrder() {}
