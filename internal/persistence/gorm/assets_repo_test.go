package gorm_test

import (
	"context"
	"testing"

	"github.com/duolacloud/crud-core/datasource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm"
	"github.com/yzimhao/trading_engine/v2/migrations"
	"go.uber.org/zap"
	_gorm "gorm.io/gorm"
)

type assetsRepoTest struct {
	suite.Suite
	ctx    context.Context
	repo   persistence.AssetsRepository
	v      *viper.Viper
	gorm   *_gorm.DB
	logger *zap.Logger
}

func (suite *assetsRepoTest) SetupTest() {
	suite.ctx = context.Background()

	suite.v = di.NewViper()
	suite.gorm = di.NewGorm(suite.v)
	suite.logger = zap.NewNop()
	redis := di.NewRedis(suite.v, suite.logger)
	cache, _ := di.NewCache(suite.v, redis)
	suite.repo = gorm.NewAssetsRepo(datasource.NewDataSource(suite.gorm), cache)
}

func TestAssetsRepo(t *testing.T) {
	suite.Run(t, new(assetsRepoTest))
}

func (suite *assetsRepoTest) TearDownTest() {
	migrations.MigrateDown(suite.gorm, suite.v, suite.logger)
}

func (suite *assetsRepoTest) TestDespoit() {
	migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

	_, err := suite.repo.Despoit(suite.ctx, "user1", "BTC", "1")
	suite.NoError(err)

	asset, err := suite.repo.FindOne(suite.ctx, "user1", "BTC")
	suite.NoError(err)
	suite.Equal(asset.UserId, "user1")
	suite.Equal(asset.Symbol, "BTC")
	suite.Equal(asset.TotalBalance.Cmp(types.Amount("1")), 0)
	suite.Equal(asset.AvailBalance.Cmp(types.Amount("1")), 0)
	suite.Equal(asset.FreezeBalance.Cmp(types.Amount("0")), 0)

	//TODO test aseets_log
}

func (suite *assetsRepoTest) TestWithdraw() {

	testCases := []struct {
		name  string
		setup func()
	}{
		{
			name: "提现用户不存在",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				_, err := suite.repo.Withdraw(suite.ctx, "user1", "BTC", "1000")
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现用户余额不足",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				_, err := suite.repo.Despoit(suite.ctx, "user1", "BTC", "1")
				suite.NoError(err)

				_, err = suite.repo.Withdraw(suite.ctx, "user1", "BTC", "1000")
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现 余额充足",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				_, err := suite.repo.Despoit(suite.ctx, "user1", "BTC", "2000")
				suite.NoError(err)

				_, err = suite.repo.Withdraw(suite.ctx, "user1", "BTC", "1000")
				suite.NoError(err)

				asset, err := suite.repo.FindOne(suite.ctx, "user1", "BTC")
				suite.NoError(err)
				suite.Equal(asset.UserId, "user1")
				suite.Equal(asset.Symbol, "BTC")
				suite.Equal(asset.TotalBalance.Cmp(types.Amount("1000")), 0)
				suite.Equal(asset.AvailBalance.Cmp(types.Amount("1000")), 0)
				suite.Equal(asset.FreezeBalance.Cmp(types.Amount("0")), 0)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.setup()
		})
	}
}
