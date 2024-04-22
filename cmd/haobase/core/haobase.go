package core

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/base"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
	"github.com/yzimhao/trading_engine/cmd/haobase/orders"
	"github.com/yzimhao/trading_engine/cmd/haobase/www"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils/app"
)

func Run() {
	base.Init()
	initDemoBaseData()
	orders.Run()
	www.Run()
}

func initDemoBaseData() {
	if config.App.Main.Mode == app.ModeDemo.String() {
		varieties.DemoData()
	}
}
