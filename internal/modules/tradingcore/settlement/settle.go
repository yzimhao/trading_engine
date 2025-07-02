package settlement

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/duolacloud/broker-core"
	"github.com/duolacloud/crud-core/cache"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/app/webws"
	"github.com/yzimhao/trading_engine/v2/internal/persistence"
	"github.com/yzimhao/trading_engine/v2/internal/persistence/database/entities"
	models_types "github.com/yzimhao/trading_engine/v2/internal/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SettleProcessor struct {
	db     *gorm.DB
	logger *zap.Logger
	// cache            cache.Cache
	productRepo   persistence.ProductRepository
	userAssetRepo persistence.UserAssetRepository
	broker        broker.Broker
	redis         *redis.Client
	locker        *SettleLocker
	ws            *webws.WsManager
}

type inSettleContext struct {
	fx.In
	DB            *gorm.DB
	Logger        *zap.Logger
	Cache         cache.Cache
	ProductRepo   persistence.ProductRepository
	UserAssetRepo persistence.UserAssetRepository
	Broker        broker.Broker
	Redis         *redis.Client
	Locker        *SettleLocker
	Ws            *webws.WsManager
}

func NewSettleProcessor(in inSettleContext) *SettleProcessor {
	return &SettleProcessor{
		db:     in.DB,
		logger: in.Logger,
		// cache:            in.Cache,
		productRepo:   in.ProductRepo,
		userAssetRepo: in.UserAssetRepo,
		broker:        in.Broker,
		redis:         in.Redis,
		locker:        in.Locker,
		ws:            in.Ws,
	}
}

func (s *SettleProcessor) Run(ctx context.Context, tradeResult matching_types.TradeResult) error {
	//自动创建结算需要的表
	tradeLog := entities.TradeRecord{Symbol: tradeResult.Symbol}
	if err := s.db.Table(tradeLog.TableName()).AutoMigrate(&tradeLog); err != nil {
		s.logger.Sugar().Errorf("auto migrate trade log table failed: %v", err)
		return err
	}

	//创建买卖双方订单标志锁
	if err := s.locker.Lock(ctx, tradeResult.AskOrderId, tradeResult.BidOrderId); err != nil {
		return err
	}
	defer s.locker.Unlock(ctx, tradeResult.AskOrderId, tradeResult.BidOrderId)

	return s.flow(ctx, tradeResult)
}

func (s *SettleProcessor) flow(ctx context.Context, tradeResult matching_types.TradeResult) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		//交易对配置
		product, err := s.productRepo.Get(tradeResult.Symbol)
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
		if err := s.orderDelivery(tx, tradeLog, askOrder, bidOrder, product); err != nil {
			return err
		}

		// notify quote cacluate kline
		notifyQuote := models_types.EventNotifyQuote{
			TradeResult: tradeResult,
		}
		body, err := json.Marshal(notifyQuote)
		if err != nil {
			return err
		}
		if err := s.broker.Publish(ctx, models_types.TOPIC_NOTIFY_QUOTE, &broker.Message{
			Body: body,
		}, broker.WithShardingKey(tradeResult.Symbol)); err != nil {
			return err
		}

		//推送交易页面上的最新成交记录
		s.ws.Broadcast(ctx, webws.MsgTradeTpl.Format(map[string]string{"symbol": tradeResult.Symbol}),
			map[string]any{
				"price":    tradeLog.Price,
				"qty":      tradeLog.Quantity,
				"amount":   tradeLog.Amount,
				"trade_at": tradeResult.TradeTime,
			},
		)
		//TODO 推送买卖双方个人结算的成交记录

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

func (s *SettleProcessor) writeTradeLog(tx *gorm.DB, tradeResult matching_types.TradeResult, askOrder, bidOrder *entities.Order) (*entities.TradeRecord, error) {

	amount := tradeResult.TradeQuantity.Mul(tradeResult.TradePrice)

	tradeLog := entities.TradeRecord{
		Symbol:     tradeResult.Symbol,
		TradeId:    generateTradeId(tradeResult.AskOrderId, tradeResult.BidOrderId),
		Ask:        tradeResult.AskOrderId,
		Bid:        tradeResult.BidOrderId,
		TradeBy:    tradeResult.TradeBy,
		AskUid:     askOrder.UserId,
		BidUid:     bidOrder.UserId,
		Price:      tradeResult.TradePrice,
		Quantity:   tradeResult.TradeQuantity,
		Amount:     amount,
		AskFeeRate: askOrder.FeeRate,
		AskFee:     amount.Mul(askOrder.FeeRate),
		BidFeeRate: bidOrder.FeeRate,
		BidFee:     amount.Mul(bidOrder.FeeRate),
	}

	if err := tx.Table(tradeLog.TableName()).Create(&tradeLog).Error; err != nil {
		return nil, err
	}
	return &tradeLog, nil
}

