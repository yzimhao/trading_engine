package tradelog

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gookit/goutil/arrutil"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/haoquote/period"
	"github.com/yzimhao/trading_engine/haoquote/ws"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils"
	"xorm.io/xorm"
)

var (
	db  *xorm.Engine
	rdc *redis.Client
)

type TradeLog struct {
	Id            int64      `xorm:"autoincr pk" json:"-"`
	Symbol        string     `xorm:"-" json:"-"`
	TradeAt       int64      `xorm:"notnull" json:"trade_at"`
	TradePrice    string     `xorm:"decimal(30, 10) notnull" json:"price"`
	TradeQuantity string     `xorm:"decimal(30, 10) notnull" json:"qty"`
	TradeAmount   string     `xorm:"decimal(30, 10) notnull" json:"amount"`
	Ask           string     `xorm:"varchar(128) notnull unique(askbid)" json:"-"`
	Bid           string     `xorm:"varchar(128) notnull unique(askbid)" json:"-"`
	CreateTime    utils.Time `xorm:"timestamp created" json:"-"`
	UpdateTime    utils.Time `xorm:"timestamp updated" json:"-"`
}

func (t *TradeLog) TableName() string {
	return fmt.Sprintf("trade_log_%s", t.Symbol)
}

func (t *TradeLog) CreateTable() error {
	if t.Symbol == "" {
		return fmt.Errorf("symbol is null")
	}

	exist, err := db.IsTableExist(t.TableName())
	if err != nil {
		return err
	}

	if !exist {
		err := db.CreateTables(t)
		if err != nil {
			return err
		}
		err = db.CreateIndexes(t)
		if err != nil {
			return err
		}
		err = db.CreateUniques(t)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *TradeLog) Save() error {
	_, err := db.Table(t.TableName()).Insert(t)
	return err
}

func Init(rc *redis.Client, d *xorm.Engine) {
	db = d
	rdc = rc
}

func Monitor(symbol string, price_digit, qty_digit int64) {
	key := types.FormatQuoteTradeResult.Format(symbol)
	logrus.Infof("正在监听%s成交日志...", symbol)

	needPeriods := viper.GetStringSlice("haoquote.period")

	for {
		func() {
			cx := context.Background()
			if n, _ := rdc.LLen(cx, key).Result(); n == 0 {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := rdc.LPop(cx, key).Bytes()

			var data trading_core.TradeResult
			err := json.Unmarshal(raw, &data)
			if err != nil {
				logrus.Warnf("%s 解析json: %s 错误: %s", key, raw, err)
				return
			}

			//todo 保存成交日志到数据库
			row := TradeLog{
				Symbol:        symbol,
				TradeAt:       data.TradeTime,
				TradePrice:    data.TradePrice.String(),
				TradeQuantity: data.TradeQuantity.String(),
				TradeAmount:   data.TradePrice.Mul(data.TradeQuantity).String(),
				Ask:           data.AskOrderId,
				Bid:           data.BidOrderId,
			}
			row.CreateTable()
			if err := row.Save(); err != nil {
				logrus.Warnf("%s成交日志保存失败: %s %s %#v", symbol, raw, err, data)
				return
			}

			// todo 更多的period
			for _, curp := range period.Periods() {
				func(cp period.PeriodType) {
					if !arrutil.StringsHas(needPeriods, string(cp)) {
						return
					}

					row := period.NewPeriod(symbol, cp, data)
					err = save_db(row)
					if err != nil {
						logrus.Errorf("保存period数据出错: %s", err)
						return
					}

					//websocket通知更新
					to := types.MsgMarketKLine.Format(string(cp), symbol)
					ws.M.Broadcast <- ws.MsgBody{
						To: to,
						Response: ws.Response{
							Type: to,
							Body: [6]any{
								time.Time(row.OpenAt).UnixMilli(),
								utils.NumberFix(row.Open, int(price_digit)),
								utils.NumberFix(row.High, int(price_digit)),
								utils.NumberFix(row.Low, int(price_digit)),
								utils.NumberFix(row.Close, int(price_digit)),
								utils.NumberFix(row.Volume, int(qty_digit)),
							},
						},
					}
				}(curp)
			}

			//成交日志通知
			tradelog_msg(symbol, row, price_digit, qty_digit)
		}()

	}
}

func tradelog_msg(symbol string, data TradeLog, pd, qd int64) {
	data.TradePrice = utils.NumberFix(data.TradePrice, int(pd))
	data.TradeAmount = utils.NumberFix(data.TradeAmount, int(pd))
	data.TradeQuantity = utils.NumberFix(data.TradeQuantity, int(qd))

	to := types.MsgTrade.Format(symbol)
	ws.M.Broadcast <- ws.MsgBody{
		To: to,
		Response: ws.Response{
			Type: to,
			Body: data,
		},
	}
}

func save_db(row *period.Period) error {
	row.CreateTable(db)

	sess := db.NewSession()
	defer sess.Close()

	var err error
	exist, _ := sess.Table(row.TableName()).Where("open_at=?", row.OpenAt.Format()).Exist()
	if exist {
		_, err = sess.Table(row.TableName()).Where("open_at=?", row.OpenAt.Format()).ForUpdate().Update(row)
	} else {
		_, err = sess.Table(row.TableName()).Insert(row)
	}
	return err
}
