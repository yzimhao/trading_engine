package entities

type AssetsLog struct {
	UUID
	Base
	UserId     string          `gorm:"type:varchar(30);not null;index" json:"user_id"`
	Symbol     string          `gorm:"type:varchar(30);not null;index" json:"symbol"`
	Before     string          `gorm:"type:decimal(40,20);default:0" json:"before"`                   // 变动前
	Amount     string          `gorm:"type:decimal(40,20);default:0" json:"amount"`                   // 变动数
	After      string          `gorm:"type:decimal(40,20);default:0" json:"after"`                    // 变动后
	TransID    string          `gorm:"type:varchar(100);not null;index:idx_trans_id" json:"trans_id"` // 业务相关的id
	ChangeType AssetChangeType `gorm:"type:varchar(15)" json:"change_type"`
	Info       string          `gorm:"type:varchar(200)" json:"info"`
}
