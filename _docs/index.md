## Haotrader交易撮合引擎
<p align="center">
    <img src="https://img.shields.io/github/stars/yzimhao/trading_engine?style=social">
    <img src="https://img.shields.io/github/forks/yzimhao/trading_engine?style=social">
	<img src="https://img.shields.io/github/issues/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/repo-size/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/license/yzimhao/trading_engine">
</p>


## 主要特点和功能：

  * 订单撮合引擎：haotrader提供了一个强大的订单撮合引擎，可以高效地处理买卖双方的报价，并根据价格和时间优先原则撮合订单。这使得haotrader适用于各种交易场景，包括股票、期货、数字货币等。

  * 优先级队列：haotrader使用优先级队列来管理订单，确保订单的执行顺序符合市场规则。这种数据结构可以高效地处理订单的添加、修改和取消请求。

  * 委托深度：haotrader支持实时生成和输出委托深度（order book），让交易用户能够查看市场上的订单情况和市场深度。

  * 最新成交价格：haotrader还提供了实时的最新成交价格，以供交易者参考，帮助他们做出更明智的交易决策。
  
  * 数据持久化：haotrader使用轻量级的boltdb/bolt存储，持久化接受到的订单数据，重启后快速从文件恢复数据。

  * 配置灵活：haotrader的配置非常灵活，用户可以根据自己的需求自定义市场规则、手续费、撮合逻辑等参数。

  * 开源：haotrader是开源项目，用户可以自由查看、使用和修改其源代码，以适应不同的交易场景和需求。

  * 跨平台支持：haotrader支持多个操作系统，包括 Linux、Windows 和 macOS，以确保在不同平台上运行稳定。


## 功能
* [x] 限价委托
* [x] 市价委托
  * [x] 市价按数量买入、卖出
  * [x] 市价按金额买入、卖出
* [x] 取消订单
* [x] 委托深度
* [x] 最新价格


## 接入方式
1. `go package`
  ```shell
    go get github.com/yzimhao/trading_engine
  ```
  具体详细使用方法参考 [Readme](https://github.com/yzimhao/trading_engine#引入包接入)


2. 独立程序
  * [使用文档](/trading_engine/haotrader)
  * [程序下载](https://github.com/yzimhao/trading_engine/releases/latest)


## Support or Contact
[需求建议](https://github.com/yzimhao/trading_engine#%E9%9C%80%E6%B1%82%E8%AE%A8%E8%AE%BA%E8%81%94%E7%B3%BB)
