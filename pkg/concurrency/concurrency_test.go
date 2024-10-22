package concurrency_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/yzimhao/trading_engine/v2/pkg/concurrency"
)

func TestExecutor(t *testing.T) {
	executor := concurrency.NewExecutor(5)

	for i := 0; i < 20; i++ {
		executor.Execute(func() any {
			time.Sleep(time.Second)
			fmt.Printf("task %d done\n", i)
			return i
		})
	}

	results := executor.Run()

	for _, result := range results {
		fmt.Println(result)
	}
}

func BenchmarkExecutor(b *testing.B) {
	for i := 0; i < b.N; i++ {
		executor := concurrency.NewExecutor(5)
		for j := 0; j < 1000; j++ {
			executor.Execute(func() any {
				return nil
			})
		}
		executor.Run()
	}
}
