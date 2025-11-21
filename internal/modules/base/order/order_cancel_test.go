package order

import (
	"encoding/json"
	"testing"

	models_types "github.com/yzimhao/trading_engine/v2/internal/types"
	matching_types "github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

func TestEventNotifyCancelOrder_Marshal(t *testing.T) {
	data := models_types.EventNotifyCancelOrder{
		Symbol:    "btcusdt",
		OrderSide: matching_types.OrderSideBuy,
		OrderId:   "B123",
		Type:      matching_types.RemoveItemTypeByUser,
	}

	b, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(b, &parsed); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if parsed["symbol"] != "btcusdt" {
		t.Fatalf("expected symbol btcusdt got %v", parsed["symbol"])
	}
	if parsed["order_id"] != "B123" {
		t.Fatalf("expected order_id B123 got %v", parsed["order_id"])
	}

	if models_types.TOPIC_NOTIFY_ORDER_CANCEL != "notify_order_cancel" {
		t.Fatalf("topic constant changed: %s", models_types.TOPIC_NOTIFY_ORDER_CANCEL)
	}
}
