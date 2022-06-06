#### 交易系统撮合引擎
  一款使用golang编写的，采用优先级队列实现的订单排序、订单撮合、输出委托深度等的开源程序。

#### 适用于场景 
  买卖双方各自报价，按照价格优先、时间优先的顺序，对买卖双方进行撮合。该程序只实现了撮合系统这部分逻辑。
  
  一个完整的交易系统是由用户系统、账户系统、订单系统、撮合系统以及清算系统等子系统构成的。各个子系统相互配合，完成交易物报价交易。（可以参考 <a href="https://www.liaoxuefeng.com/article/1185272483766752" target="_blank">廖雪峰老师的文章</a>）

  

#### Demo
  <a href="http://132.226.14.192:8080/demo" target="_blank">在线体验</a> (免费的oracle cloud机器，望大家悠着点儿测试，感谢)
#### 功能列表
  - [x] 委托深度
  - [x] 订单撮合  
  - [x] 限价单
  - [x] 取消订单
  - [x] 市价单 按数量
  - [ ] 市价单 按金额


####
```
  go get github.com/yzimhao/trading_engine
```

#### 接入相关方法介绍
```go
  var btcusdt *trading_engine.TradePair
  //初始化交易对，需要设置价格、数量的小数点位数，
  //需要将数字格式化字符串对外展示的时候，用到这两个小数位，统一数字长度
  btcusdt = trading_engine.NewTradePair("BTC_USDT", 2, 6)

  //买卖订单号最好做一个区分，方便识别订单
  
  //卖单
  uniq := fmt.Sprintf("a-%s", orderId)
  createTime := time.Now().Unix()
  //限价卖单
  item := trading_engine.NewAskLimitItem(uniq, price, quantity, createTime)
  btcusdt.PushNewOrder(item)
  //市价-按数量卖出
  item = trading_engine.NewAskMarketQtyItem(uniq, quantity, createTime)
  btcusdt.PushNewOrder(item)
  //市价-按金额卖出
  item = trading_engine.NewAskMarketAmountItem(uniq, amount, createTime)
  btcusdt.PushNewOrder(item)


  //买单
  uniq := fmt.Sprintf("b-%s", orderId)
  createTime := time.Now().Unix()
  //限价买单
  item := trading_engine.NewBidLimitItem(uniq, price, quantity, createTime)
  btcusdt.PushNewOrder(item)
  //市价-按数量买单
  item = trading_engine.NewBidMarketQtyItem(uniq, quantity, createTime)
  btcusdt.PushNewOrder(item)
  //市价-按金额买单
  item = trading_engine.NewBidMarketAmountItem(uniq, amount, createTime)
  btcusdt.PushNewOrder(item)


  //取消订单, 该操作会将uniq订单号从队列中移除，然后发出一个chan通知在ChCancelResult
  //业务代码可以通过监听取消通知，去做撤单逻辑相关的操作
  if strings.HasPrefix(orderId, "a-") {
      btcusdt.CancelOrder(trading_engine.OrderSideSell, uniq)
  } else {
      btcusdt.CancelOrder(trading_engine.OrderSideBuy, uniq)
  }


  //获取深度, 参数为深度获取的个数 ["1.0001", "19960"] => [价格，数量]
  ask := btcusdt.GetAskDepth(10)
  bid := btcusdt.GetBidDepth(10)


  //撮合系统有chan通知，监听如下
  for {
    select{
        case tradelog := <-btcusdt.ChTradeResult:
            //撮合成功，买卖双方订单信息，成交价格、数量等
            //通知结算逻辑...
            ...
        case orderId := <- btcusdt.ChCancelResult:
            //被取消的订单id, 确认队列里面没有了 会有这个通知
            ...
        default:
            time.Sleep(time.Duration(50) * time.Millisecond)
    }
    
  }

```  



#### example
  <a href="example">使用案例</a>


#### 需求讨论联系
<img src="example/me.jpg" style="width:250px;">