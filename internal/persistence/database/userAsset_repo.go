package database

import (
	"errors"

	"github.com/shopspring/decimal"
	models "github.com/yzimhao/trading_engine/v2/internal/models/asset"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userAssetRepo struct {
	db     *gorm.DB
	logger *zap.Logger
}

var _ persistence.UserAssetRepository = (*userAssetRepo)(nil)

func NewUserAssetRepo(datasource *gorm.DB, logger *zap.Logger) persistence.UserAssetRepository {

	return &userAssetRepo{
		db:     datasource,
		logger: logger,
	}

}

func (u *userAssetRepo) QueryUserAsset(userId string, symbol string) (*entities.UserAsset, error) {
	var asset entities.UserAsset

	if err := u.db.Where("user_id = ? AND symbol = ?", userId, symbol).First(&asset).Error; err != nil {
		return nil, err
	}
	return &asset, nil
}

func (u *userAssetRepo) QueryUserAssets(userId string, symbols ...string) ([]*entities.UserAsset, error) {
	var assets []*entities.UserAsset
	if err := u.db.Where("user_id = ?", userId).Where("symbol in (?)", symbols).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (r *userAssetRepo) Despoit(transId, userId, symbol string, amount decimal.Decimal) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(tx, symbol, entities.SYSTEM_USER_ROOT, userId, amount, transId)
	})
}

func (r *userAssetRepo) Withdraw(transId, userId, symbol string, amount decimal.Decimal) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(tx, symbol, userId, entities.SYSTEM_USER_ROOT, amount, transId)
	})
}

// 两个user之间的转账
func (r *userAssetRepo) Transfer(transId, from, to, symbol string, amount decimal.Decimal) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.transfer(tx, symbol, from, to, amount, transId)
	})
}

// 冻结资产
// 这里使用tx传入，方便在结算的时候事务中使用
func (r *userAssetRepo) Freeze(tx *gorm.DB, transId, userId, symbol string, amount decimal.Decimal) (*entities.UserAssetFreeze, error) {
	if amount.Cmp(decimal.Zero) < 0 {
		return nil, errors.New("amount must be >= 0")
	}

	asset := entities.UserAsset{UserId: userId, Symbol: symbol}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", userId, symbol).FirstOrCreate(&asset).Error; err != nil {
		return nil, err
	}

	//冻结金额为0，冻结全部可用
	if amount.Cmp(decimal.Zero) == 0 {
		amount = asset.AvailBalance
	}

	asset.AvailBalance = asset.AvailBalance.Sub(amount)
	asset.FreezeBalance = asset.FreezeBalance.Add(amount)

	if asset.AvailBalance.Cmp(decimal.Zero) < 0 {
		return nil, errors.New("insufficient balance")
	}

	if tx.Where("user_id = ? AND symbol = ?", userId, symbol).Updates(&asset).Error != nil {
		return nil, errors.New("update asset failed")
	}

	//freeze log
	freezeLog := &entities.UserAssetFreeze{
		UserId:       userId,
		Symbol:       symbol,
		Amount:       amount,
		FreezeAmount: amount,
		TransId:      transId,
		// FreezeType:   entities.FreezeTypeTrade, //TODO 冻结类型
	}
	if tx.Create(&freezeLog).Error != nil {
		return nil, errors.New("create freeze log failed")
	}

	return freezeLog, nil
}

