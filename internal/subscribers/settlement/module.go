package settlement

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewSettlementSubscriber,
	),
	fx.Invoke(startupSubscriber),
)

func startupSubscriber(subscriber *SettlementSubscriber) {
	subscriber.Subscribe()
}
