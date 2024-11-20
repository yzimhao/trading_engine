package di

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/v2/app/template_func"
)

func NewGinEngine(v *viper.Viper) *gin.Engine {
	v.SetDefault("app.template_path", "./app/frontend/views/**/*.html")
	v.SetDefault("app.static_path", "./app/frontend/statics/")
	templatePath := v.GetString("app.template_path")
	staticPath := v.GetString("app.static_path")

	engine := gin.New()

	templateFunc := template_func.NewTemplateFunc()
	engine.SetFuncMap(templateFunc.FuncMap())
	engine.LoadHTMLGlob(templatePath)
	engine.StaticFS("/statics", http.Dir(staticPath))
	return engine
}
