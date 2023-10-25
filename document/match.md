#### 只接入撮合引擎

```shell
  go get github.com/yzimhao/trading_engine/trading_core
```

```go
  var object *trading_core.TradePair
  //初始化交易对，需要设置价格、数量的小数点位数，
  //需要将数字格式化字符串对外展示的时候，用到这两个小数位，统一数字长度
  object = trading_core.NewTradePair("symbol", 2, 6)

  //买卖订单号最好做一个区分，方便识别订单
  
  //卖单
  uniq := fmt.Sprintf("a-%s", orderId)
  createTime := time.Now().Unix()
  //限价卖单
  item := trading_core.NewAskLimitItem(uniq, price, quantity, createTime)
  object.PushNewOrder(item)
  //市价-按数量卖出
  item = trading_core.NewAskMarketQtyItem(uniq, quantity, createTime)
  object.PushNewOrder(item)
  //市价-按金额卖出,需要用户持有的该资产最大数量
  item = trading_core.NewAskMarketAmountItem(uniq, amount, maxQty, createTime)
  object.PushNewOrder(item)


  //买单
  uniq := fmt.Sprintf("b-%s", orderId)
  createTime := time.Now().Unix()
  //限价买单
  item := trading_core.NewBidLimitItem(uniq, price, quantity, createTime)
  object.PushNewOrder(item)
  //市价-按数量买单,需要用户可用资金来限制最大买入量
  item = trading_core.NewBidMarketQtyItem(uniq, quantity, maxAmount, createTime)
  object.PushNewOrder(item)
  //市价-按金额买单
  item = trading_core.NewBidMarketAmountItem(uniq, amount, createTime)
  object.PushNewOrder(item)


  //取消订单, 该操作会将uniq订单号从队列中移除，然后发出一个chan通知在ChCancelResult
  //业务代码可以通过监听取消通知，去做撤单逻辑相关的操作
  if strings.HasPrefix(orderId, "a-") {
      object.CancelOrder(trading_core.OrderSideSell, uniq)
  } else {
      object.CancelOrder(trading_core.OrderSideBuy, uniq)
  }


  //获取深度, 参数为深度获取的个数 ["1.0001", "19960"] => [价格，数量]
  ask := object.GetAskDepth(10)
  bid := object.GetBidDepth(10)


  //撮合系统有chan通知，监听如下
  for {
    select{
        case tradelog := <-object.ChTradeResult:
            //撮合成功，买卖双方订单信息，成交价格、数量等
            //通知结算逻辑...
            ...
        case orderId := <- object.ChCancelResult:
            //被取消的订单id, 确认队列里面没有了 会有这个通知
            ...
        default:
            time.Sleep(time.Duration(50) * time.Millisecond)
    }
    
  }

```  

