package settlement

import (
	"fmt"

	"github.com/duolacloud/crud-core/cache"
	models_order "github.com/yzimhao/trading_engine/v2/internal/models/order"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SettleProcessor struct {
	db     *gorm.DB
	logger *zap.Logger
	cache  cache.Cache
}

type inContext struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Cache  cache.Cache
}

func NewSettleProcessor(in inContext) *SettleProcessor {
	return &SettleProcessor{
		db:     in.DB,
		logger: in.Logger,
		cache:  in.Cache,
	}
}

func (s *SettleProcessor) Run(tradeResult matching_types.TradeResult) error {
	//TODO 自动创建结算需要的表
	tradeLog := entities.TradeLog{Symbol: tradeResult.Symbol}
	if err := s.db.Table(tradeLog.TableName()).AutoMigrate(&tradeLog); err != nil {
		s.logger.Sugar().Errorf("auto migrate trade log table failed: %v", err)
		return err
	}

	//TODO 资产相关表

	return s.flow(tradeResult)
}

func (s *SettleProcessor) flow(tradeResult matching_types.TradeResult) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		//检查订单
		askOrder, bidOrder, err := s.checkOrder(tx, tradeResult)
		if err != nil {
			return err
		}

		//写交易记录
		tradeLog, err := s.writeTradeLog(tx, tradeResult, askOrder, bidOrder)
		if err != nil {
			return err
		}

		//更新ask订单信息
		if err := s.updateAskOrderInfo(tx, tradeLog, askOrder, tradeResult.RemainderMarketOrderId); err != nil {
			return err
		}

		//更新bid订单信息
		if err := s.updateBidOrderInfo(tx, tradeLog, bidOrder, tradeResult.RemainderMarketOrderId); err != nil {
			return err
		}

		//订单交割
		if err := s.orderDelivery(tx, tradeResult); err != nil {
			return err
		}

		return nil
	})
}

func (s *SettleProcessor) checkOrder(tx *gorm.DB, tradeResult matching_types.TradeResult) (askOrder, bidOrder *entities.Order, err error) {
	askOrder = &entities.Order{
		Symbol: tradeResult.Symbol,
	}
	bidOrder = &entities.Order{
		Symbol: tradeResult.Symbol,
	}

	if err := tx.Table(askOrder.TableName()).Where("order_id=?", tradeResult.AskOrderId).First(&askOrder).Error; err != nil {
		return nil, nil, err
	}
	if err := tx.Table(bidOrder.TableName()).Where("order_id=?", tradeResult.BidOrderId).First(&bidOrder).Error; err != nil {
		return nil, nil, err
	}

	if askOrder.Status != models_types.OrderStatusNew {
		return nil, nil, fmt.Errorf("invalid ask order status")
	}

	if bidOrder.Status != models_types.OrderStatusNew {
		return nil, nil, fmt.Errorf("invalid bid order status")
	}
	return askOrder, bidOrder, nil
}

func (s *SettleProcessor) writeTradeLog(tx *gorm.DB, tradeResult matching_types.TradeResult, askOrder, bidOrder *entities.Order) (*entities.TradeLog, error) {

	amount := models_types.Amount(tradeResult.TradeQuantity).Mul(models_types.Amount(tradeResult.TradePrice))

	tradeLog := entities.TradeLog{
		Symbol:     tradeResult.Symbol,
		TradeId:    models_order.GenerateTradeId(tradeResult.AskOrderId, tradeResult.BidOrderId),
		Ask:        tradeResult.AskOrderId,
		Bid:        tradeResult.BidOrderId,
		TradeBy:    tradeResult.TradeBy,
		AskUid:     askOrder.UserId,
		BidUid:     askOrder.UserId,
		Price:      tradeResult.TradePrice,
		Quantity:   tradeResult.TradeQuantity,
		Amount:     amount.String(),
		AskFeeRate: askOrder.FeeRate,
		AskFee:     amount.Mul(models_types.Amount(askOrder.FeeRate)).String(),
		BidFeeRate: bidOrder.FeeRate,
		BidFee:     amount.Mul(models_types.Amount(bidOrder.FeeRate)).String(),
	}

	if err := tx.Table(tradeLog.TableName()).Create(&tradeLog).Error; err != nil {
		return nil, err
	}
	return &tradeLog, nil
}

