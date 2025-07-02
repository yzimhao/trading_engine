package usercenter

import (
	"github.com/yzimhao/trading_engine/v2/internal/modules/usercenter/assets"
	"github.com/yzimhao/trading_engine/v2/internal/modules/usercenter/orders"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"usercenter",
	assets.Module,
	orders.Module,
)
