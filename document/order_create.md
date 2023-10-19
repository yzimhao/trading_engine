#### 新订单

##### 1) 请求地址

>{haobase_host}/api/v1/base/order/create

##### 2) 调用方式：HTTP POST

##### 3) 接口描述：

> 创建新订单

##### 4) 请求参数:
> 需要登陆
##### Header:
|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|
|Token|string|Y|||

##### Body:
> Content-Type: application/json; charset=utf-8
```
{
    "symbol": "usdjpy",
    "side": "sell",
    "order_type": "limit",
    "price": "1",
    "qty": "1"
}
```

|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|
|symbol|string|Y|||
|side|string|Y|值：sell、buy||
|order_type|string|Y|limit/market||
|price|string|N|只有限价单有效||
|qty|string|N|限价单或按市价买卖指定数量时||
|amount|string|N|市价买卖指定金额时||



##### 5) 请求返回结果:

```
{
    "data": "A23101211561219972141",
    "ok": 1
}
```


##### 6) 请求返回结果参数说明:
|字段名称       |类型            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|
|data|string|下单成功订单号||

  
##### END  
  
