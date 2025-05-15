
<p align="center">
    <img src="https://img.shields.io/github/stars/yzimhao/trading_engine?style=social">
    <img src="https://img.shields.io/github/forks/yzimhao/trading_engine?style=social">
	<img src="https://img.shields.io/github/issues/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/repo-size/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/license/yzimhao/trading_engine">
</p>


#### 贡献代码
  全新重构，希望更多开源爱好者能加入共同开发练手，共同进步，请联系下方wx



### 系统架构
```mermaid
graph TD
    A[交易所<br>Exchange]
    A --> AA[用户资产<br>UserAssets]
    A --> AB[用户订单<br>UserOrders]
    A --> AC[交易核心<br>TradingCore]
    A --> AD[消息推送<br>Notification]

    AC --> AC1[撮合引擎<br>MatchingEngine]
    AC --> AC2[订单结算<br>Settlement]
    AC --> AC5[行情系统<br>MarketData]
```

---

#### 撮合引擎
> 只集成撮合引擎部分

```
    go get github.com/yzimhao/trading_engine/v2/pkg/matching
```

```go
    ctx := context.Background()
    opts := []matching.Option{
        matching.WithPriceDecimals(2),
        matching.WithQuantityDecimals(2),
    }
    btcusdt = matching.NewEngine(ctx, "btcusdt", opts...)

    //添加订单
    btcusdt.AddItem(...)
    //移除订单
    btcusdt.RemoveItem(...)

    //监听事件结果
    btcusdt.OnTradeResult(func(result types.TradeResult) {
        //TODO
    })
    btcusdt.OnRemoveResult(func(result types.RemoveResult) {
        //TODO
    })

    //获取深度
    btcusdt.GetAskOrderBook(10) // [][2]string [["1.01","4.00"],["1.10","2.00"]]
    btcusdt.GetBidOrderBook(10)

```


 
  #### 交流
<img src="https://github.com/yzimhao/trading_engine/blob/main/documents/images/wechat.jpg?raw=true" width = "150"/>

  #### Star History

[![Star History Chart](https://api.star-history.com/svg?repos=yzimhao/trading_engine&type=Date)](https://star-history.com/#yzimhao/trading_engine&Date)

![Visitor's Count](https://profile-counter.glitch.me/yzimhao_trading_engine/count.svg)