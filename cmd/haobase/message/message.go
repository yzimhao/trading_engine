package message

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Subscribe() {
	topic := types.FormatWsMessage.Format("")

	go func() {
		rdc := app.RedisPool().Get()
		defer rdc.Close()

		psc := redis.PubSubConn{Conn: rdc}
		psc.Subscribe(topic)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				logrus.Infof("广播的消息 %s: message: %s\n", v.Channel, v.Data)
				var send_data ws.MsgBody
				err := json.Unmarshal(v.Data, &send_data)
				if err != nil {
					app.Logger.Errorf("解析消息出错 %s", v.Data)
					continue
				}
				ws.M.Broadcast <- send_data
			case redis.Subscription:
				logrus.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				logrus.Errorf("message %s subscribe: %s", topic, v.Error())
			}
		}
	}()
}

func Publish(msg ws.MsgBody) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()
	topic := types.FormatWsMessage.Format("")

	raw := msg.JSON()
	logrus.Infof("message %s publish: %s", topic, raw)
	if _, err := rdc.Do("Publish", topic, raw); err != nil {
		logrus.Warnf("广播%s消息: %s err: %s", topic, raw, err.Error())
	}
}
