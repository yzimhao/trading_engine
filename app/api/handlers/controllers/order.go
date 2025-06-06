package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/duolacloud/broker-core"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/app/common"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type OrderController struct {
	broker      broker.Broker
	logger      *zap.Logger
	repo        persistence.OrderRepository
	productRepo persistence.ProductRepository
}

type inOrderContext struct {
	fx.In
	Logger      *zap.Logger
	Broker      broker.Broker
	DB          *gorm.DB
	ProductRepo persistence.ProductRepository
	Repo        persistence.OrderRepository
}

func NewOrderController(in inOrderContext) *OrderController {
	return &OrderController{
		broker:      in.Broker,
		logger:      in.Logger,
		repo:        in.Repo,
		productRepo: in.ProductRepo,
	}
}

type OrderCreateRequest struct {
	Symbol    string                   `json:"symbol" binding:"required" example:"btcusdt"`
	Side      matching_types.OrderSide `json:"side" binding:"required" example:"buy"`
	OrderType matching_types.OrderType `json:"order_type" binding:"required" example:"limit"`
	Price     *decimal.Decimal         `json:"price,omitempty" example:"1.00"`
	Quantity  *decimal.Decimal         `json:"qty,omitempty" example:"12"`
	Amount    *decimal.Decimal         `json:"amount,omitempty"`
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

	userId := common.GetUserId(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseError(c, err)
		return
	}

	var (
		order *entities.Order
		err   error
		event models_types.EventOrderNew
	)

	if req.OrderType == matching_types.OrderTypeLimit {
		if req.Price == nil || req.Quantity == nil {
			common.ResponseError(c, errors.New("price and quantity are required"))
			return
		}
		order, err = ctrl.repo.CreateLimit(context.Background(), userId, req.Symbol, req.Side, *req.Price, *req.Quantity)
		if err != nil {
			ctrl.logger.Error("create limit order error", zap.Error(err), zap.Any("req", req))
			common.ResponseError(c, err)
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
			common.ResponseError(c, errors.New("amount or quantity is required"))
			return
		}

		if req.Amount != nil && req.Amount.Cmp(decimal.Zero) > 0 {
			order, err = ctrl.repo.CreateMarketByAmount(context.Background(), userId, req.Symbol, req.Side, *req.Amount)
			if err != nil {
				ctrl.logger.Error("create market by amount order error", zap.Error(err), zap.Any("req", req))
				common.ResponseError(c, err)
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
			order, err = ctrl.repo.CreateMarketByQty(context.Background(), userId, req.Symbol, req.Side, *req.Quantity)
			if err != nil {
				ctrl.logger.Error("create market by qty order error", zap.Error(err), zap.Any("req", req))
				common.ResponseError(c, err)
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
		ctrl.logger.Error("marshal order created event error", zap.Error(err), zap.Any("event", event))
		common.ResponseError(c, err)
		return
	}

	err = ctrl.broker.Publish(context.Background(), models_types.TOPIC_ORDER_NEW, &broker.Message{
		Body: body,
	}, broker.WithShardingKey(event.Symbol))

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
// @Param symbol query string true "symbol"
// @Param start query int64 true "start"
// @Param end query int64 true "end"
// @Param limit query int true "limit"
// @Success 200 {string} any
// @Router /api/v1/order/history [get]
func (ctrl *OrderController) HistoryList(c *gin.Context) {
	//TODO
	userId := common.GetUserId(c)
	symbol := c.DefaultQuery("symbol", "")
	start, _ := strconv.ParseInt(c.DefaultQuery("start", "0"), 10, 64)
	end, _ := strconv.ParseInt(c.DefaultQuery("end", "0"), 10, 64)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	orders, err := ctrl.repo.HistoryList(c, userId, symbol, start, end, limit)
	if err != nil {
		common.ResponseError(c, err)
		return
	}
	common.ResponseOK(c, orders)
}

// @Summary unfinished list
// @Description unfinished list
// @ID v1.order.unfinished
// @Tags order
// @Accept json
// @Produce json
// @Param symbol query string true "symbol"
// @Success 200 {string} any
// @Router /api/v1/order/unfinished [get]
func (ctrl *OrderController) UnfinishedList(c *gin.Context) {

	symbol := c.Query("symbol")
	// userId := common.GetUserId(c)

	// tradeVariety, err := ctrl.tradeVariety.FindBySymbol(c, symbol)
	// if err != nil {
	// 	common.ResponseError(c, err)
	// 	return
	// }

	//TODO 这个未完成订单列表 需要重新写一个方法
	orders, err := ctrl.repo.LoadUnfinishedOrders(c, symbol)
	if err != nil {
		common.ResponseError(c, err)
		return
	}
	common.ResponseOK(c, orders)
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
