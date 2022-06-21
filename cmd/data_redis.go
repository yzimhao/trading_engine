package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	te "github.com/yzimhao/trading_engine"
)

func newOrderFromRedis(wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	logrus.Info("start newOrderFromRedis")
	sub := rdb.Subscribe(ctx, fmt.Sprintf("push_new_order.%s", pair.Symbol))
	defer sub.Close()
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			logrus.Errorf("sub.ReceiveMessage error: %v", err)
			continue
		}
		logrus.Debugf("sub.ReceiveMessage: %v", msg)
		if msg.Channel == fmt.Sprintf("push_new_order.%s", pair.Symbol) {
			var item NewOrderMsgBody

			err := json.Unmarshal([]byte(msg.Payload), &item)
			if err != nil {
				logrus.Errorf("json.Unmarshal error: %v payload: %s", err, msg.Payload)
				continue
			}

			order := item.Order
			price_type := strings.ToLower(item.PriceType)
			order_type := strings.ToLower(order.Side)
			if price_type == "limit" {
				if order_type == "buy" {
					pair.ChNewOrder <- te.NewBidLimitItem(order.OrderId, str2decimal(order.Price), str2decimal(order.Quantity), str2Int64(order.CreateTime))
				} else if order_type == "sell" {
					pair.ChNewOrder <- te.NewAskLimitItem(order.OrderId, str2decimal(order.Price), str2decimal(order.Quantity), str2Int64(order.CreateTime))
				}
			} else if price_type == "market-qty" {
				if order_type == "buy" {
					maxHoldAmount := str2decimal(order.MaxHoldAmount)
					pair.ChNewOrder <- te.NewBidMarketQtyItem(order.OrderId, str2decimal(order.Quantity), maxHoldAmount, str2Int64(order.CreateTime))
				} else if order_type == "sell" {
					pair.ChNewOrder <- te.NewAskMarketQtyItem(order.OrderId, str2decimal(order.Quantity), str2Int64(order.CreateTime))
				}
			} else if price_type == "market-amount" {
				if order_type == "buy" {
					pair.ChNewOrder <- te.NewBidMarketAmountItem(order.OrderId, str2decimal(order.Amount), str2Int64(order.CreateTime))
				} else if order_type == "sell" {
					maxHoldQty := str2decimal(order.MaxHoldQty)
					pair.ChNewOrder <- te.NewAskMarketAmountItem(order.OrderId, str2decimal(order.Amount), maxHoldQty, str2Int64(order.CreateTime))
				}
			}
		}
	}
}

func cancelOrderFromRedis(wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	logrus.Info("cancelOrderFromRedis")
	sub := rdb.Subscribe(ctx, fmt.Sprintf("cancel_order.%s", pair.Symbol))
	defer sub.Close()
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			logrus.Errorf("sub.ReceiveMessage error: %v", err)
			continue
		}
		logrus.Debugf("sub.ReceiveMessage: %v", msg)

		if msg.Channel == fmt.Sprintf("cancel_order.%s", pair.Symbol) {
			var item CancelOrderMsgBody

			err := json.Unmarshal([]byte(msg.Payload), &item)
			if err != nil {
				logrus.Errorf("json.Unmarshal error: %v payload: %s", err, msg.Payload)
				continue
			}

			if item.Side == "buy" {
				pair.CancelOrder(te.OrderSideBuy, item.OrderId)
			} else if item.Side == "sell" {
				pair.CancelOrder(te.OrderSideSell, item.OrderId)
			}
		}
	}
}

func publishMsgToRedis(wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()

	for {
		select {
		case log, ok := <-pair.ChTradeResult:
			if ok {
				msg, _ := json.Marshal(log)
				rdb.Publish(ctx, fmt.Sprintf("trade_log.%s", pair.Symbol), msg)
			}
		case cancelOrderId := <-pair.ChCancelResult:
			rdb.Publish(ctx, fmt.Sprintf("cancel_result.%s", pair.Symbol), cancelOrderId)

		default:
			time.Sleep(time.Duration(100) * time.Millisecond)
		}

	}
}
