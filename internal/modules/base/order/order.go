package order

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"github.com/yzimhao/trading_engine/v2/internal/modules/middlewares"
	"github.com/yzimhao/trading_engine/v2/internal/modules/notification/ws"
	notification_ws "github.com/yzimhao/trading_engine/v2/internal/modules/notification/ws"
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
	produce   *provider.Produce
	auth      *middlewares.AuthMiddleware
	ws        *notification_ws.WsManager
}

func newOrderModule(
	router *provider.Router,
	logger *zap.Logger,
	produce *provider.Produce,
	auth *middlewares.AuthMiddleware,
	ws *notification_ws.WsManager,
	orderRepo persistence.OrderRepository) {
	o := orderModule{
		router:    router,
		logger:    logger,
		orderRepo: orderRepo,
		produce:   produce,
		auth:      auth,
		ws:        ws,
	}
	o.registerRouter()
}

func (o *orderModule) registerRouter() {
	orderGroup := o.router.APIv1.Group("/order")
	// 权限认证
	orderGroup.Use(o.auth.Auth())
	// 创建交易订单
	orderGroup.POST("/", o.create)
	orderGroup.GET("/cancel", o.cancel)
}

type CreateOrderRequest struct {
	Symbol    string                   `json:"symbol" binding:"required" example:"btcusdt"`
	Side      matching_types.OrderSide `json:"side" binding:"required" example:"SELL/BUY"`
	OrderType matching_types.OrderType `json:"order_type" binding:"required" example: "LIMIT/MARKET"`
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

		event.Price = func() decimal.Decimal {
			p := order.Price
			return p
		}()

		event.Quantity = func() decimal.Decimal {
			q := order.Quantity
			return q
		}()

		o.ws.SendTo(c, order.UserId, ws.MsgNewOrderTpl.Format(map[string]string{"symbol": order.Symbol}), map[string]any{
			"symbol":          order.Symbol,
			"order_id":        order.OrderId,
			"order_side":      order.OrderSide,
			"order_type":      order.OrderType,
			"price":           order.Price,
			"quantity":        order.Quantity,
			"avg_price":       order.AvgPrice,
			"finished_qty":    order.FinishedQty,
			"finished_amount": order.FinishedAmount,
			"at":              order.CreatedAt.Unix(),
		})

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

			event.MaxAmount = func() decimal.Decimal {
				a := order.FreezeAmount
				return a
			}()
		} else {
			order, err = o.orderRepo.CreateMarketByQty(context.Background(), userId, req.Symbol, req.Side, *req.Quantity)
			if err != nil {
				o.logger.Error("create market by qty order error", zap.Error(err), zap.Any("req", req))
				o.router.ResponseError(c, types.ErrInternalError)
				return
			}
		}

	}

	event.Symbol = order.Symbol
	event.OrderId = order.OrderId
	event.OrderSide = order.OrderSide
	event.OrderType = order.OrderType
	event.Quantity = order.Quantity
	event.Amount = order.Amount
	event.NanoTime = order.NanoTime
	event.MaxAmount = order.FreezeAmount
	event.MaxQty = order.FreezeQty

	body, err := json.Marshal(event)
	if err != nil {
		o.logger.Error("marshal order created event error", zap.Error(err), zap.Any("event", event))
		o.router.ResponseError(c, types.ErrInternalError)
		return
	}

	o.logger.Sugar().Debugf("create order event: %s", body)

	// err = o.broker.Publish(context.Background(), types.TOPIC_ORDER_NEW, &broker.Message{
	// 	Body: body,
	// }, broker.WithShardingKey(event.Symbol))

	err = o.produce.Publish(context.Background(), types.TOPIC_ORDER_NEW, body)
	if err != nil {
		o.logger.Error("publish order created event error", zap.Error(err))
		o.router.ResponseError(c, types.ErrInternalError)
		return
	}
	o.router.ResponseOk(c, gin.H{"order_id": order.OrderId})
}

// @Summary 创建订单
// @Description
// @ID v1.order
// @Tags order
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Param order_id query string true "order_id"
// @Success 200 {string} any
// @Router /api/v1/order/cancel [get]
func (o *orderModule) cancel(c *gin.Context) {
	symbol := c.Query("symbol")
	orderId := c.Query("order_id")
	o.orderRepo.Cancel(c, symbol, orderId, matching_types.RemoveItemTypeByUser)
	o.router.ResponseOk(c, nil)
}
