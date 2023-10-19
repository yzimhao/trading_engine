package matching

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func http_start(addr string) {
	if viper.GetBool("haotrader.http.debug") {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	logrus.Infof("http服务监听: %s", addr)
	web_start(router)
	router.Run(addr)
}

func web_start(router *gin.Engine) {
	api := router.Group("api/v1")
	{
		// api.GET("/depth", symbol_depth)
		api.GET("/db_stats", db_stats)
	}

}

func symbol_depth(ctx *gin.Context) {
	_limit := ctx.Query("limit")
	symbol := strings.ToLower(ctx.Query("symbol"))

	if _, ok := teps[symbol]; !ok {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	limit, _ := strconv.Atoi(_limit)
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	ctx.JSON(http.StatusOK, gin.H{
		"asks": teps[symbol].GetAskDepth(limit),
		"bids": teps[symbol].GetBidDepth(limit),
	})
}

func db_stats(ctx *gin.Context) {
	queue := make(map[string]any)

	for symbol, _ := range teps {
		queue[symbol] = map[string]any{
			"AsksLength": teps[symbol].AskLen(),
			"BidsLength": teps[symbol].BidLen(),
		}
	}

	ctx.JSON(200, gin.H{
		"db":    localdb.Stats(),
		"queue": queue,
	})
}
