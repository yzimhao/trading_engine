package view

import (
	"strings"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/view/admin"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func Run() {
	router := gin.Default()

	if config.App.Main.Mode != config.ModeProd {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	setupRouter(router)
	setupPages(router)

	router.Run(config.App.Haoadm.Listen)
}

type Router interface {
	Handle(httpMethod, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
}

func setMethods(r Router, methods []string, relativePath string, handlers ...gin.HandlerFunc) {
	for _, m := range methods {
		r.Handle(strings.ToUpper(m), relativePath, handlers...)
	}
}

func setupRouter(router *gin.Engine) {
	templateDir := "./template/default"
	router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		// Delims:    gintemplate.Delims{Left: "{{", Right: "}}"},
		Root:      templateDir, //template root path
		Extension: ".html",     //file extension
		// Master:    "",          //master layout file
		// Partials: []string{}, //partial files
		Funcs: templateFuncMap(),
		DisableCache: func() bool {
			return true
		}(),
	})
	router.Use(static.Serve("/", static.LocalFile(templateDir, false)))
	// router.Use(static.Serve("/uploads", static.LocalFile("./uploads", false)))
}

func setupPages(router *gin.Engine) {
	//admin
	radmin := router.Group("/admin")

	auth, _ := admin.AuthMiddleware()
	setMethods(radmin, []string{"GET"}, "/login", admin.Login)
	setMethods(radmin, []string{"POST"}, "/login", auth.LoginHandler)
	setMethods(radmin, []string{"GET"}, "/logout", auth.LogoutHandler)
	setMethods(radmin, []string{"GET"}, "/refresh_token", auth.RefreshHandler)

	// radmin.Use(auth.MiddlewareFunc())
	{
		setMethods(radmin, []string{"GET"}, "/index", admin.Index)
		setMethods(radmin, []string{"GET"}, "/welcome", admin.Welcome)
		setMethods(radmin, []string{"GET"}, "/system/settings", admin.Index)
		setMethods(radmin, []string{"GET"}, "/varieties/list", admin.VarietiesList)
		setMethods(radmin, []string{"GET"}, "/tradingvarieties/list", admin.TradingVarietiesList)
		setMethods(radmin, []string{"GET"}, "/user/assets", admin.AssetsList)
		setMethods(radmin, []string{"GET"}, "/user/order", admin.UserOrder)
		setMethods(radmin, []string{"GET"}, "/user/unfinished", admin.UserOrderUnfinished)
	}

	api := router.Group("/api/v1/admin")
	{
		setMethods(api, []string{"GET"}, "/system/menu", admin.SystemMenu)
	}
}
