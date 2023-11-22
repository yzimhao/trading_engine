#### Websocket 行情推送

  - 订阅地址 ws://{haoquote_host}/quote/ws
  - 用户可以侦听/订阅数个数据流

```go
type WebSocketMsgType string

const (
	MsgDepth       WebSocketMsgType = "depth.{symbol}"
	MsgTrade       WebSocketMsgType = "trade.{symbol}"
	MsgLatestPrice WebSocketMsgType = "price.{symbol}"
	MsgMarketKLine WebSocketMsgType = "kline.{period}.{symbol}"
	MsgMarket24H   WebSocketMsgType = "market.24h.{symbol}"
	MsgOrderCancel WebSocketMsgType = "order.cancel.{symbol}"
	MsgToken       WebSocketMsgType = "token.{token}"
	MsgUser        WebSocketMsgType = "_user.{user_id}" //特殊的类型，通过后端程序设置的属性
)
```
> 详细参考该文件： https://github.com/yzimhao/trading_engine/blob/master/types/websocket_msg.go


##### 订阅消息格式
```
{
    "sub":
    [
        "depth.usdjpy",
        "trade.usdjpy",
        "price.usdjpy",
        "kline.m1.usdjpy",
        "market.24h.usdjpy",
        "token." + token, //和用户相关的一些消息订阅，
    ]
}
```

##### 取消订阅消息格式
```
{
    "unsub":
    [
        "trade.usdjpy",
        "market.24h.usdjpy"
    ]
}
```