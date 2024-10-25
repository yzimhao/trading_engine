package order

import (
	"context"
	"time"

	"github.com/pkg/errors"
	models_order "github.com/yzimhao/trading_engine/v2/internal/models/order"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/models/variety"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type orderRepository struct {
	db               *gorm.DB
	logger           *zap.Logger
	tradeVarietyRepo persistence.TradeVarietyRepository
	assetRepo        persistence.AssetRepository
}

var _ persistence.OrderRepository = (*orderRepository)(nil)

func NewOrderRepo(
	db *gorm.DB,
	logger *zap.Logger,
	tradeVarietyRepo persistence.TradeVarietyRepository,
	assetRepo persistence.AssetRepository,
) persistence.OrderRepository {
	return &orderRepository{
		db:               db,
		logger:           logger,
		tradeVarietyRepo: tradeVarietyRepo,
		assetRepo:        assetRepo,
	}
}

func (o *orderRepository) CreateLimit(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, price, qty string) (order *entities.Order, err error) {
	// 查询交易对配置
	tradeInfo, err := o.tradeVarietyRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	data := entities.Order{
		OrderId:   models_order.GenerateOrderId(side),
		UserId:    user_id,
		Symbol:    symbol,
		OrderSide: side,
		OrderType: matching_types.OrderTypeLimit,
		Price:     price,
		Quantity:  qty,
		NanoTime:  time.Now().UnixNano(),
		FeeRate:   tradeInfo.FeeRate,
		Status:    models_types.OrderStatusNew,
	}

	unfinished := entities.UnfinishedOrder{
		Order: data,
	}

	if err := o.validateOrder(ctx, tradeInfo, &data); err != nil {
		return nil, errors.Wrap(err, "validate order failed")
	}

	o.logger.Sugar().Infof("auto create tables: %s, %s", data.TableName(), unfinished.TableName())

	//auto create tables
	if err := o.db.Table(data.TableName()).AutoMigrate(&entities.Order{}); err != nil {
		return nil, errors.Wrap(err, "auto migrate order table failed")
	}
	if err := o.db.Table(unfinished.TableName()).AutoMigrate(&entities.UnfinishedOrder{}); err != nil {
		return nil, errors.Wrap(err, "auto migrate unfinished order table failed")
	}

	// 开启事务
	err = o.db.Transaction(func(tx *gorm.DB) (err error) {
		//冻结资产
		if data.OrderSide == matching_types.OrderSideSell {
			data.FreezeQty = data.Quantity
			err = o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.TargetVariety.Symbol, models_types.Amount(data.Quantity))
		} else {
			amount := models_types.Amount(data.Price).Mul(models_types.Amount(data.Quantity))
			fee := amount.Mul(models_types.Amount(data.FeeRate))
			data.FreezeAmount = amount.Add(fee).String()
			err = o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.BaseVariety.Symbol, models_types.Amount(data.FreezeAmount))
		}
		if err != nil {
			return err
		}

		if err := tx.Table(data.TableName()).Create(&data).Error; err != nil {
			return err
		}

		unfinished.Order = data
		if err := tx.Table(unfinished.TableName()).Create(&unfinished).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (o *orderRepository) CreateMarketByAmount(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, amount string) (order *entities.Order, err error) {
	//TODO implement me
	return nil, nil
}

func (o *orderRepository) CreateMarketByQty(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, qty string) (order *entities.Order, err error) {
	//TODO implement me
	return nil, nil
}

func (o *orderRepository) Cancel(ctx context.Context, order_id string, user_id *string) error {
	//TODO implement me
	return nil
}

func (o *orderRepository) validateOrder(ctx context.Context, tradeInfo *variety.TradeVariety, data *entities.Order) (err error) {

	//TODO 数量检查

	//TODO 价格检查

	//TODO 对向订单检查，防止自己的买单和卖单成交

	return nil
}
