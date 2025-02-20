package modules

import (
	"github.com/yzimhao/trading_engine/v2/internal/modules/matching"
	"github.com/yzimhao/trading_engine/v2/internal/modules/quote"
	"github.com/yzimhao/trading_engine/v2/internal/modules/settlement"
	"go.uber.org/fx"
)

var Load = fx.Module(
	"modules",
	settlement.Module,
	matching.Module,
	quote.Module,
)
