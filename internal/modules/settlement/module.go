package settlement

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewSettlementSubscriber,
		NewCancelOrderSubscriber,
	),
	fx.Invoke(startupSubscriber),
)

func startupSubscriber(subscriber *SettlementSubscriber, cancelOrder *CancelOrderSubscriber) {
	subscriber.Subscribe()
	cancelOrder.Subscribe()
}
