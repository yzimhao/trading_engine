package period

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"xorm.io/xorm"
)

type Period struct {
	Symbol        string                   `xorm:"-" json:"-"`
	OpenAt        utils.Time               `xorm:"notnull timestamp unique(open_at) default CURRENT_TIMESTAMP" json:"open_at"` //开盘时间
	CloseAt       utils.Time               `xorm:"notnull timestamp default CURRENT_TIMESTAMP" json:"close_at"`                // 收盘时间
	Open          string                   `xorm:"decimal(30, 10) notnull" json:"open"`                                        //开盘价
	High          string                   `xorm:"decimal(30, 10) notnull" json:"high"`                                        // 最高价
	Low           string                   `xorm:"decimal(30, 10) notnull" json:"low"`                                         //最低价
	Close         string                   `xorm:"decimal(30, 10) notnull" json:"close"`                                       //收盘价(当前K线未结束的即为最新价)
	Volume        string                   `xorm:"decimal(30, 10) notnull" json:"volume"`                                      //成交量
	Amount        string                   `xorm:"decimal(30, 10) notnull" json:"amount"`                                      //成交额
	CreateTime    utils.Time               `xorm:"timestamp created" json:"-"`
	UpdateTime    utils.Time               `xorm:"timestamp updated" json:"-"`
	Interval      PeriodType               `xorm:"-" json:"-"`
	raw           trading_core.TradeResult `xorm:"-" json:"-"`
	LastOpenTime  int64                    `xorm:"-" json:"last_open_time"`
	LastCloseTime int64                    `xorm:"-" json:"last_close_time"`
}

func NewPeriod(symbol string, p PeriodType, tr trading_core.TradeResult) *Period {
	tradetime := time.Unix(int64(tr.TradeTime/1e9), 0)
	open_at, close_at := parse_start_end_time(tradetime, p)

	data := Period{}

	ckey := formatKLineKey(p, symbol, open_at.Unix(), close_at.Unix())
	cache_data, err := getKLinePeriod(ckey)
	if err != nil {
		app.Logger.Errorf("get %s cache error: %s", ckey, err)
	}
	json.Unmarshal(cache_data, &data)

	app.Logger.Infof("get %s cache: [%#v]", ckey, data)
	defer func() {
		raw, _ := json.Marshal(data)
		app.Logger.Infof("set %s cache: [%#v cur_price:%s]", ckey, data, tr.TradePrice.String())

		ttl := close_at.Unix() - time.Now().Unix() + 60*60
		setKLinePeriod(ckey, raw, ttl)
	}()

	data.raw = tr
	data.Interval = p
	data.Symbol = symbol
	data.OpenAt = utils.Time(open_at)
	data.CloseAt = utils.Time(close_at)

	data.get_open()
	data.get_high()
	data.get_low()
	data.get_close()
	data.get_volume()
	data.get_amount()

	return &data
}

func (p *Period) TableName() string {
	return fmt.Sprintf("%squote_period_%s_%s", app.TablePrefix(), p.Symbol, p.Interval)
}

func (p *Period) CreateTable(db *xorm.Session) error {
	if p.Symbol == "" || p.Interval == "" {
		return fmt.Errorf("symbol or period is null")
	}

	return dbtables.AutoCreateTable(db, p)
}

func (p *Period) get_open() {
	if p.Open == "" {
		p.Open = p.raw.TradePrice.String()
		p.High = p.Open
		p.Low = p.Open
		p.Close = p.Open

		p.Volume = "0"
		p.Amount = "0"

		p.LastOpenTime = p.raw.TradeTime
	}
	if p.raw.TradeTime < p.LastOpenTime {
		p.Open = p.raw.TradePrice.String()
		p.LastOpenTime = p.raw.TradeTime
	}
}

func (p *Period) get_high() {
	if p.raw.TradePrice.Cmp(utils.D(p.High)) > 0 {
		p.High = p.raw.TradePrice.String()
	}
}

func (p *Period) get_low() {
	if p.raw.TradePrice.Cmp(utils.D(p.Low)) < 0 {
		p.Low = p.raw.TradePrice.String()
	}
}

func (p *Period) get_close() {
	if p.raw.TradeTime > p.LastCloseTime {
		p.Close = p.raw.TradePrice.String()
		p.LastCloseTime = p.raw.TradeTime
	}
}

func (p *Period) get_volume() {
	v := utils.D(p.Volume).Add(p.raw.TradeQuantity)
	p.Volume = v.String()
}

func (p *Period) get_amount() {
	v := utils.D(p.Amount).Add(p.raw.TradePrice.Mul(p.raw.TradeQuantity))
	p.Amount = v.String()
}
