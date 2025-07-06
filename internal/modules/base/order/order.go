package order

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/duolacloud/broker-core"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	"github.com/yzimhao/trading_engine/v2/internal/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"base.order",
	fx.Invoke(newOrderModule),
)

type orderModule struct {
	router    *provider.Router
	logger    *zap.Logger
	orderRepo persistence.OrderRepository
	broker    broker.Broker
	auth      *middlewares.AuthMiddleware
}

func newOrderModule(
	router *provider.Router,
	logger *zap.Logger,
	broker broker.Broker,
	auth *middlewares.AuthMiddleware,
	orderRepo persistence.OrderRepository) {
	o := orderModule{
		router:    router,
		logger:    logger,
		orderRepo: orderRepo,
		broker:    broker,
		auth:      auth,
	}
	o.registerRouter()
}

func (o *orderModule) registerRouter() {
	orderGroup := o.router.APIv1.Group("/order")
	// 权限认证
	orderGroup.Use(o.auth.Auth())
	// 创建交易订单
	orderGroup.POST("/", o.create)
}

type CreateOrderRequest struct {
	Symbol    string                   `json:"symbol" binding:"required" example:"btcusdt"`
	Side      matching_types.OrderSide `json:"side" binding:"required" example:"buy"`
	OrderType matching_types.OrderType `json:"order_type" binding:"required" example:"limit"`
	Price     *decimal.Decimal         `json:"price,omitempty" example:"1.00"`
	Quantity  *decimal.Decimal         `json:"qty,omitempty" example:"12"`
	Amount    *decimal.Decimal         `json:"amount,omitempty"`
}

// @Summary 创建订单
// @Description
// @ID v1.order
// @Tags order
// @Accept json
// @Produce json
// @Param args body CreateOrderRequest true "args"
// @Success 200 {string} any
// @Router /api/v1/order [post]
func (o *orderModule) create(c *gin.Context) {
	var req CreateOrderRequest

	userId := o.router.ParseUserID(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		o.router.ResponseError(c, types.ErrInvalidParam)
		return
	}

	var (
		order *entities.Order
		err   error
		event types.EventOrderNew
	)

	req.Symbol = strings.ToLower(req.Symbol)
	if req.OrderType == matching_types.OrderTypeLimit {
		if req.Price == nil || req.Quantity == nil {
			o.logger.Sugar().Warnf("price or quantity is required")
			o.router.ResponseError(c, types.ErrInvalidParam)
			return
		}
		order, err = o.orderRepo.CreateLimit(context.Background(), userId, req.Symbol, req.Side, *req.Price, *req.Quantity)
		if err != nil {
			o.logger.Error("create limit order error", zap.Error(err), zap.Any("req", req))
			o.router.ResponseError(c, types.ErrInternalError)
			return
		}

		event.Price = func() *decimal.Decimal {
			p := order.Price
			return &p
		}()

		event.Quantity = func() *decimal.Decimal {
			q := order.Quantity
			return &q
		}()

	} else {
		if req.Amount == nil && req.Quantity == nil {
			o.logger.Sugar().Warnf("amount or quantity is required")
			o.router.ResponseError(c, types.ErrInvalidParam)
			return
		}

		if req.Amount != nil && req.Amount.Cmp(decimal.Zero) > 0 {
			order, err = o.orderRepo.CreateMarketByAmount(context.Background(), userId, req.Symbol, req.Side, *req.Amount)
			if err != nil {
				o.logger.Error("create market by amount order error", zap.Error(err), zap.Any("req", req))
				o.router.ResponseError(c, types.ErrInternalError)
				return
			}

			event.Amount = func() *decimal.Decimal {
				a := order.Amount
				return &a
			}()
			event.MaxAmount = func() *decimal.Decimal {
				a := order.FreezeAmount
				return &a
			}()
		} else {
			order, err = o.orderRepo.CreateMarketByQty(context.Background(), userId, req.Symbol, req.Side, *req.Quantity)
			if err != nil {
				o.logger.Error("create market by qty order error", zap.Error(err), zap.Any("req", req))
				o.router.ResponseError(c, types.ErrInternalError)
				return
			}

			event.Quantity = func() *decimal.Decimal {
				q := order.Quantity
				return &q
			}()
			event.MaxQty = func() *decimal.Decimal {
				q := order.FreezeQty
				return &q
			}()
		}
	}

	event.Symbol = order.Symbol
	event.OrderId = order.OrderId
	event.OrderSide = order.OrderSide
	event.OrderType = order.OrderType
	event.NanoTime = order.NanoTime

	body, err := json.Marshal(event)
	if err != nil {
		o.logger.Error("marshal order created event error", zap.Error(err), zap.Any("event", event))
		o.router.ResponseError(c, types.ErrInternalError)
		return
	}

	err = o.broker.Publish(context.Background(), types.TOPIC_ORDER_NEW, &broker.Message{
		Body: body,
	}, broker.WithShardingKey(event.Symbol))

	if err != nil {
		o.logger.Error("publish order created event error", zap.Error(err))
		o.router.ResponseError(c, types.ErrInternalError)
		return
	}
	o.router.ResponseOk(c, gin.H{"order_id": order.OrderId})
}
