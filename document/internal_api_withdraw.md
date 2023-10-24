#### 提现

##### 1) 请求地址

>{haobase_host}/api/v1/internal/withdraw

##### 2) 调用方式：HTTP POST

##### 3) 接口描述：

> 内部接口，用户提现资产


##### 4) 请求参数:


##### Body:
> Content-Type: application/json; charset=utf-8
```
{
    "user_id": "100001", #提现用户
    "order_id": "w0003", #提现订单号，唯一
    "symbol": "usd", 
    "amount": "100" #提现的金额 必须大于0
}
```


##### 5) 请求返回结果:

```
{
    "data": "",
    "ok": 1
}
```


##### 6) 请求返回结果参数说明:

  
##### END  
  
