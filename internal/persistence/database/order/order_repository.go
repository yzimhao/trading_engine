package order

import (
	"context"
	"time"

	"github.com/pkg/errors"
	models_order "github.com/yzimhao/trading_engine/v2/internal/models/order"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type orderRepository struct {
	db            *gorm.DB
	logger        *zap.Logger
	productRepo   persistence.ProductRepository
	userAssetRepo persistence.UserAssetRepository
}

var _ persistence.OrderRepository = (*orderRepository)(nil)

func NewOrderRepo(
	db *gorm.DB,
	logger *zap.Logger,
	productRepo persistence.ProductRepository,
	userAssetRepo persistence.UserAssetRepository,
) persistence.OrderRepository {
	return &orderRepository{
		db:            db,
		logger:        logger,
		productRepo:   productRepo,
		userAssetRepo: userAssetRepo,
	}
}

func (o *orderRepository) LoadUnfinishedOrders(ctx context.Context, symbol string) (orders []*entities.Order, err error) {
	//TODO 分批读取
	unfinished := entities.UnfinishedOrder{}
	o.db.Table(unfinished.TableName()).Where("symbol=?", symbol).Order("nano_time asc").Find(&orders)
	return orders, nil
}

func (o *orderRepository) HistoryList(ctx context.Context, user_id, symbol string, start, end int64, limit int) (orders []*entities.Order, err error) {
	//TODO 分批读取
	entity := entities.Order{Symbol: symbol}
	o.db.Table(entity.TableName()).Where("user_id=? and symbol=? and nano_time>=? and nano_time<=?", user_id, symbol, start, end).Order("nano_time asc").Limit(limit).Find(&orders)
	return orders, nil
}

func (o *orderRepository) CreateLimit(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, price, qty string) (order *entities.Order, err error) {
	// 查询交易对配置
	product, err := o.productRepo.Get(symbol)
	if err != nil {
		return nil, errors.Wrap(err, "find trade variety failed")
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
		FeeRate:   product.FeeRate,
		Status:    models_types.OrderStatusNew,
	}

	unfinished := entities.UnfinishedOrder{
		Order: data,
	}

	if err := o.validateOrderLimit(ctx, product, &data); err != nil {
		return nil, errors.Wrap(err, "validate order failed")
	}

	o.logger.Sugar().Infof("auto create tables: %s, %s", data.TableName(), unfinished.TableName())

	//auto create tables
	if !o.db.Migrator().HasTable(data.TableName()) {
		if err := o.db.Table(data.TableName()).AutoMigrate(&entities.Order{}); err != nil {
			return nil, errors.Wrap(err, "auto migrate order table failed")
		}
	}

	if !o.db.Migrator().HasTable(unfinished.TableName()) {
		if err := o.db.Table(unfinished.TableName()).AutoMigrate(&entities.UnfinishedOrder{}); err != nil {
			return nil, errors.Wrap(err, "auto migrate unfinished order table failed")
		}
	}

	// 开启事务
	err = o.db.Transaction(func(tx *gorm.DB) (err error) {
		//冻结资产
		if data.OrderSide == matching_types.OrderSideSell {
			data.FreezeQty = data.Quantity
			_, err := o.userAssetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, product.Target.Symbol, models_types.Numeric(data.Quantity))
			if err != nil {
				return errors.Wrap(err, "freeze asset failed")
			}
		} else {
			amount := models_types.Numeric(data.Price).Mul(models_types.Numeric(data.Quantity))
			fee := amount.Mul(models_types.Numeric(data.FeeRate))
			data.FreezeAmount = amount.Add(fee).String()
			_, err := o.userAssetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, product.Base.Symbol, models_types.Numeric(data.FreezeAmount))
			if err != nil {
				return errors.Wrap(err, "freeze asset failed")
			}
		}

		if err := tx.Table(data.TableName()).Create(&data).Error; err != nil {
			return errors.Wrap(err, "create order failed")
		}

		unfinished.Order = data
		if err := tx.Table(unfinished.TableName()).Create(&unfinished).Error; err != nil {
			return errors.Wrap(err, "create unfinished order failed")
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "create order failed")
	}

	return &data, nil
}

