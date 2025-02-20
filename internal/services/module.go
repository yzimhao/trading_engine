package services

import "go.uber.org/fx"

var Module = fx.Module(
	"services",
	fx.Provide(),
)
