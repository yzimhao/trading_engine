
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
    classDef done fill:#2ecc71,stroke:#27ae60,color:#fff  
    classDef pending fill:#bdc3c7,stroke:#95a5a6,color:#333

    A[交易所<br>Exchange]

    A --> aa[基础数据<br>Base]
    aa --> aa1[资产种类<br>Asset]
    aa --> aa2[交易品种<br>Product]

    A --> AA[用户中心<br>UserCenter]
    AA --> AAa[用户资产<br>UserAssets]
    class AAa done
    AA--> AAb[用户订单<br>UserOrders]
    class AAb done
    A --> AC[交易核心<br>TradingCore]
    A --> AC3[行情系统<br>Quote]
    class AC3 pending
    A --> AD[消息推送<br>Notification]
    class AD pending
    A --> AZ[...<br>Other]

    AC --> AC1[撮合引擎<br>MatchingEngine]
    class AC1 done
    AC --> AC2[订单结算<br>Settlement]
    class AC2 done
    
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