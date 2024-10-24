package controllers

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/duolacloud/broker-core"
	rocketmq "github.com/duolacloud/broker-rocketmq"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/common"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	gorm_order "github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/order"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderController struct {
	broker broker.Broker
	logger *zap.Logger
	repo   persistence.OrderRepository
}

type inOrderContext struct {
	fx.In
	Logger *zap.Logger
	Broker broker.Broker
	DB     *gorm.DB
}

func NewOrderController(in inOrderContext) *OrderController {
	repo := gorm_order.NewOrderRepo(in.DB, in.Logger)
	return &OrderController{
		broker: in.Broker,
		logger: in.Logger,
		repo:   repo,
	}
}

type OrderCreateRequest struct {
	Symbol    string                   `json:"symbol" binding:"required"`
	Side      matching_types.OrderSide `json:"side" binding:"required"`
	OrderType matching_types.OrderType `json:"order_type" binding:"required"`
	Price     *string                  `json:"price,omitempty" example:"1.00"`
	Quantity  *string                  `json:"qty,omitempty" example:"12"`
	Amount    *string                  `json:"amount,omitempty" example:"100.00"`
}

// @Summary create order
// @Description create order
// @ID v1.order.create
// @Tags order
// @Accept json
// @Produce json
// @Param args body OrderCreateRequest true "args"
// @Success 200 {string} any
// @Router /api/v1/order/create [post]
func (ctrl *OrderController) Create(c *gin.Context) {
	var req OrderCreateRequest

	userId := c.MustGet("userId").(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, err)
		return
	}

	var order *entities.Order
	var err error
	if req.OrderType == matching_types.OrderTypeLimit {
		if req.Price == nil || req.Quantity == nil {
			common.ResponseError(c, errors.New("price and quantity are required"))
			return
		}
		order, err = ctrl.repo.CreateLimit(context.Background(), userId, req.Symbol, req.Side, *req.Price, *req.Quantity)
	} else {
		if req.Amount == nil && req.Quantity == nil {
			common.ResponseError(c, errors.New("amount or quantity is required"))
			return
		}

		if req.Amount != nil && models_types.Amount(*req.Amount).Cmp(models_types.Amount("0")) > 0 {
			order, err = ctrl.repo.CreateMarketByAmount(context.Background(), userId, req.Symbol, req.Side, *req.Amount)
		} else {
			order, err = ctrl.repo.CreateMarketByQty(context.Background(), userId, req.Symbol, req.Side, *req.Quantity)
		}
	}
	if err != nil {
		ctrl.logger.Error("create order error", zap.Error(err))
		common.ResponseError(c, err)
		return
	}

	event := models_types.EventOrderNew{
		Symbol:    order.Symbol,
		OrderId:   order.OrderId,
		OrderSide: order.OrderSide,
		OrderType: order.OrderType,
		Price:     &order.Price,
		Quantity:  &order.Quantity,
		Amount:    &order.Amount,
		NanoTime:  order.NanoTime,
	}

	body, err := json.Marshal(event)
	if err != nil {
		ctrl.logger.Error("marshal order created event error", zap.Error(err))
		common.ResponseError(c, err)
		return
	}

	err = ctrl.broker.Publish(context.Background(), models_types.TOPIC_ORDER_NEW, &broker.Message{
		Body: body,
	}, rocketmq.WithShardingKey(event.Symbol))

	if err != nil {
		ctrl.logger.Error("publish order created event error", zap.Error(err))
		common.ResponseError(c, err)
		return
	}

	common.ResponseOK(c, gin.H{"order_id": order.OrderId})
}

// @Summary history list
// @Description history list
// @ID v1.order.history
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/history [get]
func (ctrl *OrderController) HistoryList(c *gin.Context) {
	common.ResponseOK(c, "test")
}

// @Summary unfinished list
// @Description unfinished list
// @ID v1.order.unfinished
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/unfinished [get]
func (ctrl *OrderController) UnfinishedList(c *gin.Context) {
	common.ResponseOK(c, "test")
}

// @Summary trade history list
// @Description trade history list
// @ID v1.order.trade_history
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/trade/history [get]
func (ctrl *OrderController) TradeHistoryList(c *gin.Context) {
	common.ResponseOK(c, "test")
}
