package gorm

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	k_repo "github.com/duolacloud/crud-core-gorm/repositories"
	"github.com/duolacloud/crud-core/cache"
	"github.com/duolacloud/crud-core/datasource"
	b_mappers "github.com/duolacloud/crud-core/mappers"
	"github.com/duolacloud/crud-core/repositories"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	"gorm.io/gorm"
)

type gormAssetRepo struct {
	*repositories.MapperRepository[models.Asset, models.CreateAsset, models.UpdateAsset, entities.Asset, entities.Asset, map[string]any]
	datasource      datasource.DataSource[gorm.DB]
	assetLogRepo    *gormAssetLogRepo
	assetFreezeRepo *gormAssetFreezeRepo
}

type gormAssetLogRepo struct {
	*repositories.MapperRepository[models.AssetLog, models.CreateAssetLog, models.UpdateAssetLog, entities.AssetLog, entities.AssetLog, map[string]any]
}

type gormAssetFreezeRepo struct {
	*repositories.MapperRepository[models.AssetFreeze, models.CreateAssetFreeze, models.UpdateAssetFreeze, entities.AssetFreeze, entities.AssetFreeze, map[string]any]
}

func NewAssetRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) persistence.AssetRepository {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.Asset, entities.Asset, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.Asset, models.CreateAsset, models.UpdateAsset, entities.Asset, entities.Asset, map[string]any](),
	)

	return &gormAssetRepo{
		MapperRepository: mapperRepo,
		datasource:       datasource,
		assetLogRepo:     newAssetLogRepo(datasource, cache),
		assetFreezeRepo:  newAssetFreezeRepo(datasource, cache),
	}

}

func newAssetLogRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) *gormAssetLogRepo {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.AssetLog, entities.AssetLog, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.AssetLog, models.CreateAssetLog, models.UpdateAssetLog, entities.AssetLog, entities.AssetLog, map[string]any](),
	)

	return &gormAssetLogRepo{
		MapperRepository: mapperRepo,
	}
}

func newAssetFreezeRepo(datasource datasource.DataSource[gorm.DB], cache cache.Cache) *gormAssetFreezeRepo {
	cacheWrapperRepo := repositories.NewCacheRepository(
		k_repo.NewGormCrudRepository[entities.AssetFreeze, entities.AssetFreeze, map[string]any](datasource),
		cache,
	)

	mapperRepo := repositories.NewMapperRepository(
		cacheWrapperRepo,
		b_mappers.NewJSONMapper[models.AssetFreeze, models.CreateAssetFreeze, models.UpdateAssetFreeze, entities.AssetFreeze, entities.AssetFreeze, map[string]any](),
	)

	return &gormAssetFreezeRepo{
		MapperRepository: mapperRepo,
	}
}

func (r *gormAssetRepo) Despoit(ctx context.Context, userId, symbol string, amount string) (order_id string, err error) {
	order_id = uuid.New().String()
	err = r.transfer(ctx, symbol, entities.SYSTEM_USER_ID, userId, types.Amount(amount), order_id)
	return
}

func (r *gormAssetRepo) Withdraw(ctx context.Context, userId, symbol, amount string) (order_id string, err error) {
	order_id = uuid.New().String()
	err = r.transfer(ctx, symbol, userId, entities.SYSTEM_USER_ID, types.Amount(amount), order_id)
	return
}

func (r *gormAssetRepo) Transfer(ctx context.Context, from, to, symbol, amount string) error {
	return r.transfer(ctx, symbol, from, to, types.Amount(amount), uuid.New().String())
}

func (r *gormAssetRepo) transfer(ctx context.Context, symbol, from, to string, amount types.Amount, transId string) error {

	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return errors.Wrap(err, "get gorm db")
	}

	//临时有要求使用原生sql 不使用orm
	rawDb, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "get rawdb")
	}

	fromUser, err := r.queryOne(ctx, rawDb, from, symbol)
	if err != nil {
		return errors.Wrap(err, "query from user")
	}

	toUser, err := r.queryOne(ctx, rawDb, to, symbol)
	if err != nil {
		return errors.Wrap(err, "query to user")
	}

	tx, err := rawDb.Begin()
	if err != nil {
		return errors.Wrap(err, "begin tx")
	}

	fromUser.TotalBalance = fromUser.TotalBalance.Sub(amount)
	fromUser.AvailBalance = fromUser.AvailBalance.Sub(amount)

	if fromUser.UserId != entities.SYSTEM_USER_ID {
		if fromUser.AvailBalance.Cmp(types.Amount("0")) < 0 {
			return errors.New("insufficient balance")
		}
	}

	toUser.TotalBalance = toUser.TotalBalance.Add(amount)
	toUser.AvailBalance = toUser.AvailBalance.Add(amount)

	if err := r.update(ctx, tx, fromUser); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update from user")
	}
	if err := r.update(ctx, tx, toUser); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update to user")
	}

	//create logs
	err = r.createLog(ctx, tx, &entities.AssetLog{
		UserId:        from,
		Symbol:        symbol,
		BeforeBalance: fromUser.TotalBalance.Add(amount),
		Amount:        amount,
		AfterBalance:  fromUser.TotalBalance,
		TransID:       transId,
		ChangeType:    "despoit",
	})
	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "create from log")
	}

	err = r.createLog(ctx, tx, &entities.AssetLog{
		UserId:        to,
		Symbol:        symbol,
		BeforeBalance: toUser.TotalBalance.Sub(amount),
		Amount:        amount,
		AfterBalance:  toUser.TotalBalance,
		TransID:       transId,
		ChangeType:    "withdraw",
	})

	if err != nil {
		tx.Rollback()
		return errors.Wrap(err, "create to log")
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "commit tx")
	}

	return nil
}

