package quote

import (
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"quote",
	fx.Invoke(run),
)

func run(router *provider.Router, logger *zap.Logger) {
	newQuoteApi(router, logger)
}
