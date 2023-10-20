#### 当前挂单

##### 1) 请求地址

>{haobase_host}/api/v1/base/order/unfinshed

##### 2) 调用方式：HTTP GET

##### 3) 接口描述：

> 获取还未完全成交的订单列表

##### 4) 请求参数:
> 需要登陆
##### Header:
|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|
|Token|string|Y|||


##### GET参数:
|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|
|symbol|string|N||为空则返回所有交易对的挂单|
|limit|int|N||默认 500; 最大 1000|



##### 5) 请求返回结果:

```
{
    "data": [
        {
            "symbol": "usdjpy",
            "order_id": "A23102015240522027723",
            "order_side": "sell",
            "order_type": "limit",
            "price": "1.000",
            "quantity": "10.00",
            "finished_qty": "0.00",
            "finished_amount": "0.000",
            "status": "new",
            "create_time": 1697786645236014000
        },
        {
            "symbol": "usdjpy",
            "order_id": "A23102015240331493947",
            "order_side": "sell",
            "order_type": "limit",
            "price": "1.000",
            "quantity": "10.00",
            "finished_qty": "0.00",
            "finished_amount": "0.000",
            "status": "new",
            "create_time": 1697786643331832000
        }
    ],
    "ok": 1
}
```


##### 6) 请求返回结果参数说明:


  
##### END  
  
