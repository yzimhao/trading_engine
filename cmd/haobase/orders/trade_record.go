package orders

import (
	"fmt"

	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

type TradeBy int

const (
	TradeBySell TradeBy = 1
	TradeByBuy  TradeBy = 2
)

// 成交记录表
type TradeLog struct {
	Id         int64      `xorm:"pk autoincr bigint" json:"id"`
	Symbol     string     `xorm:"-" json:"-"`
	TradeId    string     `xorm:"varchar(30) unique(trade)" json:"trade_id"`
	Ask        string     `xorm:"varchar(30) unique(trade)" json:"ask"`
	Bid        string     `xorm:"varchar(30) unique(trade)" json:"bid"`
	TradeBy    TradeBy    `xorm:"tinyint(1)" json:"trade_by"`
	AskUid     string     `xorm:"notnull" json:"ask_uid"`
	BidUid     string     `xorm:"notnull" json:"bid_uid"`
	Price      string     `xorm:"decimal(40,20) notnull default(0)" json:"price"`
	Quantity   string     `xorm:"decimal(40,20) notnull default(0)" json:"quantity"`
	Amount     string     `xorm:"decimal(40,20) notnull default(0)" json:"amount"`
	AskFeeRate string     `xorm:"decimal(40,20) notnull default(0)" json:"ask_fee_rate"`
	AskFee     string     `xorm:"decimal(40,20) notnull default(0)" json:"ask_fee"`
	BidFeeRate string     `xorm:"decimal(40,20) notnull default(0)" json:"bid_fee_rate"`
	BidFee     string     `xorm:"decimal(40,20) notnull default(0)" json:"bid_fee"`
	CreateTime utils.Time `xorm:"timestamp created" json:"create_time"`
	UpdateTime utils.Time `xorm:"timestamp updated" json:"update_time"`
}

func (tr *TradeLog) Save(db *xorm.Session) error {
	_, err := db.Table(tr).Insert(tr)
	if err != nil {
		return err
	}
	return nil
}

func (tr *TradeLog) FormatDecimal(price_digit, qty_digit int) TradeLog {
	tr.Amount = utils.FormatDecimal(tr.Amount, price_digit)
	tr.Price = utils.FormatDecimal(tr.Price, price_digit)
	tr.Quantity = utils.FormatDecimal(tr.Quantity, qty_digit)

	tr.AskFeeRate = utils.FormatDecimal(tr.AskFeeRate, price_digit)
	tr.AskFee = utils.FormatDecimal(tr.AskFee, price_digit)
	tr.BidFeeRate = utils.FormatDecimal(tr.BidFeeRate, price_digit)
	tr.BidFee = utils.FormatDecimal(tr.BidFee, price_digit)
	return *tr
}

func (tr *TradeLog) TableName() string {
	return fmt.Sprintf("%strade_log_%s", app.TablePrefix(), tr.Symbol)
}
