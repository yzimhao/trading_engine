package admin

import (
	"github.com/gin-gonic/gin"
)

// layui response
func render(ctx *gin.Context, code int, msg string, total int, data interface{}) {
	ctx.JSON(200, gin.H{
		"code":  code,
		"msg":   msg,
		"count": total,
		"data":  data,
	})
}

func Index(ctx *gin.Context) {
	ctx.HTML(200, "layout_index", gin.H{})
}
