package www

import (
	"strings"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/www/admin"
	"github.com/yzimhao/trading_engine/utils/app"
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
	setupApi(router)
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
		Delims:    gintemplate.Delims{Left: "{%", Right: "%}"},
		Root:      templateDir, //template root path
		Extension: ".html",     //file extension
		// Master:    "layouts/",                //master layout file
		// Partials: []string{"partials/head"}, //partial files
		Funcs: templateFuncMap(),
		DisableCache: func() bool {
			return gin.Mode() == gin.DebugMode
		}(),
	})
	router.Use(static.Serve("/", static.LocalFile(templateDir, false)))
	// router.Use(static.Serve("/uploads", static.LocalFile("./uploads", false)))
}

func setupApi(router *gin.Engine) {
	apiv1 := router.Group("/api/")
	{
		app.Logger.Infof("%v", apiv1)
	}
}

func setupPages(router *gin.Engine) {
	//admin
	radmin := router.Group("/admin")

	auth, _ := admin.AuthMiddleware()
	setMethods(radmin, []string{"GET"}, "/login", admin.Login)
	setMethods(radmin, []string{"POST"}, "/login", auth.LoginHandler)
	setMethods(radmin, []string{"GET"}, "/logout", auth.LogoutHandler)
	setMethods(radmin, []string{"GET"}, "/refresh_token", auth.RefreshHandler)

	radmin.Use(auth.MiddlewareFunc())
	{
		// setMethods(radmin, []string{"GET"}, "/", admin.Index)

		// //添加、修改账号分组
		// setMethods(radmin, []string{"GET", "POST"}, "/account/group/add", admin.AccountGroupAdd)
		// setMethods(radmin, []string{"GET", "POST"}, "/account/group/account_setting", admin.AccountGroupSetting)
		// setMethods(radmin, []string{"GET", "POST"}, "/account/group/strategy_setting", admin.AccountGroupStrategySetting)
		// setMethods(radmin, []string{"GET", "POST"}, "/account/list", admin.AccountList)
		// setMethods(radmin, []string{"GET", "POST"}, "/account/position/list", admin.AccountPositionList)
		// setMethods(radmin, []string{"GET", "POST"}, "/account/position/history/list", admin.AccountPositionHistoryList)

		// //管理后台api
		// rapi := router.Group("/aapi")
		// {
		// 	setMethods(rapi, []string{"GET"}, "/system/init", admin.SystemInit)
		// 	setMethods(rapi, []string{"GET", "POST"}, "/account/group/list", admin.AccountGroupList)
		// 	setMethods(rapi, []string{"GET", "POST"}, "/account/group/delete", admin.AccountGroupDelete)
		// 	setMethods(rapi, []string{"GET", "POST"}, "/account/edit", admin.AccountEdit)
		// 	setMethods(rapi, []string{"GET", "POST"}, "/signals/list", admin.SignalsLogs)
		// }

	}
}
