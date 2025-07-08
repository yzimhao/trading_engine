package matching

import (
	"encoding/json"

	"github.com/shopspring/decimal"
	"github.com/yzimhao/trading_engine/v2/pkg/matching/types"
)

type order struct {
	OrderId      string             `json:"orderId"`
	OrderType    types.OrderType    `json:"orderType"`
	SubOrderType types.SubOrderType `json:"subOrderType"`
	Price        decimal.Decimal    `json:"price"`
	Quantity     decimal.Decimal    `json:"quantity"`
	Amount       decimal.Decimal    `json:"amount"`
	HoldQty      decimal.Decimal    `json:"HoldQty"`    //用户持有的数量，市价订单卖出的时候用到，用来退出撮合循环
	HoldAmount   decimal.Decimal    `json:"HoldAmount"` //用户持有的现金  市价订单买入的时候用到，用来退出撮合循环
	CreateTime   int64              `json:"createTime"`
	index        int                `json:"-"`
}

func (o *order) GetIndex() int {
	return o.index
}

func (o *order) SetIndex(index int) {
	o.index = index
}

func (o *order) SetQuantity(qnt decimal.Decimal) {
	o.Quantity = qnt
}

func (o *order) SetAmount(amount decimal.Decimal) {
	o.Amount = amount
}

func (o *order) SetHoldAmount(a decimal.Decimal) {
	o.HoldAmount = a
}

func (o *order) SetHoldQty(q decimal.Decimal) {
	o.HoldQty = q
}

func (o *order) GetUniqueId() string {
	return o.OrderId
}

func (o *order) GetPrice() decimal.Decimal {
	return o.Price
}

func (o *order) GetQuantity() decimal.Decimal {
	return o.Quantity
}

func (o *order) GetAmount() decimal.Decimal {
	return o.Amount
}

func (o *order) GetCreateTime() int64 {
	return o.CreateTime
}

func (o *order) GetOrderType() types.OrderType {
	return o.OrderType
}

func (o *order) GetSubOrderType() types.SubOrderType {
	return o.SubOrderType
}

func (o *order) GetHoldAmount() decimal.Decimal {
	return o.HoldAmount
}
func (o *order) GetHoldQty() decimal.Decimal {
	return o.HoldQty
}

func (o *order) Marshal() []byte {
	raw, _ := json.Marshal(o)
	return raw
}

// 这个方法留在具体的 ask/bid 队列中实现
// func (o *Order) Less() {}

type AskItem struct {
	order
}

func (a *AskItem) Less(other QueueItem) bool {
	//价格优先，时间优先原则
	//价格低的在最上面
	return (a.Price.Cmp(other.(*AskItem).Price) == -1) || (a.Price.Cmp(other.(*AskItem).Price) == 0 && a.CreateTime < other.(*AskItem).CreateTime)
}

func (a *AskItem) GetOrderSide() types.OrderSide {
	return types.OrderSideSell
}

type BidItem struct {
	order
}

func (a *BidItem) Less(other QueueItem) bool {
	//价格优先，时间优先原则
	//价格高的在最上面
	return (a.Price.Cmp(other.(*BidItem).Price) == 1) || (a.Price.Cmp(other.(*BidItem).Price) == 0 && a.CreateTime < other.(*BidItem).CreateTime)
}

func (a *BidItem) GetOrderSide() types.OrderSide {
	return types.OrderSideBuy
}

func NewAskLimitItem(uniq string, price, quantity decimal.Decimal, createTime int64) *AskItem {
	return newAskItem(types.OrderTypeLimit, uniq, price, quantity, decimal.Zero, decimal.Zero, decimal.Zero, createTime, types.SubOrderTypeUnknown)
}

// 市价， 按成交个数卖出
func NewAskMarketQtyItem(uniq string, quantity decimal.Decimal, createTime int64) *AskItem {
	return newAskItem(types.OrderTypeMarket, uniq, decimal.Zero, quantity, decimal.Zero, decimal.Zero, decimal.Zero, createTime, types.SubOrderTypeMarketByQty)
}

// 市价 按成交金额卖出订单，需要用户持有交易物的数量，在撮合时达到成交金额后退出循环
func NewAskMarketAmountItem(uniq string, amount, maxHoldQty decimal.Decimal, createTime int64) *AskItem {
	return newAskItem(types.OrderTypeMarket, uniq, decimal.Zero, decimal.Zero, amount, maxHoldQty, decimal.Zero, createTime, types.SubOrderTypeMarketByAmount)
}

func NewBidLimitItem(uniq string, price, quantity decimal.Decimal, createTime int64) *BidItem {
	return newBidItem(types.OrderTypeLimit, uniq, price, quantity, decimal.Zero, decimal.Zero, decimal.Zero, createTime, types.SubOrderTypeUnknown)
}

// 市价，按成交个数买入
func NewBidMarketQtyItem(uniq string, quantity, holdAmount decimal.Decimal, createTime int64) *BidItem {
	return newBidItem(types.OrderTypeMarket, uniq, decimal.Zero, quantity, decimal.Zero, decimal.Zero, holdAmount, createTime, types.SubOrderTypeMarketByQty)
}

// 市价，按成交金额买入
func NewBidMarketAmountItem(uniq string, amount decimal.Decimal, createTime int64) *BidItem {
	return newBidItem(types.OrderTypeMarket, uniq, decimal.Zero, decimal.Zero, amount, decimal.Zero, amount, createTime, types.SubOrderTypeMarketByAmount)
}

func newAskItem(
	pt types.OrderType,
	uniqId string, price, quantity, amount decimal.Decimal,
	holdQty, holdAmount decimal.Decimal,
	createTime int64, subT types.SubOrderType) *AskItem {
	return &AskItem{
		order: order{
			OrderId:      uniqId,
			Price:        price,
			Quantity:     quantity,
			CreateTime:   createTime,
			OrderType:    pt,
			SubOrderType: subT,
			Amount:       amount,
			HoldAmount:   holdAmount,
			HoldQty:      holdQty,
		},
	}
}
func newBidItem(
	pt types.OrderType,
	uniqId string, price, quantity, amount decimal.Decimal,
	holdQty, holdAmount decimal.Decimal,
	createTime int64, subT types.SubOrderType) *BidItem {
	return &BidItem{
		order: order{
			OrderId:      uniqId,
			Price:        price,
			Quantity:     quantity,
			CreateTime:   createTime,
			OrderType:    pt,
			SubOrderType: subT,
			Amount:       amount,
			HoldAmount:   holdAmount,
			HoldQty:      holdQty,
		}}
}
