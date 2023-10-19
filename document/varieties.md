#### 所有交易品种


##### 1) 请求地址

>{haobase_host}/api/v1/base/trading/varieties

##### 2) 调用方式：HTTP GET

##### 3) 接口描述：

> 获取所有交易品种

##### 4) 请求参数:

##### GET参数:
|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|



##### 5) 请求返回结果:

```
{
    "data": [
        {
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
        {
            "symbol": "eurusd",
            "name": "欧美",
            "target_symbol_id": 2,
            "base_symbol_id": 1,
            "price_precision": 5,
            "qty_precision": 2,
            "allow_min_qty": "0.01",
            "allow_max_qty": "0",
            "allow_min_amount": "1",
            "allow_max_amount": "0",
            "fee_rate": "0.001",
            "update_at": 1697263181,
            "target": {
                "id": 2,
                "symbol": "eur",
                "name": "欧元",
                "show_precision": 2,
                "min_precision": 8,
                "base": false,
                "update_at": "2023-10-14T13:59:41+08:00"
            },
            "base": {
                "id": 1,
                "symbol": "usd",
                "name": "美元",
                "show_precision": 2,
                "min_precision": 8,
                "base": false,
                "update_at": "2023-10-14T13:59:41+08:00"
            }
        }
    ],
    "ok": 1
}
```


##### 6) 请求返回结果参数说明:
|字段名称       |类型            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|

  
##### END  
  

