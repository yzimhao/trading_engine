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
	db     *gorm.DB
	logger *zap.Logger
}

var _ persistence.OrderRepository = (*orderRepository)(nil)

func NewOrderRepo(db *gorm.DB, logger *zap.Logger) persistence.OrderRepository {
	return &orderRepository{
		db:     db,
		logger: logger,
	}
}

func (o *orderRepository) CreateLimit(ctx context.Context, user_id, symbol string, side matching_types.OrderSide, price, qty string) (order *entities.Order, err error) {
	data := entities.Order{
		UserId:    user_id,
		Symbol:    symbol,
		OrderSide: side,
		OrderType: matching_types.OrderTypeLimit,
		Price:     price,
		Quantity:  qty,
		NanoTime:  time.Now().UnixNano(),
		FeeRate:   "0", //TODO
		Status:    models_types.OrderStatusNew,
	}
	unfinished := entities.UnfinishedOrder{
		Order: data,
	}

	//auto create tables
	o.db.AutoMigrate(&unfinished, &data)

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
