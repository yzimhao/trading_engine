package view

import (
	"fmt"
	"html/template"

	"github.com/yzimhao/trading_engine/utils/app"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"app_power_by": func() string {
			return fmt.Sprintf("powered by Haotrader %s", app.Version)
		},
		"app_short_name": func() string {
			return "HT"
		},
		"version": func() string {
			return app.Version
		},
		"site_name": func() string {
			return config.App.Haoadm.SiteName
		},
	}
}
