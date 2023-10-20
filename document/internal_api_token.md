#### 设置用户登陆Token

##### 1) 请求地址

>{haobase_host}/api/v1/internal/settoken

##### 2) 调用方式：HTTP POST

##### 3) 接口描述：

> 一个服务端内部接口，用于已有自己的用户系统，对接本系统的交易服务。


##### 4) 请求参数:



##### Body:
> Content-Type: application/json; charset=utf-8
```
{
    "user_id": "10012",
    "token": "user_token_value", //前端传递到交易系统的token值，能和user_id唯一对应
    "ttl": 600 //token过期的时间 单位秒
}
```

|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|
|order_id|string|Y|需要取消的订单号||



##### 5) 请求返回结果:

```
{
    "data": "",
    "ok": 1
}
```


##### 6) 请求返回结果参数说明:

  
##### END  
  