func (s *SettleProcessor) updateAskOrderInfo(tx *gorm.DB, tradeLog *entities.TradeLog, askOrder *entities.Order, remainderMarketOrderId string) error {

	askOrder.Fee = models_types.Amount(askOrder.Fee).Add(models_types.Amount(tradeLog.AskFee)).String()
	askOrder.FinishedQty = models_types.Amount(askOrder.FinishedQty).Add(models_types.Amount(tradeLog.Quantity)).String()
	askOrder.FinishedAmount = models_types.Amount(askOrder.FinishedAmount).Add(models_types.Amount(tradeLog.Amount)).String()
	askOrder.AvgPrice = models_types.Amount(askOrder.FinishedAmount).Div(models_types.Amount(askOrder.FinishedQty)).String()
	//初始状态为部分成交
	askOrder.Status = models_types.OrderStatusPartialFill

	if askOrder.OrderType == matching_types.OrderTypeLimit {
		be := models_types.Amount(askOrder.FinishedQty).Cmp(models_types.Amount(askOrder.Quantity))
		if be > 0 {
			return fmt.Errorf("invalid ask order finished qty")
		}
		if be == 0 {
			askOrder.Status = models_types.OrderStatusFilled
		}

		if err := tx.Table(askOrder.TableName()).Where("order_id=?", askOrder.OrderId).Updates(askOrder).Error; err != nil {
			return err
		}

		if askOrder.Status == models_types.OrderStatusNew {
			err := tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", askOrder.OrderId).Updates(askOrder).Error
			if err != nil {
				return err
			}
		} else {
			err := tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", askOrder.OrderId).Delete(askOrder).Error
			if err != nil {
				return err
			}
		}
	} else {
		//市价单结算
		if models_types.Amount(askOrder.Quantity).Equal(models_types.Amount(askOrder.FinishedQty)) || models_types.Amount(askOrder.Amount).Equal(models_types.Amount(askOrder.FinishedAmount)) {
			askOrder.Status = models_types.OrderStatusFilled
		}

		if remainderMarketOrderId == askOrder.OrderId {
			askOrder.Status = models_types.OrderStatusFilled
		}

		err := tx.Table(askOrder.TableName()).Where("order_id=?", askOrder.OrderId).Updates(askOrder).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SettleProcessor) updateBidOrderInfo(tx *gorm.DB, tradeLog *entities.TradeLog, bidOrder *entities.Order, remainderMarketOrderId string) error {

	bidOrder.Fee = models_types.Amount(bidOrder.Fee).Add(models_types.Amount(tradeLog.BidFee)).String()
	bidOrder.FinishedQty = models_types.Amount(bidOrder.FinishedQty).Add(models_types.Amount(tradeLog.Quantity)).String()
	bidOrder.FinishedAmount = models_types.Amount(bidOrder.FinishedAmount).Add(models_types.Amount(tradeLog.Amount)).String()
	bidOrder.AvgPrice = models_types.Amount(bidOrder.FinishedAmount).Div(models_types.Amount(bidOrder.FinishedQty)).String()
	//初始状态为部分成交
	bidOrder.Status = models_types.OrderStatusPartialFill

	if bidOrder.OrderType == matching_types.OrderTypeLimit {
		be := models_types.Amount(bidOrder.FinishedQty).Cmp(models_types.Amount(bidOrder.Quantity))
		if be > 0 {
			return fmt.Errorf("invalid bid order finished qty")
		}
		if be == 0 {
			bidOrder.Status = models_types.OrderStatusFilled
		}

		err := tx.Table(bidOrder.TableName()).Where("order_id=?", bidOrder.OrderId).Updates(bidOrder).Error
		if err != nil {
			return err
		}

		if bidOrder.Status == models_types.OrderStatusNew {
			err := tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", bidOrder.OrderId).Updates(bidOrder).Error
			if err != nil {
				return err
			}
		} else {
			err := tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", bidOrder.OrderId).Delete(bidOrder).Error
			if err != nil {
				return err
			}
		}
	} else {
		//市价单结算
		if models_types.Amount(bidOrder.Quantity).Equal(models_types.Amount(bidOrder.FinishedQty)) || models_types.Amount(bidOrder.Amount).Equal(models_types.Amount(bidOrder.FinishedAmount)) {
			bidOrder.Status = models_types.OrderStatusFilled
		}

		if remainderMarketOrderId == bidOrder.OrderId {
			bidOrder.Status = models_types.OrderStatusFilled
		}

		err := tx.Table(bidOrder.TableName()).Where("order_id=?", bidOrder.OrderId).Updates(bidOrder).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SettleProcessor) orderDelivery(tx *gorm.DB, tradeResult matching_types.TradeResult) error {
	//TODO
	return nil
}