func (s *SettleProcessor) updateAskOrderInfo(tx *gorm.DB, tradeLog *entities.TradeRecord, askOrder *entities.Order, remainderMarketOrderId string) error {

	askOrder.Fee = askOrder.Fee.Add(tradeLog.AskFee)
	askOrder.FinishedQty = askOrder.FinishedQty.Add(tradeLog.Quantity)
	askOrder.FinishedAmount = askOrder.FinishedAmount.Add(tradeLog.Amount)
	askOrder.AvgPrice = askOrder.FinishedAmount.Div(askOrder.FinishedQty)
	//初始状态为部分成交
	askOrder.Status = models_types.OrderStatusPartialFill

	if askOrder.OrderType == matching_types.OrderTypeLimit {
		cmp := askOrder.FinishedQty.Cmp(askOrder.Quantity)
		if cmp > 0 {
			return fmt.Errorf("invalid ask order finished qty")
		}
		if cmp == 0 {
			askOrder.Status = models_types.OrderStatusFilled
		}

		if err := tx.Table(askOrder.TableName()).Where("order_id=?", askOrder.OrderId).Updates(askOrder).Error; err != nil {
			return err
		}

		err := tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", askOrder.OrderId).Updates(askOrder).Error
		if err != nil {
			return err
		}

		if askOrder.Status == models_types.OrderStatusFilled {
			err := tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", askOrder.OrderId).Delete(askOrder).Error
			if err != nil {
				return err
			}
		}
	} else {
		//市价单结算
		if askOrder.Quantity.Equal(askOrder.FinishedQty) || askOrder.Amount.Equal(askOrder.FinishedAmount) {
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

func (s *SettleProcessor) updateBidOrderInfo(tx *gorm.DB, tradeLog *entities.TradeRecord, bidOrder *entities.Order, remainderMarketOrderId string) error {

	bidOrder.Fee = bidOrder.Fee.Add(tradeLog.BidFee)
	bidOrder.FinishedQty = bidOrder.FinishedQty.Add(tradeLog.Quantity)
	bidOrder.FinishedAmount = bidOrder.FinishedAmount.Add(tradeLog.Amount)
	bidOrder.AvgPrice = bidOrder.FinishedAmount.Div(bidOrder.FinishedQty)
	//初始状态为部分成交
	bidOrder.Status = models_types.OrderStatusPartialFill

	if bidOrder.OrderType == matching_types.OrderTypeLimit {
		cmp := bidOrder.FinishedQty.Cmp(bidOrder.Quantity)
		if cmp > 0 {
			return fmt.Errorf("invalid bid order finished qty")
		}
		if cmp == 0 {
			bidOrder.Status = models_types.OrderStatusFilled
		}

		err := tx.Table(bidOrder.TableName()).Where("order_id=?", bidOrder.OrderId).Updates(bidOrder).Error
		if err != nil {
			return err
		}

		err = tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", bidOrder.OrderId).Updates(bidOrder).Error
		if err != nil {
			return err
		}

		if bidOrder.Status == models_types.OrderStatusFilled {
			err := tx.Table(new(entities.UnfinishedOrder).TableName()).Where("order_id=?", bidOrder.OrderId).Delete(bidOrder).Error
			if err != nil {
				return err
			}
		}
	} else {
		//市价单结算
		if bidOrder.Quantity.Equal(bidOrder.FinishedQty) || bidOrder.Amount.Equal(bidOrder.FinishedAmount) {
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
	tradeLog *entities.TradeRecord,
	ask, bid *entities.Order,
	product *entities.Product,
) error {

	//结算被交易物品
	// 1.解冻卖家的冻结数量
	// 2.将解冻的数量转移给买家
	err := s.userAssetRepo.UnFreeze(tx, tradeLog.Ask, ask.UserId, product.Target.Symbol, tradeLog.Quantity)
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery target variety unFreeze: %v %s", tradeLog, err.Error())
		return err
	}
	err = s.userAssetRepo.TransferWithTx(tx, tradeLog.TradeId, ask.UserId, bid.UserId,
		product.Target.Symbol, tradeLog.Quantity)
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery target variety transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//结算本位币
	// 1.解冻买家冻结的金额
	// 2.将买家解冻的金额转入卖家账户
	amount := tradeLog.Amount
	err = s.userAssetRepo.UnFreeze(tx, tradeLog.Bid, bid.UserId, product.Base.Symbol, amount.Add(tradeLog.BidFee))
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery base variety unFreeze: %v %s", tradeLog, err.Error())
		return err
	}

	err = s.userAssetRepo.TransferWithTx(tx, tradeLog.TradeId, bid.UserId, ask.UserId, product.Base.Symbol, amount)
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery base variety transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//ask手续费收入到系统账号里
	err = s.userAssetRepo.TransferWithTx(tx, tradeLog.TradeId, ask.UserId, entities.SYSTEM_USER_FEE, product.Base.Symbol, tradeLog.AskFee)
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery ask fee transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//bid的手续费收入到系统账号里
	err = s.userAssetRepo.TransferWithTx(tx, tradeLog.TradeId, bid.UserId, entities.SYSTEM_USER_FEE, product.Base.Symbol, tradeLog.BidFee)
	if err != nil {
		s.logger.Sugar().Errorf("orderDelivery bid fee transfer: %v %s", tradeLog, err.Error())
		return err
	}

	//订单状态为已成交，则解冻该订单冻结的全部数量
	if ask.OrderType == matching_types.OrderTypeMarket && ask.Status == models_types.OrderStatusFilled {
		err = s.userAssetRepo.UnFreeze(tx, tradeLog.Ask, ask.UserId, product.Target.Symbol, decimal.Zero)
		if err != nil {
			s.logger.Sugar().Errorf("orderDelivery ask unFreeze: %v %s", tradeLog, err.Error())
			return err
		}
	}
	if bid.OrderType == matching_types.OrderTypeMarket && bid.Status == models_types.OrderStatusFilled {
		err = s.userAssetRepo.UnFreeze(tx, tradeLog.Bid, bid.UserId, product.Base.Symbol, decimal.Zero)
		if err != nil {
			s.logger.Sugar().Errorf("orderDelivery bid unFreeze: %v %s", tradeLog, err.Error())
			return err
		}
	}

	return nil
}

func generateTradeId(ask, bid string) string {
	date := time.Now().Format("060102")
	raw := fmt.Sprintf("%s%s", ask, bid)

	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", raw)))
	hashed := fmt.Sprintf("%x", hash.Sum(nil))
	return fmt.Sprintf("T%s%s", date, hashed[0:17])
}
