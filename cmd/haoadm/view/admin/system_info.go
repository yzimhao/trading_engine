package admin

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/keepalive"
)

func SystemInfo(ctx *gin.Context) {

	data := make(map[string][]keepalive.App)

	rdc := app.RedisPool().Get()
	defer rdc.Close()

	data["haoadm"] = make([]keepalive.App, 0)
	data["haobase"] = make([]keepalive.App, 0)
	data["haomatch"] = make([]keepalive.App, 0)
	data["haoquote"] = make([]keepalive.App, 0)

	for _, topic := range keepalive.AppInfoTopic() {
		var ap keepalive.App
		info, _ := redis.Bytes(rdc.Do("get", topic))
		json.Unmarshal([]byte(info), &ap)

		data[ap.Name] = append(data[ap.Name], ap)
	}

	db := app.Database().NewSession()
	defer db.Close()

	utils.ResponseOkJson(ctx, gin.H{
		"module":       data,
		"version":      app.Version,
		"build_at":     app.Build,
		"assets_stats": assets.AssetsCheck(),
	})
}

func SystemTradeStats(ctx *gin.Context) {
	data := orders.NewTradeStats()
	render(ctx, 0, "", 0, data)
}
