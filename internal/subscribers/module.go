package subscribers

import (
	"github.com/yzimhao/trading_engine/v2/internal/subscribers/settlement"
	"go.uber.org/fx"
)

var Module = fx.Options(
	settlement.Module,
)