func (r *gormAssetRepo) createLog(ctx context.Context, db *sql.Tx, log *entities.AssetLog) error {
	stmt, err := db.Prepare("insert into assets_logs (id, user_id, symbol, before_balance, amount, after_balance, trans_id, change_type, info, created_at, updated_at) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)")
	if err != nil {
		return errors.Wrap(err, "prepare insert user")
	}

	_, err = stmt.Exec(uuid.New().String(), log.UserId, log.Symbol, log.BeforeBalance, log.Amount, log.AfterBalance, log.TransID, log.ChangeType, log.Info, time.Now(), time.Now())
	if err != nil {
		return errors.Wrap(err, "exec insert user")
	}

	return nil
}

func (r *gormAssetRepo) FindOne(ctx context.Context, userId, symbol string) (*entities.Asset, error) {
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get gorm db")
	}

	rawDb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get raw db")
	}

	return r.queryOne(ctx, rawDb, userId, symbol)
}

func (r *gormAssetRepo) FindAssetHistory(ctx context.Context) ([]entities.AssetLog, error) {
	db, err := r.datasource.GetDB(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "get gorm db")
	}

	rawDb, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "get raw db")
	}

	rows, err := rawDb.QueryContext(ctx, "SELECT * FROM assets_logs limit 10")
	if err != nil {
		return nil, errors.Wrap(err, "query assets logs")
	}

	defer rows.Close()

	var logs []entities.AssetLog

	for rows.Next() {
		var log entities.AssetLog
		err := rows.Scan(&log.Id, &log.UserId, &log.Symbol, &log.BeforeBalance, &log.Amount, &log.AfterBalance, &log.TransID, &log.ChangeType, &log.Info, &log.CreatedAt, &log.UpdatedAt)
		if err != nil {
			return nil, errors.Wrap(err, "scan assets log")
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *gormAssetRepo) queryOne(ctx context.Context, rawDb *sql.DB, userId, symbol string) (*entities.Asset, error) {

	// 查询是否存在指定的资产记录
	row := rawDb.QueryRowContext(ctx, "SELECT * FROM assets WHERE user_id = $1 AND symbol = $2 LIMIT 1", userId, symbol)
	var user entities.Asset
	err := row.Scan(&user.Id, &user.UserId, &user.Symbol, &user.TotalBalance, &user.FreezeBalance, &user.AvailBalance, &user.CreatedAt, &user.UpdatedAt)

	// 如果出现数据库查询错误
	if err != nil && err != sql.ErrNoRows {
		return nil, errors.Wrap(err, "query user")
	}

	if err == sql.ErrNoRows {
		return &entities.Asset{
			UserId: userId,
			Symbol: symbol,
		}, nil
	}

	return &user, nil
}

func (r *gormAssetRepo) update(ctx context.Context, tx *sql.Tx, user *entities.Asset) error {
	// 查询是否存在指定的资产记录
	row := tx.QueryRowContext(ctx, "SELECT id FROM assets WHERE user_id = $1 AND symbol = $2 LIMIT 1", user.UserId, user.Symbol)
	var id string
	err := row.Scan(&id)

	// 如果出现数据库查询错误
	if err != nil && err != sql.ErrNoRows {
		return errors.Wrap(err, "query user")
	}

	if err == sql.ErrNoRows {
		// 如果记录不存在，执行插入操作
		_, err := tx.ExecContext(ctx, "INSERT INTO assets (id, user_id, symbol, total_balance, avail_balance, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			uuid.New().String(), user.UserId, user.Symbol, user.TotalBalance, user.AvailBalance, time.Now(), time.Now())
		if err != nil {
			return errors.Wrap(err, "exec insert user")
		}
	} else {
		// 如果记录存在，执行更新操作
		_, err := tx.ExecContext(ctx, "UPDATE assets SET total_balance = $1, avail_balance = $2, updated_at = $3 WHERE id = $4",
			user.TotalBalance, user.AvailBalance, time.Now(), id)
		if err != nil {
			return errors.Wrap(err, "exec update user")
		}
	}

	return nil
}
