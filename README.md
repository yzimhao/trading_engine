## HaoTrader
  
<p align="center">
    <img src="https://img.shields.io/github/stars/yzimhao/trading_engine?style=social">
    <img src="https://img.shields.io/github/forks/yzimhao/trading_engine?style=social">
	<img src="https://img.shields.io/github/issues/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/repo-size/yzimhao/trading_engine">
	<img src="https://img.shields.io/github/license/yzimhao/trading_engine">
</p>

  HaoTrader适用于各种金融证券交易场景。拥有高性能的订单撮合、实时结算、行情计算、实时推送等功能。
  配置灵活，允许用户根据自身需求自定义配置各模块独立运行。

  该程序目前只包含服务端API，无UI。

  __征集web端UI__

##
  ![image](https://github.com/yzimhao/trading_engine/blob/master/document/images/haotrader.png?raw=true)

## 演示DEMO
  <a href="http://144.91.108.90:20001/" target="_blank">在线体验</a> 
  > 感谢[9cat](https://github.com/9cat)大佬提供免费测试服务器 

  <a href="http://144.91.108.90:20010/admin/index" target="_blank">运营后台</a> 



## HaoTrader系统包含模块
  - [ ] [haoadm]
    - [ ] 后台登陆认证+权限控制
    - [ ] 系统设置
    - [ ] 交易报表
  - [x] [haobase]
    - [x] 交易品种
    - [x] 资产模块
    - [x] 订单模块
    - [x] 结算模块
    - [x] 充值提现

  - [x] [haomatch]
    - [x] 撮合模块

  - [x] [haoquote]
    - [x] 深度行情
    - [x] Kline数据
    - [x] websocket推送

## 撮合
```
    go get github.com/yzimhao/trading_engine/trading_core
```
- <a href="/document/match.md">文档</a>


## 开发文档
- <a href="https://yzimhao.github.io/trading_engine/">交易系统文档</a>



## 参考
- <a href="https://www.liaoxuefeng.com/article/1185272483766752" target="_blank">证券交易系统设计与开发</a>

## 需求讨论
   <img src="https://github.com/yzimhao/trading_engine/blob/master/document/images/wechat.jpg?raw=true" width = "150"/>

## 声明
- 本项目仅供参考和学习之用，不建议将其用于生产环境或重要交易场景。
