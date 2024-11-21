package di

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/yzimhao/trading_engine/v2/app/template_func"
	"go.uber.org/zap"
)

func NewGinEngine(v *viper.Viper, logger *zap.Logger) *gin.Engine {
	v.SetDefault("app.template_path", "app/frontend/views/")
	templatePath := v.GetString("app.template_path")
	v.SetDefault("app.static_path", "./app/frontend/statics/")

	staticPath := v.GetString("app.static_path")

	engine := gin.New()

	templateFunc := template_func.NewTemplateFunc()
	engine.SetFuncMap(templateFunc.FuncMap())

	engine.HTMLRender = renderer(templatePath, logger)

	engine.StaticFS("/statics", http.Dir(staticPath))
	return engine
}

func renderer(templatePath string, logger *zap.Logger) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	// 确保路径为规范化的绝对路径
	absTemplatePath, err := filepath.Abs(templatePath)
	if err != nil {
		logger.Sugar().Panic("Invalid template path", zap.Error(err))
	}

	logger.Sugar().Debug("Resolved template path", zap.String("path", absTemplatePath))

	// 加载所有 HTML 模板
	tpls, err := filepath.Glob(filepath.Join(absTemplatePath, "**/*.html"))
	if err != nil {
		logger.Sugar().Panic("Error loading templates", zap.Error(err))
	}

	if len(tpls) == 0 {
		logger.Sugar().Panic("No templates found in path", zap.String("path", absTemplatePath))
	}

	// 为每个模板文件设置唯一名称
	for _, tpl := range tpls {
		name := strings.TrimPrefix(tpl, absTemplatePath)
		name = strings.TrimPrefix(name, string(filepath.Separator)) // 去除首部的路径分隔符
		logger.Sugar().Debug("Registering template", zap.String("name", name), zap.String("path", tpl))
		r.AddFromFiles(name, tpl)
	}
	return r
}
