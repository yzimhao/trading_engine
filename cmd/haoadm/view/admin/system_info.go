package admin

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haobase/assets"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app"
)

func SystemInfo(ctx *gin.Context) {
	info := app.KeepaliveInfo()
	data := make([]string, 0)

	for _, item := range info {
		name := strings.Split(item, ".")[1]
		data = append(data, name)
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
