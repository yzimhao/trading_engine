## HaoTrader
  
<p align="center">
    <img src="https://img.shields.io/github/stars/yzimhao/trading_engine?style=social">
    <img src="https://img.shields.io/github/forks/yzimhao/trading_engine?style=social">
	<img src="https://img.shields.io/github/issues/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/repo-size/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/license/yzimhao/trading_engine">
</p>

  HaoTrader适用于各种金融证券交易场景。拥有高性能的订单撮合、实时结算、行情计算、实时推送等功能。支持数据持久化，故障重启快速恢复数据。
  
  配置灵活，允许用户根据自身需求自定义配置各模块独立运行。
  

<!-- ## 流程
  ![image](https://github.com/yzimhao/trading_engine/blob/master/docs/images/haotrader.png?raw=true) -->

## 演示
  <a href="http://144.91.108.90:20001/" target="_blank">在线体验</a> 
  > 感谢[9cat](https://github.com/9cat)大佬提供免费测试服务器 



## HaoTrader系统包含模块
  - [x] 交易品种 --> [haobase]
  - [x] 资产模块 --> [haobase]
  - [x] 订单模块 --> [haobase]
  - [x] 撮合模块 --> [haomatch]
  - [x] 结算模块 --> [haobase]
  - [x] 行情模块 --> [haoquote]
  


## 引入包接入
```
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



## 参考
  <a href="https://www.liaoxuefeng.com/article/1185272483766752" target="_blank">证券交易系统设计与开发</a>

## 需求讨论
   <img src="https://github.com/yzimhao/trading_engine/blob/master/document/images/wechat.jpg?raw=true" width = "150"/>

## 声明
  - 本项目仅供参考和学习之用，不建议将其用于生产环境或重要交易场景。
