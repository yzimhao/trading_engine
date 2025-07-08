package types

type OrderStatus int

const (
	OrderStatusNew           OrderStatus = iota + 1 //新建  但未被提交到市场
	OrderStatusPending                              //等待触发  等待触发
	OrderStatusSubmitted                            //已提交  已提交市场，等待执行
	OrderStatusPartialFill                          //部分成交  部分成交，但尚未完全执行
	OrderStatusFilled                               //已成交  订单已经完全执行，所有股票或合约已经被买入或卖出
	OrderStatusExpired                              //已过期  如果订单设置了有效期，且在规定时间内未能成交，订单可能会被标记为已过期
	OrderStatusRejected                             //已拒绝  交易所或经纪商可能会拒绝执行某些类型的订单，例如超过限制的市价单
	OrderStatusPartialCancel                        //部分取消   在部分成交后，交易者可能取消尚未执行的部分订单
	OrderStatusCanceled                             //已取消  交易者或系统取消了订单，订单不再有效
)
