package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
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

func FormatStrNumber(n string, p int) string {
	d, _ := decimal.NewFromString(n)
	//TODO check err
	return d.StringFixed(int32(p))
}
