package matching_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/pkg/matching"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

func BenchmarkMatchingEngine(b *testing.B) {
	ctx := context.Background()
	engine := matching.NewEngine(ctx, "BTCUSDT")

	initialOrders := 10000
	for i := 0; i < initialOrders; i++ {
		order := matching.NewAskLimitItem(fmt.Sprintf("%d", i), decimal.NewFromInt(10), decimal.NewFromInt(1), 1)
		engine.AddItem(order)
	}

	engine.OnTradeResult(func(result types.TradeResult) {
		b.Logf("trade result: %v", result)
	})

	b.ResetTimer()

	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			order := matching.NewBidLimitItem(fmt.Sprintf("%d", i), decimal.NewFromInt(10), decimal.NewFromInt(1), 1)
			engine.AddItem(order)
		}(i)
	}

	wg.Wait()
	b.StopTimer()

	b.Logf("Finished matching %d orders", b.N)
}