func (o *orderRepository) CreateMarketByAmount(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, amount string) (order *entities.Order, err error) {
	product, err := o.productRepo.Get(symbol)
	if err != nil {
		return nil, err
	}

	data := entities.Order{
		OrderId:      models_order.GenerateOrderId(side),
		UserId:       user_id,
		Symbol:       symbol,
		OrderSide:    side,
		OrderType:    matching_types.OrderTypeMarket,
		FeeRate:      product.FeeRate,
		FreezeAmount: amount,
		Status:       models_types.OrderStatusNew,
		NanoTime:     time.Now().UnixNano(),
	}

	if err := o.db.Table(data.TableName()).AutoMigrate(&entities.Order{}); err != nil {
		return nil, errors.Wrap(err, "auto migrate order table failed")
	}

	if err := o.validateOrderMarketAmount(ctx, product, &data); err != nil {
		return nil, errors.Wrap(err, "validate order failed")
	}

	err = o.db.Transaction(func(tx *gorm.DB) (err error) {
		if data.OrderSide == matching_types.OrderSideSell {
			f, err := o.userAssetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, product.Target.Symbol, models_types.NumericZero)
			if err != nil {
				return err
			}
			data.FreezeQty = f.FreezeAmount.String()
		} else {
			f, err := o.userAssetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, product.Base.Symbol, models_types.Numeric(data.FreezeAmount))
			if err != nil {
				return err
			}
			data.FreezeAmount = f.FreezeAmount.String()
			fee := models_types.Numeric(data.FreezeAmount).Mul(models_types.Numeric(data.FeeRate))
			data.Amount = models_types.Numeric(data.FreezeAmount).Sub(fee).String()
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
	product, err := o.productRepo.Get(symbol)
	if err != nil {
		return nil, err
	}

	data := entities.Order{
		OrderId:   models_order.GenerateOrderId(side),
		UserId:    user_id,
		Symbol:    symbol,
		OrderSide: side,
		OrderType: matching_types.OrderTypeMarket,
		FeeRate:   product.FeeRate,
		Quantity:  qty,
		Status:    models_types.OrderStatusNew,
		NanoTime:  time.Now().UnixNano(),
	}

	if err := o.validateOrderMarketQty(ctx, product, &data); err != nil {
		return nil, errors.Wrap(err, "validate order failed")
	}

	if err := o.db.Table(data.TableName()).AutoMigrate(&entities.Order{}); err != nil {
		return nil, errors.Wrap(err, "auto migrate order table failed")
	}

	err = o.db.Transaction(func(tx *gorm.DB) (err error) {
		if data.OrderSide == matching_types.OrderSideSell {
			f, err := o.userAssetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, product.Target.Symbol, models_types.Numeric(data.Quantity))
			if err != nil {
				return err
			}
			data.FreezeQty = f.FreezeAmount.String()
		} else {
			f, err := o.userAssetRepo.Freeze(ctx, tx, data.OrderId, data.UserId, product.Base.Symbol, models_types.NumericZero)
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

	product, err := o.productRepo.Get(symbol)
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
			err := o.userAssetRepo.UnFreeze(ctx, tx, order.OrderId, order.UserId, product.Target.Symbol, models_types.NumericZero)
			if err != nil {
				return err
			}
		} else {
			err := o.userAssetRepo.UnFreeze(ctx, tx, order.OrderId, order.UserId, product.Base.Symbol, models_types.NumericZero)
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

func (o *orderRepository) validateOrderLimit(ctx context.Context, product *entities.Product, data *entities.Order) (err error) {

	//TODO 数量检查

	//TODO 价格检查

	//TODO 对向订单检查，防止自己的买单和卖单成交

	return nil
}

func (o *orderRepository) validateOrderMarketAmount(ctx context.Context, product *entities.Product, data *entities.Order) (err error) {
	//TODO implement me
	return nil
}

func (o *orderRepository) validateOrderMarketQty(ctx context.Context, product *entities.Product, data *entities.Order) (err error) {
	//TODO implement me
	return nil
}
