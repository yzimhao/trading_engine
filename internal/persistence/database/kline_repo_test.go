package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/duolacloud/crud-core/datasource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	kline_types "github.com/yzimhao/trading_engine/v2/pkg/kline/types"
	"go.uber.org/zap"
	_gorm "gorm.io/gorm"
)

var (
	testSymbol = "BTCUSDT"
	testPeriod = kline_types.PERIOD_M1
)

type klineRepoTest struct {
	suite.Suite
	ctx    context.Context
	repo   persistence.KlineRepository
	v      *viper.Viper
	gorm   *_gorm.DB
	logger *zap.Logger
}

func TestKlineRepo(t *testing.T) {
	suite.Run(t, new(klineRepoTest))
}

func (suite *klineRepoTest) SetupTest() {
	_ = gotenv.Load(provider.Root() + "/.env")

	suite.ctx = context.Background()

	logger := zap.NewExample()
	suite.v = provider.NewViper(logger)
	suite.gorm = provider.NewGorm(suite.v)
	suite.logger = logger
	redis := provider.NewRedis(suite.v, suite.logger)
	cache, _ := provider.NewCache(suite.v, redis)

	suite.repo = database.NewKlineRepo(datasource.NewDataSource(suite.gorm), cache, logger)
}

func (suite *klineRepoTest) TearDownTest() {
}

func (suite *klineRepoTest) cleanKlineTable() {
	table := entities.Kline{
		Symbol: testSymbol,
		Period: testPeriod,
	}

	indexes, err := suite.gorm.Migrator().GetIndexes(table.TableName())
	suite.Require().NoError(err)
	for _, index := range indexes {
		err := suite.gorm.Migrator().DropIndex(table.TableName(), index.Name())
		suite.Require().NoError(err)
	}
	err = suite.gorm.Migrator().DropTable(table.TableName())
	suite.Require().NoError(err)
}

func (suite *klineRepoTest) TestSaveKline() {
	suite.cleanKlineTable()

	now := time.Now()
	openAt, closeAt := kline_types.ParsePeriodTime(now, testPeriod)

	err := suite.repo.Save(suite.ctx, &entities.Kline{
		Symbol:  testSymbol,
		Period:  testPeriod,
		OpenAt:  openAt,
		CloseAt: closeAt,
		Open:    "1",
		High:    "2",
		Low:     "0.5",
		Close:   "1.5",
		Volume:  "1000",
		Amount:  "10000",
	})
	suite.Require().NoError(err)

	err = suite.repo.Save(suite.ctx, &entities.Kline{
		Symbol:  testSymbol,
		Period:  testPeriod,
		OpenAt:  openAt,
		CloseAt: closeAt,
		Open:    "1",
		High:    "10",
		Low:     "0.01",
		Close:   "8",
		Volume:  "1001",
		Amount:  "10000",
	})
	suite.Require().NoError(err)
}
