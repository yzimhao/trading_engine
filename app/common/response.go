package common

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ResponseOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"data": data,
	})
}

func ResponseError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"ok":  false,
		"msg": err.Error(),
	})
}

func NumberFix(n string, toFix int) string {
	b, _ := strconv.ParseFloat(n, 64)
	format := "%." + fmt.Sprintf("%d", toFix) + "f"
	return fmt.Sprintf(format, b)
}
