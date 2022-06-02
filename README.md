#### 交易系统撮合引擎
  一款使用golang编写的，采用优先级队列实现的订单排序、订单撮合、输出委托深度等的开源程序。
  <a href="http://132.226.14.192:8080/demo">在线demo</a> (免费的oracle cloud机器，望大家悠着点儿测试，感谢)
#### 功能列表
  - [x] 输出委托深度
  - [x] 订单撮合  
  - [x] 支持限价单
  - [x] 取消订单
  - [x] 统一数字小数位
  - [ ] 支持市价单


####
```
  go get github.com/yzimhao/trading_engine
```

#### 接入流程
```
  var btcusdt *trading_engine.TradePair
  btcusdt = trading_engine.NewTradePair("BTC_USDT", 2, 6)

  //买卖订单号最好做一个区分，方便识别订单
  //卖单 
  orderId = fmt.Sprintf("a-%s", orderId)
  item := trading_engine.NewAskItem(orderId, string2decimal(price), string2decimal(quantity), time.Now().Unix())
  btcusdt.PushNewOrder(trading_engine.OrderSideSell, item)

  //买单
  orderId = fmt.Sprintf("b-%s", orderId)
  item := trading_engine.NewBidItem(orderId, string2decimal(price), string2decimal(quantity), time.Now().Unix())
  btcusdt.PushNewOrder(trading_engine.OrderSideBuy, item)

  //获取深度, 参数为深度获取的个数 ["1.0001", "19960"] => [价格，数量]
  ask := btcusdt.GetAskDepth(10)
  bid := btcusdt.GetBidDepth(10)


  //买卖双方价格成交后会chan通知，监听如下
  for {
    if log, ok := <-btcusdt.ChTradeResult; ok {
      //其他通知，通知结算逻辑...
      ...
    }
  }

  //取消订单
  if strings.HasPrefix(orderId, "a-") {
    btcusdt.CancelOrder(trading_engine.OrderSideSell, orderId)
  } else {
    btcusdt.CancelOrder(trading_engine.OrderSideBuy, orderId)
  }

```  



#### example
  <a href="example">使用案例</a>