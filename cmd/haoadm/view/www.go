package view

import (
	"strings"

	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/models"
	"github.com/yzimhao/trading_engine/cmd/haoadm/view/admin"
	"github.com/yzimhao/trading_engine/utils"
	"github.com/yzimhao/trading_engine/utils/app/config"
)

func Run() {
	router := gin.Default()

	if config.App.Main.Mode != config.ModeProd {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	models.Init()

	// store := cookie.NewStore([]byte(config.App.Main.SecretKey))
	// router.Use(sessions.Sessions("tk", store))

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
		Delims:    gintemplate.Delims{Left: "{%", Right: "%}"},
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
}

func runModeCheck(ctx *gin.Context) {
	if config.App.Main.Mode == config.ModeDemo && config.App.Haoadm.Readonly {
		if ctx.Request.Method == "POST" && ctx.Request.RequestURI != "/admin/login" {
			ctx.Abort()
			utils.ResponseFailJson(ctx, "demo模式，禁止修改数据")
			return
		}
	}
}

func setupPages(router *gin.Engine) {
	//后台页面
	radmin := router.Group("/admin")
	//后台接口
	api := router.Group("/api/v1/admin")
	radmin.Use(runModeCheck)
	api.Use(runModeCheck)

	auth, _ := admin.AuthMiddleware()
	setMethods(radmin, []string{"GET"}, "/login", admin.Login)
	setMethods(radmin, []string{"POST"}, "/login", auth.LoginHandler)
	setMethods(radmin, []string{"GET"}, "/logout", auth.LogoutHandler)
	setMethods(radmin, []string{"GET"}, "/refresh_token", auth.RefreshHandler)

	radmin.Use(auth.MiddlewareFunc())
	{
		setMethods(radmin, []string{"GET"}, "/index", admin.Index)
		setMethods(radmin, []string{"GET"}, "/welcome", admin.Welcome)
		setMethods(radmin, []string{"GET"}, "/system/settings", admin.Index)
		setMethods(radmin, []string{"GET"}, "/varieties/list", admin.VarietiesList)
		setMethods(radmin, []string{"GET", "POST"}, "/varieties/add", admin.VarietiesAdd)
		setMethods(radmin, []string{"GET"}, "/tradingvarieties/list", admin.TradingVarietiesList)
		setMethods(radmin, []string{"GET", "POST"}, "/tradingvarieties/add", admin.TradingVarietiesAdd)

		setMethods(radmin, []string{"GET"}, "/user/assets", admin.AssetsList)
		setMethods(radmin, []string{"GET"}, "/user/assets/freeze", admin.AssetsFreezeList)
		setMethods(radmin, []string{"GET"}, "/user/assets/logs", admin.AssetsLogsList)

		setMethods(radmin, []string{"GET"}, "/user/order/history", admin.UserOrderHistory)
		setMethods(radmin, []string{"GET"}, "/user/trade/history", admin.TradeHistory)
		setMethods(radmin, []string{"GET"}, "/user/unfinished", admin.UserOrderUnfinished)
	}
	api.Use(auth.MiddlewareFunc())
	{
		setMethods(api, []string{"GET"}, "/system/menu", admin.SystemMenu)
		setMethods(api, []string{"GET"}, "/system/info", admin.SystemInfo)
		setMethods(api, []string{"GET"}, "/system/tradestats", admin.SystemTradeStats)
		setMethods(api, []string{"POST"}, "/user/unfinished/cancel", admin.CancelUserOrder)
	}
}
