package quote

import "go.uber.org/fx"

var Module = fx.Module(
	"quote",
	fx.Provide(NewQuote),
	fx.Invoke(startupQuote),
)

func startupQuote(quote *Quote) {
	quote.Subscribe()
}
