package example

import (
	"go.uber.org/fx"
)

var Module = fx.Module(
	"example",
	fx.Invoke(newExample),
)
