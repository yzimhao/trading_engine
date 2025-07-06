package quote

import (
	"github.com/duolacloud/crud-core/cache"
	"github.com/yzimhao/trading_engine/v2/internal/di/provider"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"quote",
	fx.Invoke(run),
)

func run(router *provider.Router, logger *zap.Logger, c cache.Cache) {
	newQuoteApi(router, logger, c)
}
