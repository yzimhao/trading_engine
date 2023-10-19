#### 交易品种规格


##### 1) 请求地址

>{haobase_host}/api/v1/base/varieties/config

##### 2) 调用方式：HTTP GET

##### 3) 接口描述：

> 获取指定交易品种规格信息

##### 4) 请求参数:

##### GET参数:
|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|
|symbol|string|N|交易对symbol||


##### 5) 请求返回结果:

```
{
    "data":{
        "symbol": "usdjpy",
        "name": "美日",
        "target_symbol_id": 1,
        "base_symbol_id": 3,
        "price_precision": 3,
        "qty_precision": 2,
        "allow_min_qty": "0.01",
        "allow_max_qty": "0",
        "allow_min_amount": "1",
        "allow_max_amount": "0",
        "fee_rate": "0.005",
        "update_at": 1697263181,
        "target": {
            "id": 1,
            "symbol": "usd",
            "name": "美元",
            "show_precision": 2,
            "min_precision": 8,
            "base": false,
            "update_at": "2023-10-14T13:59:41+08:00"
        },
        "base": {
            "id": 3,
            "symbol": "jpy",
            "name": "日元",
            "show_precision": 2,
            "min_precision": 8,
            "base": false,
            "update_at": "2023-10-14T13:59:41+08:00"
        }
    },
    "ok": 1
}
```


##### 6) 请求返回结果参数说明:
|字段名称       |类型            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|

  
##### END  
  

