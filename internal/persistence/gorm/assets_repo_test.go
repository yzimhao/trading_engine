package gorm_test

import (
	"context"
	"testing"

	"github.com/duolacloud/crud-core/datasource"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"github.com/yzimhao/trading_engine/v2/migrations"
	"go.uber.org/zap"
	_gorm "gorm.io/gorm"
)

type assetsRepoTest struct {
	suite.Suite
	ctx    context.Context
	repo   persistence.AssetRepository
	v      *viper.Viper
	gorm   *_gorm.DB
	logger *zap.Logger
}

func (suite *assetsRepoTest) SetupTest() {
	_ = gotenv.Load("../../../.env")

	suite.ctx = context.Background()

	suite.v = di.NewViper()
	suite.gorm = di.NewGorm(suite.v)
	suite.logger = zap.NewNop()
	redis := di.NewRedis(suite.v, suite.logger)
	cache, _ := di.NewCache(suite.v, redis)
	logger := zap.NewNop()
	suite.repo = gorm.NewAssetRepo(datasource.NewDataSource(suite.gorm), cache, logger)
}

func TestAssetsRepo(t *testing.T) {
	suite.Run(t, new(assetsRepoTest))
}

func (suite *assetsRepoTest) TearDownTest() {
	// migrations.MigrateDown(suite.gorm, suite.v, suite.logger)
}

func (suite *assetsRepoTest) TestDespoit() {
	migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

	err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("1"))
	suite.NoError(err)

	asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(asset.UserId, "user1")
	suite.Equal(asset.Symbol, "BTC")
	suite.Equal(asset.TotalBalance.Cmp(types.Amount("1")), 0)
	suite.Equal(asset.AvailBalance.Cmp(types.Amount("1")), 0)
	suite.Equal(asset.FreezeBalance.Cmp(types.Amount("0")), 0)

	systemAsset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": entities.SYSTEM_USER_ID,
		},
	})
	suite.NoError(err)
	suite.Equal(systemAsset.UserId, entities.SYSTEM_USER_ID)
	suite.Equal(systemAsset.Symbol, "BTC")
	suite.Equal(systemAsset.TotalBalance.Cmp(types.Amount("-1")), 0)
	suite.Equal(systemAsset.AvailBalance.Cmp(types.Amount("-1")), 0)
	suite.Equal(systemAsset.FreezeBalance.Cmp(types.Amount("0")), 0)
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

				err := suite.repo.Withdraw(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("1000"))
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现用户余额不足",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("1"))
				suite.NoError(err)

				err = suite.repo.Withdraw(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("1000"))
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现 余额充足",
			setup: func() {
				migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

				err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("2000"))
				suite.NoError(err)

				err = suite.repo.Withdraw(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("1000"))
				suite.NoError(err)

				asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
					"symbol": map[string]any{
						"eq": "BTC",
					},
					"user_id": map[string]any{
						"eq": "user1",
					},
				})
				// fmt.Printf("asset: %+v\n", asset)
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

func (suite *assetsRepoTest) TestFreeze() {
	migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

	_, err := suite.repo.Freeze(suite.ctx, suite.gorm, uuid.New().String(), "user1", "BTC", types.Amount("1000"))
	suite.Equal(err.Error(), "insufficient balance")

	err = suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("1000"))
	suite.NoError(err)

	_, err = suite.repo.Freeze(suite.ctx, suite.gorm, uuid.New().String(), "user1", "BTC", types.Amount("1"))
	suite.NoError(err)

	asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(asset.FreezeBalance.Cmp(types.Amount("1")), 0)
	suite.Equal(asset.AvailBalance.Cmp(types.Amount("999")), 0)

	// 冻结全部
	_, err = suite.repo.Freeze(suite.ctx, suite.gorm, uuid.New().String(), "user1", "BTC", types.Amount("0"))
	suite.NoError(err)

	asset, err = suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(asset.FreezeBalance.Cmp(types.Amount("1000")), 0)
	suite.Equal(asset.AvailBalance.Cmp(types.Amount("0")), 0)
}

func (suite *assetsRepoTest) TestTransfer() {
	migrations.MigrateUp(suite.gorm, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.gorm, suite.v, suite.logger)

	err := suite.repo.Despoit(suite.ctx, uuid.New().String(), "user1", "BTC", types.Amount("1000"))
	suite.NoError(err)

	transId := uuid.New().String()
	_, err = suite.repo.Freeze(suite.ctx, suite.gorm, transId, "user1", "BTC", types.Amount("900"))
	suite.NoError(err)

	err = suite.repo.UnFreeze(suite.ctx, suite.gorm, transId, "user1", "BTC", types.Amount("1"))
	suite.NoError(err)

	asset, err := suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(asset.FreezeBalance.Cmp(types.Amount("899")), 0)
	suite.Equal(asset.AvailBalance.Cmp(types.Amount("101")), 0)

	//解冻全部
	err = suite.repo.UnFreeze(suite.ctx, suite.gorm, transId, "user1", "BTC", types.Amount("0"))
	suite.NoError(err)

	asset, err = suite.repo.QueryOne(suite.ctx, map[string]any{
		"symbol": map[string]any{
			"eq": "BTC",
		},
		"user_id": map[string]any{
			"eq": "user1",
		},
	})
	suite.NoError(err)
	suite.Equal(asset.FreezeBalance.Cmp(types.Amount("0")), 0)
	suite.Equal(asset.AvailBalance.Cmp(types.Amount("1000")), 0)
}
