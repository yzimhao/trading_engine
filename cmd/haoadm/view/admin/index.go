package admin

import (
	"github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context) {
	ctx.HTML(200, "layout_index", gin.H{})
}
