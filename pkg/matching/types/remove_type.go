package types

type RemoveType int

const (
	RemoveTypeBySystem  RemoveType = iota + 1 //系统自动取消
	RemoveTypeByUser                          //用户主动取消
	RemoveTypeByPartial                       //部分成交
)