// 解冻资产
// amount为0，则解冻这条记录的全部剩余
func (r *userAssetRepo) UnFreeze(tx *gorm.DB, transId, userId, symbol string, amount decimal.Decimal) error {
	if amount.Cmp(decimal.Zero) < 0 {
		return errors.New("amount must be > 0")
	}

	freeze := entities.UserAssetFreeze{UserId: userId, Symbol: symbol, TransId: transId}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ? AND trans_id = ?", userId, symbol, transId).First(&freeze).Error; err != nil {
		return err
	}

	if freeze.Status == entities.FreezeStatusDone {
		return errors.New("unfreeze already done")
	}

	//解冻金额为0，则全部金额解冻
	if amount.Cmp(decimal.Zero) == 0 {
		amount = freeze.FreezeAmount
	}

	freeze.FreezeAmount = freeze.FreezeAmount.Sub(amount)
	if freeze.FreezeAmount.Cmp(decimal.Zero) == 0 {
		freeze.Status = entities.FreezeStatusDone
	}

	if tx.Where("user_id = ? AND symbol = ? AND trans_id = ?", userId, symbol, transId).Updates(&freeze).Error != nil {
		return errors.New("update freeze failed")
	}
	// 冻结资产变可用资产
	asset := entities.UserAsset{UserId: userId, Symbol: symbol}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", userId, symbol).FirstOrCreate(&asset).Error; err != nil {
		return err
	}
	asset.FreezeBalance = asset.FreezeBalance.Sub(amount)
	asset.AvailBalance = asset.AvailBalance.Add(amount)
	if tx.Where("user_id = ? AND symbol = ?", userId, symbol).Updates(&asset).Error != nil {
		return errors.New("update asset failed")
	}

	return nil
}

func (r *userAssetRepo) QueryFreeze(filter map[string]any) (assetFreezes []*models.AssetFreeze, err error) {
	// query := &datasource_types.PageQuery{
	// 	Filter: filter,
	// }
	// data, err := r.assetFreezeRepo.Query(ctx, query)
	// return data, err
	return nil, nil
}

func (r *userAssetRepo) TransferWithTx(tx *gorm.DB, transId, from, to, symbol string, amount decimal.Decimal) error {
	return r.transfer(tx, symbol, from, to, amount, transId)
}

func (r *userAssetRepo) transfer(tx *gorm.DB, symbol, from, to string, amount decimal.Decimal, transId string) error {
	if from == to {
		return errors.New("from and to cannot be the same")
	}

	if amount.Cmp(decimal.Zero) <= 0 {
		return errors.New("amount must be greater than 0")
	}

	//TODO transId去重

	fromAsset := entities.UserAsset{UserId: from, Symbol: symbol}
	//TODO tx.Clauses(clause.Locking{Strength: "FOR UPDATE"})
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", from, symbol).FirstOrCreate(&fromAsset).Error; err != nil {
		return err
	}

	toAsset := entities.UserAsset{UserId: to, Symbol: symbol}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ? AND symbol = ?", to, symbol).FirstOrCreate(&toAsset).Error; err != nil {
		return err
	}

	fromAsset.TotalBalance = fromAsset.TotalBalance.Sub(amount)
	fromAsset.AvailBalance = fromAsset.AvailBalance.Sub(amount)

	if fromAsset.UserId != entities.SYSTEM_USER_ROOT {
		if fromAsset.AvailBalance.Cmp(decimal.Zero) < 0 {
			return errors.New("insufficient balance")
		}
	}

	if tx.Where("user_id = ? AND symbol = ?", from, symbol).Updates(&fromAsset).Error != nil {
		return errors.New("update from asset failed")
	}

	toAsset.TotalBalance = toAsset.TotalBalance.Add(amount)
	toAsset.AvailBalance = toAsset.AvailBalance.Add(amount)
	if tx.Where("user_id = ? AND symbol = ?", to, symbol).Updates(&toAsset).Error != nil {
		return errors.New("update to asset failed")
	}

	fromLog := &entities.UserAssetLog{
		UserId:        from,
		Symbol:        symbol,
		BeforeBalance: fromAsset.TotalBalance.Add(amount),
		Amount:        amount.Neg(),
		AfterBalance:  fromAsset.TotalBalance,
		TransID:       transId,
		ChangeType:    entities.AssetChangeTypeTransfer,
	}
	if tx.Create(&fromLog).Error != nil {
		return errors.New("create from asset log failed")
	}

	toLog := &entities.UserAssetLog{
		UserId:        to,
		Symbol:        symbol,
		BeforeBalance: toAsset.TotalBalance.Sub(amount),
		Amount:        amount,
		AfterBalance:  toAsset.TotalBalance,
		TransID:       transId,
		ChangeType:    entities.AssetChangeTypeTransfer,
	}
	if tx.Create(&toLog).Error != nil {
		return errors.New("create to asset log failed")
	}

	return nil
}
