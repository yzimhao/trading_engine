package matching

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewMatching),
	fx.Invoke(startupMatching),
)

func startupMatching(matching *Matching) {
	//TODO init engine
	matching.InitEngine()
	matching.Subscribe()
}
