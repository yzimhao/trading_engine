package message

import (
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/types"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Subscribe() {
	rdc := app.RedisPool().Get()
	defer rdc.Close()
	topic := types.FormatWsMessage.Format("")

	go func() {
		psc := redis.PubSubConn{Conn: rdc}
		psc.Subscribe(topic)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				logrus.Infof("%s: message: %s\n", v.Channel, v.Data)
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
