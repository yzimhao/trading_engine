package database_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"github.com/subosito/gotenv"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"github.com/yzimhao/trading_engine/v2/migrations"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type assetsRepoTest struct {
	suite.Suite
	ctx             context.Context
	userAssetrepo   persistence.UserAssetRepository
	v               *viper.Viper
	db              *gorm.DB
	logger          *zap.Logger
	testAssetSymbol string
	testUserId      string
}

func (suite *assetsRepoTest) SetupTest() {
	_ = gotenv.Load(provider.Root() + "/.env")

	suite.ctx = context.Background()
	logger := zap.NewExample()
	suite.v = provider.NewViper(logger)
	suite.db = provider.NewGorm(suite.v)
	suite.logger = logger
	suite.testAssetSymbol = "BTC"
	suite.testUserId = "testuser"

	suite.userAssetrepo = database.NewUserAssetRepo(suite.db, logger)
}

func TestAssetsRepo(t *testing.T) {
	suite.Run(t, new(assetsRepoTest))
}

func (suite *assetsRepoTest) TearDownTest() {
	// migrations.MigrateDown(suite.db, suite.v, suite.logger)
}

func (suite *assetsRepoTest) TestDespoit() {
	migrations.MigrateUp(suite.db, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.db, suite.v, suite.logger)

	err := suite.userAssetrepo.Despoit(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1"))
	suite.NoError(err)

	asset, err := suite.userAssetrepo.QueryUserAsset(suite.testUserId, suite.testAssetSymbol)
	suite.NoError(err)
	suite.Equal(suite.testUserId, asset.UserId)
	suite.Equal(suite.testAssetSymbol, asset.Symbol)
	suite.Equal(0, asset.TotalBalance.Cmp(d("1")))
	suite.Equal(0, asset.AvailBalance.Cmp(d("1")))
	suite.Equal(0, asset.FreezeBalance.Cmp(d("0")))

	systemAsset, err := suite.userAssetrepo.QueryUserAsset(entities.SYSTEM_USER_ROOT, suite.testAssetSymbol)
	suite.NoError(err)
	suite.Equal(entities.SYSTEM_USER_ROOT, systemAsset.UserId)
	suite.Equal(suite.testAssetSymbol, systemAsset.Symbol)
	suite.Equal(0, systemAsset.TotalBalance.Cmp(d("-1")))
	suite.Equal(0, systemAsset.AvailBalance.Cmp(d("-1")))
	suite.Equal(0, systemAsset.FreezeBalance.Cmp(d("0")))
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
				migrations.MigrateUp(suite.db, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.db, suite.v, suite.logger)

				err := suite.userAssetrepo.Withdraw(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1000"))
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现用户余额不足",
			setup: func() {
				migrations.MigrateUp(suite.db, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.db, suite.v, suite.logger)

				err := suite.userAssetrepo.Despoit(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1"))
				suite.NoError(err)

				err = suite.userAssetrepo.Withdraw(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1000"))
				suite.Equal(err.Error(), "insufficient balance")
			},
		},
		{
			name: "提现 余额充足",
			setup: func() {
				migrations.MigrateUp(suite.db, suite.v, suite.logger)
				defer migrations.MigrateDown(suite.db, suite.v, suite.logger)

				err := suite.userAssetrepo.Despoit(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("2000"))
				suite.NoError(err)

				err = suite.userAssetrepo.Withdraw(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1000"))
				suite.NoError(err)

				asset, err := suite.userAssetrepo.QueryUserAsset(suite.testUserId, suite.testAssetSymbol)
				suite.NoError(err)
				suite.Equal(suite.testUserId, asset.UserId)
				suite.Equal(suite.testAssetSymbol, asset.Symbol)
				suite.Equal(0, asset.TotalBalance.Cmp(d("1000")))
				suite.Equal(0, asset.AvailBalance.Cmp(d("1000")))
				suite.Equal(0, asset.FreezeBalance.Cmp(d("0")))
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
	migrations.MigrateUp(suite.db, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.db, suite.v, suite.logger)

	_, err := suite.userAssetrepo.Freeze(suite.db, uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1000"))
	suite.Equal(err.Error(), "insufficient balance")

	err = suite.userAssetrepo.Despoit(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1000"))
	suite.NoError(err)

	_, err = suite.userAssetrepo.Freeze(suite.db, uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1"))
	suite.NoError(err)

	asset, err := suite.userAssetrepo.QueryUserAsset(suite.testUserId, suite.testAssetSymbol)
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(d("1")))
	suite.Equal(0, asset.AvailBalance.Cmp(d("999")))

	// 冻结全部
	_, err = suite.userAssetrepo.Freeze(suite.db, uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("0"))
	suite.NoError(err)

	asset, err = suite.userAssetrepo.QueryUserAsset(suite.testUserId, suite.testAssetSymbol)
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(d("1000")))
	suite.Equal(0, asset.AvailBalance.Cmp(d("0")))
}

func (suite *assetsRepoTest) TestTransfer() {
	migrations.MigrateUp(suite.db, suite.v, suite.logger)
	defer migrations.MigrateDown(suite.db, suite.v, suite.logger)

	err := suite.userAssetrepo.Despoit(uuid.New().String(), suite.testUserId, suite.testAssetSymbol, d("1000"))
	suite.NoError(err)

	transId := uuid.New().String()
	_, err = suite.userAssetrepo.Freeze(suite.db, transId, suite.testUserId, suite.testAssetSymbol, d("900"))
	suite.NoError(err)

	err = suite.userAssetrepo.UnFreeze(suite.db, transId, suite.testUserId, suite.testAssetSymbol, d("1"))
	suite.NoError(err)

	asset, err := suite.userAssetrepo.QueryUserAsset(suite.testUserId, suite.testAssetSymbol)
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(d("899")))
	suite.Equal(0, asset.AvailBalance.Cmp(d("101")))

	//解冻全部
	err = suite.userAssetrepo.UnFreeze(suite.db, transId, suite.testUserId, suite.testAssetSymbol, d("0"))
	suite.NoError(err)

	asset, err = suite.userAssetrepo.QueryUserAsset(suite.testUserId, suite.testAssetSymbol)
	suite.NoError(err)
	suite.Equal(0, asset.FreezeBalance.Cmp(d("0")))
	suite.Equal(0, asset.AvailBalance.Cmp(d("1000")))
}

func d(f string) decimal.Decimal {
	d, _ := decimal.NewFromString(f)
	return d
}
