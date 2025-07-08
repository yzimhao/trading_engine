package types

type RemoveItemType int

const (
	RemoveItemTypeByUser    RemoveItemType = iota + 1 //用户主动取消
	RemoveItemTypeBySystem                            //系统取消
	RemoveItemTypeByExpired                           //超时取消
	RemoveItemTypeByMarket                            //市场取消
	RemoveItemTypeByForce                             //强平取消
	RemoveItemTypeByOther                             //其他
)
