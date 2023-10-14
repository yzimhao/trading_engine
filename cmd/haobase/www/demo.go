package www

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/base/symbols"
	"github.com/yzimhao/trading_engine/utils/app"
)

func demoBaseData() {
	if app.RunMode == app.ModeDemo {
		symbols.DemoData()
	}
}
