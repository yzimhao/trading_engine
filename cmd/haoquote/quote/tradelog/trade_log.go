package tradelog

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gookit/goutil/arrutil"
	"github.com/yzimhao/trading_engine/cmd/haobase/message"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/cmd/haoquote/quote/period"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/trading_core"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/types/dbtables"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

var ()

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
	//和haobase中订单结算后的成交日志有一点点冲突，这里只保存最近的记录就好了
	return fmt.Sprintf("%squote_tradelog_%s", app.TablePrefix(), t.Symbol)
}

func (t *TradeLog) Save() error {
	db := app.Database().NewSession()
	defer db.Close()

	if err := dbtables.AutoCreateTable(db, t); err != nil {
		return err
	}

	_, err := db.Table(t.TableName()).Insert(t)
	return err
}

func Monitor(symbol string, price_digit, qty_digit int64) {
	key := types.FormatQuoteTradeResult.Format(symbol)
	app.Logger.Infof("监听 %s 成交日志...", symbol)

	for {
		func() {
			rdc := app.RedisPool().Get()
			defer rdc.Close()
			if n, _ := redis.Int64(rdc.Do("LLen", key)); n == 0 {
				time.Sleep(time.Duration(50) * time.Millisecond)
				return
			}

			raw, _ := redis.Bytes(rdc.Do("LPop", key))

			var data trading_core.TradeResult
			err := json.Unmarshal(raw, &data)
			if err != nil {
				app.Logger.Errorf("%s 解析: %s 错误: %s", key, raw, err)
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

			row.Save()

			// todo 更多的period
			for _, curp := range period.Periods() {
				func(cp period.PeriodType) {
					if !arrutil.StringsHas(config.App.Haoquote.Period, string(cp)) {
						return
					}

					row := period.NewPeriod(symbol, cp, data)
					err = save_db(row)
					if err != nil {
						app.Logger.Errorf("保存period数据出错: %s", err)
						return
					}

					//websocket通知更新
					to := types.MsgMarketKLine.Format(map[string]string{
						"symbol": symbol,
						"period": string(cp),
					})

					message.Publish(ws.MsgBody{
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
					})

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

	to := types.MsgTrade.Format(map[string]string{
		"symbol": symbol,
	})
	message.Publish(ws.MsgBody{
		To: to,
		Response: ws.Response{
			Type: to,
			Body: data,
		},
	})
}

func save_db(row *period.Period) error {

	db := app.Database().NewSession()
	defer db.Close()

	if err := row.CreateTable(db); err != nil {
		return err
	}

	var err error
	exist, _ := db.Table(row.TableName()).Where("open_at=?", row.OpenAt.Format()).Exist()
	if exist {
		_, err = db.Table(row.TableName()).Where("open_at=?", row.OpenAt.Format()).ForUpdate().Update(row)
	} else {
		_, err = db.Table(row.TableName()).Insert(row)
	}
	return err
}
