#### Websocket 行情推送

  - 订阅地址 ws://{haoquote_host}/quote/ws
  - 用户可以侦听/订阅数个数据流

##### 订阅消息格式
```
{
    "sub":
    [
        "depth.usdjpy",
        "tradelog.usdjpy",
        "latest_price.usdjpy",
        "kline.m1.usdjpy",
        "market.24h.usdjpy"
    ]
}
```

##### 取消订阅消息格式
```
{
    "unsub":
    [
        "tradelog.usdjpy",
        "market.24h.usdjpy"
    ]
}
```