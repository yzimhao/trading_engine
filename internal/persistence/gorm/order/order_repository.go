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

	if err := o.validateOrderLimit(ctx, tradeInfo, &data); err != nil {
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
			_, err := o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.TargetVariety.Symbol, models_types.Amount(data.Quantity))
			if err != nil {
				return err
			}
		} else {
			amount := models_types.Amount(data.Price).Mul(models_types.Amount(data.Quantity))
			fee := amount.Mul(models_types.Amount(data.FeeRate))
			data.FreezeAmount = amount.Add(fee).String()
			_, err := o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.BaseVariety.Symbol, models_types.Amount(data.FreezeAmount))
			if err != nil {
				return err
			}
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
	tradeInfo, err := o.tradeVarietyRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	data := entities.Order{
		OrderId:      models_order.GenerateOrderId(side),
		UserId:       user_id,
		Symbol:       symbol,
		OrderSide:    side,
		OrderType:    matching_types.OrderTypeMarket,
		FeeRate:      tradeInfo.FeeRate,
		FreezeAmount: amount,
		Status:       models_types.OrderStatusNew,
		NanoTime:     time.Now().UnixNano(),
	}

	if err := o.db.Table(data.TableName()).AutoMigrate(&entities.Order{}); err != nil {
		return nil, errors.Wrap(err, "auto migrate order table failed")
	}

	if err := o.validateOrderMarketAmount(ctx, tradeInfo, &data); err != nil {
		return nil, errors.Wrap(err, "validate order failed")
	}

	err = o.db.Transaction(func(tx *gorm.DB) (err error) {
		if data.OrderSide == matching_types.OrderSideSell {
			f, err := o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.TargetVariety.Symbol, models_types.Amount("0"))
			if err != nil {
				return err
			}
			data.FreezeQty = f.FreezeAmount.String()
		} else {
			f, err := o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.BaseVariety.Symbol, models_types.Amount(data.Amount))
			if err != nil {
				return err
			}
			data.FreezeAmount = f.FreezeAmount.String()
			fee := models_types.Amount(data.FreezeAmount).Mul(models_types.Amount(data.FeeRate))
			data.Amount = models_types.Amount(data.FreezeAmount).Sub(fee).String()
		}

		if err := tx.Table(data.TableName()).Create(&data).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (o *orderRepository) CreateMarketByQty(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, qty string) (order *entities.Order, err error) {
	tradeInfo, err := o.tradeVarietyRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	data := entities.Order{
		OrderId:   models_order.GenerateOrderId(side),
		UserId:    user_id,
		Symbol:    symbol,
		OrderSide: side,
		OrderType: matching_types.OrderTypeMarket,
		FeeRate:   tradeInfo.FeeRate,
		Quantity:  qty,
		Status:    models_types.OrderStatusNew,
		NanoTime:  time.Now().UnixNano(),
	}

	if err := o.validateOrderMarketQty(ctx, tradeInfo, &data); err != nil {
		return nil, errors.Wrap(err, "validate order failed")
	}

	if err := o.db.Table(data.TableName()).AutoMigrate(&entities.Order{}); err != nil {
		return nil, errors.Wrap(err, "auto migrate order table failed")
	}

	err = o.db.Transaction(func(tx *gorm.DB) (err error) {
		if data.OrderSide == matching_types.OrderSideSell {
			f, err := o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.TargetVariety.Symbol, models_types.Amount(data.Quantity))
			if err != nil {
				return err
			}
			data.FreezeQty = f.FreezeAmount.String()
		} else {
			f, err := o.assetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, tradeInfo.BaseVariety.Symbol, models_types.Amount("0"))
			if err != nil {
				return err
			}
			data.FreezeAmount = f.FreezeAmount.String()
		}

		if err := tx.Table(data.TableName()).Create(&data).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (o *orderRepository) Cancel(ctx context.Context, symbol, order_id string, cancelType models_types.CancelType) error {
	o.logger.Sugar().Infof("cancel order: %s, %s, %s, %d", symbol, order_id, cancelType)

	tradeInfo, err := o.tradeVarietyRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return err
	}

	order := &entities.Order{
		Symbol:  symbol,
		OrderId: order_id,
	}
	unfinished := &entities.UnfinishedOrder{
		Order: *order,
	}
	//TODO检查是否有结算锁

	if err := o.db.Table(order.TableName()).Where("order_id=?", order_id).First(order).Error; err != nil {
		o.logger.Sugar().Errorf("order query error: %v, symbol: %s, order_id: %s", err, symbol, order_id)
		return err
	}

	err = o.db.Transaction(func(tx *gorm.DB) (err error) {
		//解冻资产
		if order.OrderSide == matching_types.OrderSideSell {
			err := o.assetRepo.UnFreeze(ctx, tx, order.OrderId, order.UserId, tradeInfo.TargetVariety.Symbol, models_types.Amount("0"))
			if err != nil {
				return err
			}
		} else {
			err := o.assetRepo.UnFreeze(ctx, tx, order.OrderId, order.UserId, tradeInfo.BaseVariety.Symbol, models_types.Amount("0"))
			if err != nil {
				return err
			}
		}

		//更新订单状态
		if err := tx.Table(order.TableName()).Where("order_id=?", order_id).Update("status", models_types.OrderStatusCanceled).Error; err != nil {
			return err
		}

		//删除未完成订单
		if err := tx.Table(unfinished.TableName()).Where("order_id=?", order_id).Delete(unfinished).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (o *orderRepository) validateOrderLimit(ctx context.Context, tradeInfo *variety.TradeVariety, data *entities.Order) (err error) {

	//TODO 数量检查

	//TODO 价格检查

	//TODO 对向订单检查，防止自己的买单和卖单成交

	return nil
}

func (o *orderRepository) validateOrderMarketAmount(ctx context.Context, tradeInfo *variety.TradeVariety, data *entities.Order) (err error) {
	//TODO implement me
	return nil
}

func (o *orderRepository) validateOrderMarketQty(ctx context.Context, tradeInfo *variety.TradeVariety, data *entities.Order) (err error) {
	//TODO implement me
	return nil
}
