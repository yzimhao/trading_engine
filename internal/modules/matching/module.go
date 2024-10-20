package matching

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewMatchingSubscriber),
	fx.Invoke(startupSubscriber),
)

func startupSubscriber(subscriber *MatchingSubscriber) {
	subscriber.Subscribe()
}
