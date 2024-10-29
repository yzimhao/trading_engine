package settlement

import (
	"github.com/duolacloud/crud-core/cache"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type settlement struct {
	db     *gorm.DB
	logger *zap.Logger
	cache  cache.Cache
}

type inContext struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Cache  cache.Cache
}

func NewSettlement(in inContext) *settlement {
	return &settlement{
		db:     in.DB,
		logger: in.Logger,
		cache:  in.Cache,
	}
}

func (s *settlement) Run(tradeResult matching_types.TradeResult) {
	//TODO 自动创建结算需要的表
}

func (s *settlement) flow() error {

	return nil
}

func (s *settlement) checkOrder() error {
	return nil
}
func (s *settlement) writeTradeLog() error {
	return nil
}

func (s *settlement) updateAskOrderInfo() error {
	return nil
}

func (s *settlement) updateBidOrderInfo() error {
	return nil
}

func (s *settlement) orderDelivery() error {
	return nil
}
