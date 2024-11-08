package quote

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewQuote),
	fx.Invoke(startupQuote),
)

func startupQuote(quote *Quote) {
	quote.Subscribe()
}
