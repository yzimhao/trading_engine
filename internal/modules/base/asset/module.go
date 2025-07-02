package asset

import "go.uber.org/fx"

var Module = fx.Module(
	"base.asset",
	fx.Invoke(newAssetModule),
)
