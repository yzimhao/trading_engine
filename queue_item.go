package trading_engine

import (
	"github.com/shopspring/decimal"
)

type Order struct {
	orderId    string
	price      decimal.Decimal
	quantity   decimal.Decimal
	createTime int64
	index      int

	priceType PriceType
	amount    decimal.Decimal
}

func (o *Order) GetIndex() int {
	return o.index
}

func (o *Order) SetIndex(index int) {
	o.index = index
}

func (o *Order) SetQuantity(qnt decimal.Decimal) {
	o.quantity = qnt
}

func (o *Order) GetUniqueId() string {
	return o.orderId
}

func (o *Order) GetPrice() decimal.Decimal {
	return o.price
}

func (o *Order) GetQuantity() decimal.Decimal {
	return o.quantity
}

func (o *Order) GetCreateTime() int64 {
	return o.createTime
}

func (o *Order) GetPriceType() PriceType {
	return o.priceType
}
func (o *Order) GetAmount() decimal.Decimal {
	return o.amount
}

// 这个方法留在具体的 ask/bid 队列中实现
// func (o *Order) Less() {}

type AskItem struct {
	Order
}

func (a *AskItem) Less(other QueueItem) bool {
	//价格优先，时间优先原则
	//价格低的在最上面
	return (a.price.Cmp(other.(*AskItem).price) == -1) || (a.price.Cmp(other.(*AskItem).price) == 0 && a.createTime < other.(*AskItem).createTime)
}

func (a *AskItem) GetOrderSide() OrderSide {
	return OrderSideSell
}

type BidItem struct {
	Order
}

func (a *BidItem) Less(other QueueItem) bool {
	//价格优先，时间优先原则
	//价格高的在最上面
	return (a.price.Cmp(other.(*BidItem).price) == 1) || (a.price.Cmp(other.(*BidItem).price) == 0 && a.createTime < other.(*BidItem).createTime)
}

func (a *BidItem) GetOrderSide() OrderSide {
	return OrderSideBuy
}

func NewAskItem(pt PriceType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *AskItem {
	return &AskItem{
		Order: Order{
			orderId:    uniqId,
			price:      price,
			quantity:   quantity,
			createTime: createTime,
			priceType:  pt,
			amount:     amount,
		},
	}
}

func NewBidItem(pt PriceType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *BidItem {
	return &BidItem{
		Order: Order{
			orderId:    uniqId,
			price:      price,
			quantity:   quantity,
			createTime: createTime,
			priceType:  pt,
			amount:     amount,
		}}
}
