package tradingcore

import (
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/matching"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/orderlock"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/settlement"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tradingcore",
	fx.Provide(
		orderlock.NewOrderLock,
	),
	settlement.Module,
	matching.Module,
)
