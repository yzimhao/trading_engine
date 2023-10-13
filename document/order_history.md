#### 历史订单

##### 1) 请求地址

>{haobase_host}/api/v1/base/order/history

##### 2) 调用方式：HTTP GET

##### 3) 接口描述：

> 历史订单

##### 4) 请求参数:

##### GET参数:
|字段名称       |类型            |必填            |字段说明         |备注     |
| -------------|:--------------:|:--------------:|:--------------:|:------:|
|symbol|string|Y|||
|order_id|string|N|||
|start_time|string|N|||
|end_time|string|N||默认 500; 最大 1000|
|limit|int|N|||

注意：
 * 如设置 order_id , 订单量将 >= order_id
 * 如果设置 start_time 和 end_time, order_id 就不需要设置。


##### 5) 请求返回结果:

```

```


##### 6) 请求返回结果参数说明:


  
##### END  
  
