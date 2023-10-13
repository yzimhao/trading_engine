package orders

// func trade_market_buy(user_id int64, trade_pair_id int, qty, max_amount, fee_rate string) (order_id string, err error) {
// 	order := TradeOrder{
// 		OrderId:       order_id_by_side(OrderSideBid),
// 		TradingPair:   trade_pair_id,
// 		OrderSide:     OrderSideBid,
// 		OrderType:     OrderTypeMarket,
// 		UserId:        user_id,
// 		Price:         "-1",
// 		Quantity:      qty,
// 		UnfinishedQty: qty,
// 		TotalAmount:   max_amount,
// 		TradeAmount:   "0",
// 		FeeRate:       fee_rate,
// 		Status:        OrderStatusNew,
// 	}

// 	sess := db_engine.NewSession()
// 	defer sess.Close()

// 	//todo 开启事务

// 	tp := base.GetTradePairById(trade_pair_id)
// 	if tp == nil {
// 		return "", fmt.Errorf("invalid trade pair id")
// 	}

// 	//冻结相应资产
// 	_, err = assets.FreeeBalance(sess, user_id, tp.BaseSymbolId, max_amount, order.OrderId, "trade order")
// 	if err != nil {
// 		return "", err
// 	}

// 	//save order
// 	_, err = sess.Table(new(TradeOrder)).Insert(&order)
// 	if err != nil {
// 		return "", err
// 	}
// 	return order.OrderId, nil
// }

// func trade_market_sell(user_id int64, trade_pair_id int, qty, max_amount, fee_rate string) (order_id string, err error) {
// 	order := TradeOrder{
// 		OrderId:       order_id_by_side(OrderSideAsk),
// 		TradingPair:   trade_pair_id,
// 		OrderSide:     OrderSideAsk,
// 		OrderType:     OrderTypeMarket,
// 		UserId:        user_id,
// 		Price:         "-1",
// 		Quantity:      qty,
// 		UnfinishedQty: qty,
// 		TotalAmount:   max_amount,
// 		TradeAmount:   "0",
// 		FeeRate:       fee_rate,
// 		Status:        OrderStatusNew,
// 	}

// 	sess := db_engine.NewSession()
// 	defer sess.Close()

// 	tp := base.GetTradePairById(trade_pair_id)
// 	if tp == nil {
// 		return "", fmt.Errorf("invalid trade pair id")
// 	}

// 	//冻结相应资产
// 	_, err = assets.FreeeBalance(sess, user_id, tp.TradeSymbolId, qty, order.OrderId, "trade order")
// 	if err != nil {
// 		return "", err
// 	}

// 	//save order
// 	_, err = sess.Table(new(TradeOrder)).Insert(&order)
// 	if err != nil {
// 		return "", err
// 	}
// 	return order.OrderId, nil
// }
