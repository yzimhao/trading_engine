
<p align="center">
    <img src="https://img.shields.io/github/stars/yzimhao/trading_engine?style=social">
    <img src="https://img.shields.io/github/forks/yzimhao/trading_engine?style=social">
	<img src="https://img.shields.io/github/issues/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/repo-size/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/license/yzimhao/trading_engine">
</p>




#### 撮合
```
    go get github.com/yzimhao/trading_engine/v2/pkg/matching
```
#### example
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