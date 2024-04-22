package message

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/cmd/haobase/message/ws"
	"github.com/yzimhao/trading_engine/types/redisdb"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Subscribe() {
	topic := redisdb.WsMessageQueue.Format(redisdb.Replace{})

	go func() {
		rdc := app.RedisPool().Get()
		defer rdc.Close()

		psc := redis.PubSubConn{Conn: rdc}
		psc.Subscribe(topic)
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				// app.Logger.Infof("收到消息[topic:%s]: %s", v.Channel, v.Data)
				var send_data ws.MsgBody
				err := json.Unmarshal(v.Data, &send_data)
				if err != nil {
					app.Logger.Errorf("解析消息出错 %s", v.Data)
					continue
				}
				ws.M.Broadcast <- send_data
			case redis.Subscription:
				// app.Logger.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				//todo: bug=> message queue.ws.message subscribe: read tcp 10.10.10.50:54276->10.10.10.15:6379: use of closed network connection
				app.Logger.Errorf("message %s subscribe: %s", topic, v.Error())
			}
		}
	}()
}

func Publish(msg ws.MsgBody) {
	rdc := app.RedisPool().Get()
	defer rdc.Close()
	topic := redisdb.WsMessageQueue.Format(redisdb.Replace{})

	raw := msg.JSON()
	// app.Logger.Infof("广播消息[topic:%s]: %s", topic, raw)
	if _, err := rdc.Do("Publish", topic, raw); err != nil {
		app.Logger.Warnf("Publish: %s %s err: %s", topic, raw, err.Error())
	}
}
