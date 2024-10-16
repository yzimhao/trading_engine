package entities

type AssetsFreeze struct {
	UUID
	Base
	UserId       string       `gorm:"type:varchar(30);index;not null" json:"user_id"`
	Symbol       string       `gorm:"type:varchar(30);index;not null" json:"symbol"`
	Amount       string       `gorm:"type:decimal(40,20);default:0;not null" json:"amount"`            // 冻结总量
	FreezeAmount string       `gorm:"type:decimal(40,20);default:0;not null" json:"freeze_amount"`     // 冻结中的量
	Status       FreezeStatus `gorm:"type:tinyint(1)" json:"status"`                                   // 状态 冻结中, 已解冻
	TransId      string       `gorm:"type:varchar(100);uniqueIndex:trans_id;not null" json:"trans_id"` // 业务相关的id
	FreezeType   FreezeType   `gorm:"type:varchar(15)" json:"freeze_type"`
	Info         string       `gorm:"type:varchar(200)" json:"info"`
}
