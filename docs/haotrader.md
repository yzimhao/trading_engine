## 下载地址
  Haotrader <a href="https://github.com/yzimhao/trading_engine/releases/latest">最新版本</a>

## 流程
  ![image](https://github.com/yzimhao/trading_engine/blob/master/docs/images/haotrader.png?raw=true)

## 交互数据结构
  redis数据key名称里的变量{prefix}、{symbol} 是在[config.toml](https://github.com/yzimhao/trading_engine/blob/master/cmd/config.toml)文件中配置
  

 * ### {prefix}.order.list.{symbol}

  推送新订单到撮合系统使用的队列,redis队列中的json数据
  1. 限价订单
  ```
{
    "order_id": "a-bfa82033-91bf-4c90-848e-221ca391815e", 
    "side": "ask", 
    "order_type":"limit", 
    "price": "1.424292", 
    "qty": "1.23", 
    "at": 1694534512
}
  ```
  2. 市价订单
      1. 按成交量卖出
      ```
      {
          "order_id": "a-bfa82033-91bf-4c90-848e-221ca391815e", 
          "side": "ask", 
          "order_type":"market_qty", 
          "qty": "100",
          "at": 1694534512
      }
      ```
      > 例：usdjpy交易对，用户需要卖出持有的100个usd
      2. 按成交量买入
      ```
      {
          "order_id": "a-bfa82033-91bf-4c90-848e-221ca391815e", 
          "side": "bid", 
          "order_type":"market_qty", 
          "qty": "100",
          "max_amount": "2387.8901",
          "at": 1694534512
      }
      ```
      > 例：usdjpy交易对，用户持有2387.8901个jpy，需要市价成交100个usd，最终买入的usd个数小于等于100。
      3. 按金额卖出
      ```
      {
          "order_id": "a-bfa82033-91bf-4c90-848e-221ca391815e", 
          "side": "ask", 
          "order_type":"market_amount", 
          "max_qty": "100",
          "amount": "2000",
          "at": 1694534512
      }
      ```
      > 例：用户最大持有100个usd，需要市价卖出价值2000个jpy的usd，最终卖出的usd个数小于等于100
      4. 按金额买入
      ```
      {
          "order_id": "a-bfa82033-91bf-4c90-848e-221ca391815e", 
          "side": "bid", 
          "order_type":"market_amount", 
          "amount": "2000",
          "at": 1694534512
      }
      ```
      > 例：需要按市价买入价值2000个jpy的usd
  > side只能传递ask/bid


 * ### {prefix}trade.result.{symbol}

  撮合成功队列，json数据
  ```
{
    "ask": "a-bfa82033-91bf-4c90-848e-221ca391815e", 
    "bid": "b-d368358c-7b90-4d94-85c2-087a6d0ddbb6", 
    "trade_amount": "1.751879", 
    "trade_price": "1.424292", 
    "trade_quantity": "1.23", 
    "trade_time": 1694535017388092
}
  ```

 * ### {prefix}need.cancel.{symbol}
  取消订单队列，json数据
  ```
{
    "side": "ask",
    "order_id": "a-bfa82033-91bf-4c90-848e-221ca391815e"
}
  ```

 * ### {prefix}cancel.result.{symbol}
  撮合系统取消成功通知队列，json数据
  ```
{
    "order_id": "a-bfa82033-91bf-4c90-848e-221ca391815e",
    "cancel": "success",
}
  ```


## Linux/Macos使用
  下载最新版本程序后，解压缩包，进入文件夹。
  ```
  ➜  mv config.toml_sample config.toml
  #修改配置config.toml文件后
  ➜  ./haotrader 
time="2023-09-15 16:00:09" level="info" msg="当前运行在dev模式下，生产环境时main.mode请务必成prod"
time="2023-09-15 16:00:09" level="info" msg="启动撮合程序成功! 如需帮助请参考: https://github.com/yzimhao/trading_engine"
time="2023-09-15 16:00:09" level="info" msg="正在恢复[usd]数据，共0条"
time="2023-09-15 16:00:09" level="info" msg="[usd]数据恢复 已完成"
time="2023-09-15 16:00:09" level="info" msg="正在恢复[usdjpy]数据，共0条"
time="2023-09-15 16:00:09" level="info" msg="[usdjpy]数据恢复 已完成"
time="2023-09-15 16:00:09" level="info" msg="正在监听redis队列: order.list.usdjpy"
time="2023-09-15 16:00:09" level="info" msg="正在监听redis队列: order.list.usd"
time="2023-09-15 16:00:09" level="info" msg="正在监听redis队列: need.cancel.usdjpy"
time="2023-09-15 16:00:09" level="info" msg="正在监听redis队列: need.cancel.usd"
time="2023-09-15 16:00:09" level="info" msg="http服务监听: 0.0.0.0:8081"
  ```


## Windows
  请参考Linux的步骤


## 调试模拟数据
  在dev模式下，可以在命令行下使用下面的参数批量插入数据,快速调试。
  ```shell
  # 向usd交易对，插入限价卖单 价格范围是1.xxxx 每单数量2，循环1000单
  ./haotrader test -s usd --side=ask -p 1 -q 2 -n 1000
  # 向usd交易对，插入限价买单 价格范围是1.xxxx 每单数量2，循环1000单
  ./haotrader test -s usd --side=bid -p 1 -q 2 -n 1000
  ```
  更多命令支持，请
  ```
  ./haotrader help
  ```
  