package www

import (
	"html/template"

	"github.com/yzimhao/trading_engine/utils/app/config"
)

func templateFuncMap() template.FuncMap {
	return template.FuncMap{
		"app_name": func() string {
			return "HaoTrader"
		},

		"site_name": func() string {
			return config.App.Haoadm.SiteName
		},
	}
}
