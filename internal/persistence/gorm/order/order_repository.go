package order

import (
	"context"
	"time"

	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
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
}

var _ persistence.OrderRepository = (*orderRepository)(nil)

func NewOrderRepo(db *gorm.DB, logger *zap.Logger, tradeVarietyRepo persistence.TradeVarietyRepository) persistence.OrderRepository {
	return &orderRepository{
		db:               db,
		logger:           logger,
		tradeVarietyRepo: tradeVarietyRepo,
	}
}

func (o *orderRepository) CreateLimit(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, price, qty string) (order *entities.Order, err error) {
	// 查询交易对配置
	tradeInfo, err := o.tradeVarietyRepo.FindBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	data := entities.Order{
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

	//auto create tables
	if err := o.db.AutoMigrate(&unfinished, &data); err != nil {
		return nil, err
	}

	// 开启事务
	err = o.db.Transaction(func(tx *gorm.DB) error {
		//冻结资产

		return nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
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
