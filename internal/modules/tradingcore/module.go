package tradingcore

import (
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/matching"
	"github.com/yzimhao/trading_engine/v2/internal/modules/tradingcore/settlement"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"tradingcore",
	settlement.Module,
	matching.Module,
)
