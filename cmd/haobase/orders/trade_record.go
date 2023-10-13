package orders

import (
	"fmt"
	"time"

	"github.com/yzimhao/trading_engine/utils"
	"xorm.io/xorm"
)

type TradeBy int

const (
	TradeBySell TradeBy = 1
	TradeByBuy  TradeBy = 2
)

// 成交记录表
type TradeLog struct {
	Id         int64      `xorm:"pk autoincr bigint" json:"-"`
	Symbol     string     `xorm:"-" json:"-"`
	TradeId    string     `xorm:"varchar(30) unique(trade)" json:"trade_id"`
	Ask        string     `xorm:"varchar(30) unique(trade)" json:"ask"`
	Bid        string     `xorm:"varchar(30) unique(trade)" json:"bid"`
	TradeBy    TradeBy    `xorm:"tinyint(1)" json:"trade_by"`
	AskUid     string     `xorm:"notnull" json:"-"`
	BidUid     string     `xorm:"notnull" json:"-"`
	Price      string     `xorm:"decimal(40,20) notnull default(0)" json:"price"`
	Quantity   string     `xorm:"decimal(40,20) notnull default(0)" json:"quantity"`
	Amount     string     `xorm:"decimal(40,20) notnull default(0)" json:"amount"`
	AskFeeRate string     `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	AskFee     string     `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	BidFeeRate string     `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	BidFee     string     `xorm:"decimal(40,20) notnull default(0)" json:"-"`
	CreateTime utils.Time `xorm:"timestamp created" json:"trade_at"`
	UpdateTime time.Time  `xorm:"timestamp updated" json:"-"`
}

func (tr *TradeLog) Save(db *xorm.Session) error {
	if tr.Symbol == "" {
		return fmt.Errorf("symbol not set")
	}
	//todo 频繁查询表是否存在，后面考虑缓存一下
	exist, err := db.IsTableExist(tr.TableName())
	if err != nil {
		return err
	}
	if !exist {
		err := db.CreateTable(tr)
		if err != nil {
			return err
		}

		err = db.CreateIndexes(tr)
		if err != nil {
			return err
		}

		err = db.CreateUniques(tr)
		if err != nil {
			return err
		}
	}

	_, err = db.Table(tr).Insert(tr)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TradeLog) TableName() string {
	return fmt.Sprintf("trade_log_%s", tr.Symbol)
}

func GetTradelogTableName(symbol string) string {
	t := TradeLog{Symbol: symbol}
	return t.TableName()
}
