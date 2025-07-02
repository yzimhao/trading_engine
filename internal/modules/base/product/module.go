package product

import "go.uber.org/fx"

var Module = fx.Module(
	"base.product",
	fx.Invoke(
		newProductModule,
	),
)
