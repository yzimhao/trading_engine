package base

import (
	"github.com/yzimhao/trading_engine/cmd/haobase/base/settings"
	"github.com/yzimhao/trading_engine/cmd/haobase/base/varieties"
)

func Init() {
	varieties.Init()
	settings.Init()
}
