package settlement

import (
	"context"
	"fmt"

	"github.com/duolacloud/crud-core/cache"
	models_order "github.com/yzimhao/trading_engine/v2/internal/models/order"
	models_types "github.com/yzimhao/trading_engine/v2/internal/models/types"
	"github.com/yzimhao/trading_engine/v2/internal/models/variety"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/gorm/entities"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SettleProcessor struct {
	db               *gorm.DB
	logger           *zap.Logger
	cache            cache.Cache
	tradeVarietyRepo persistence.TradeVarietyRepository
	assetRepo        persistence.AssetRepository
}

type inContext struct {
	DB               *gorm.DB
	Logger           *zap.Logger
	Cache            cache.Cache
	TradeVarietyRepo persistence.TradeVarietyRepository
	AssetRepo        persistence.AssetRepository
}

func NewSettleProcessor(in inContext) *SettleProcessor {
	return &SettleProcessor{
		db:               in.DB,
		logger:           in.Logger,
		cache:            in.Cache,
		tradeVarietyRepo: in.TradeVarietyRepo,
		assetRepo:        in.AssetRepo,
	}
}

func (s *SettleProcessor) Run(ctx context.Context, tradeResult matching_types.TradeResult) error {
	//TODO 自动创建结算需要的表
	tradeLog := entities.TradeLog{Symbol: tradeResult.Symbol}
	if err := s.db.Table(tradeLog.TableName()).AutoMigrate(&tradeLog); err != nil {
		s.logger.Sugar().Errorf("auto migrate trade log table failed: %v", err)
		return err
	}

	//TODO 资产相关表

	return s.flow(ctx, tradeResult)
}

func (s *SettleProcessor) flow(ctx context.Context, tradeResult matching_types.TradeResult) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		//交易对配置
		tradePairInfo, err := s.tradeVarietyRepo.FindBySymbol(ctx, tradeResult.Symbol)
		if err != nil {
			return err
		}

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
		if err := s.orderDelivery(tx, tradeLog, askOrder, bidOrder, tradePairInfo); err != nil {
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

func (s *SettleProcessor) orderDelivery(
	tx *gorm.DB,
	tradeLog *entities.TradeLog,
	ask, bid *entities.Order,
	tradePairInfo *variety.TradeVariety,
) error {
	ctx := context.Background()
	//买家结算被交易物品
	// 1.解冻卖家的冻结数量
	// 2.将解冻的数量转移给买家
	err := s.assetRepo.UnFreeze(ctx, tx, tradeLog.TradeId, ask.UserId,
		tradePairInfo.TargetVariety.Symbol, models_types.Amount(tradeLog.Quantity))
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery target variety unFreeze: %v %s", tradeLog, err.Error())
		return err
	}
	err = s.assetRepo.TransferWithTx(ctx, tx, tradeLog.TradeId, ask.UserId, bid.UserId,
		tradePairInfo.TargetVariety.Symbol, models_types.Amount(tradeLog.Quantity))
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery target variety transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//卖家结算本位币
	// 1.解冻买家冻结的金额
	// 2.将解冻的金额扣除双方手续费后，转入卖家账户
	amount := models_types.Amount(tradeLog.Amount)
	err = s.assetRepo.UnFreeze(ctx, tx,
		tradeLog.TradeId, bid.UserId, tradePairInfo.BaseVariety.Symbol,
		amount.Add(models_types.Amount(tradeLog.AskFee).Add(models_types.Amount(tradeLog.BidFee))))
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery base variety unFreeze: %v %s", tradeLog, err.Error())
		return err
	}

	err = s.assetRepo.TransferWithTx(ctx, tx, tradeLog.TradeId, bid.UserId, ask.UserId,
		tradePairInfo.BaseVariety.Symbol, amount.Sub(models_types.Amount(tradeLog.AskFee)))
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery base variety transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//ask手续费收入到系统账号里
	err = s.assetRepo.TransferWithTx(ctx, tx, tradeLog.TradeId, ask.UserId, entities.SYSTEM_USER_FEE,
		tradePairInfo.BaseVariety.Symbol, models_types.Amount(tradeLog.AskFee))
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery ask fee transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//bid的手续费收入到系统账号里
	err = s.assetRepo.TransferWithTx(ctx, tx, tradeLog.TradeId, bid.UserId, entities.SYSTEM_USER_FEE,
		tradePairInfo.BaseVariety.Symbol, models_types.Amount(tradeLog.BidFee))
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery bid fee transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//订单状态为已成交，则解冻该订单冻结的全部数量
	if ask.Status == models_types.OrderStatusFilled {
		err = s.assetRepo.UnFreeze(ctx, tx, tradeLog.TradeId, ask.UserId, tradePairInfo.TargetVariety.Symbol, "0")
		if err != nil {
			s.logger.Sugar().Errorf("orderDelivery ask unFreeze: %v %s", tradeLog, err.Error())
			return err
		}
	}
	if bid.Status == models_types.OrderStatusFilled {
		err = s.assetRepo.UnFreeze(ctx, tx, tradeLog.TradeId, bid.UserId, tradePairInfo.BaseVariety.Symbol, "0")
		if err != nil {
			s.logger.Sugar().Errorf("orderDelivery bid unFreeze: %v %s", tradeLog, err.Error())
			return err
		}
	}

	return nil
}