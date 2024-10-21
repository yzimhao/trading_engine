package controllers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/duolacloud/broker-core"
	rocketmq "github.com/duolacloud/broker-rocketmq"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/v2/app/api/handlers/common"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type OrderController struct {
	broker broker.Broker
	logger *zap.Logger
}

type inOrderContext struct {
	fx.In
	Logger *zap.Logger
	Broker broker.Broker
}

func NewOrderController(in inOrderContext) *OrderController {
	return &OrderController{
		broker: in.Broker,
		logger: in.Logger,
	}
}

// @Summary create order
// @Description create order
// @ID v1.order.create
// @Tags order
// @Accept json
// @Produce json
// @Success 200 {string} any
// @Router /api/v1/order/create [post]
func (ctrl *OrderController) Create(c *gin.Context) {
	//TODO

	event := models_types.EventOrderNew{
		Symbol: "BTCUSDT",
		At:     time.Now().Unix(),
		// ...
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

	common.ResponseOK(c, "ok")
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
