package utils

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

func D(a string) decimal.Decimal {
	d, _ := decimal.NewFromString(a)
	return d
}

func S2F64(a string) float64 {
	v, _ := strconv.ParseFloat(a, 64)
	return v
}
func S2Int(a string) int {
	v, _ := strconv.ParseInt(a, 10, 64)
	return int(v)
}

func S2Int64(a string) int64 {
	v, _ := strconv.ParseInt(a, 10, 64)
	return v
}

func Float2String(f float64, toFix int) string {
	format := "%." + fmt.Sprintf("%d", toFix) + "f"
	return fmt.Sprintf(format, f)
}

func NumberFix(a string, toFix int) string {
	b, _ := strconv.ParseFloat(a, 64)
	return Float2String(b, toFix)
}

func ResponseOkJson(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, gin.H{
		"ok":   1,
		"data": data,
	})
}

func ResponseFailJson(ctx *gin.Context, reason string) {
	ctx.JSON(http.StatusOK, gin.H{
		"ok":     0,
		"reason": reason,
	})
}

func Hash256(data any) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%v", data)))
	hashed := fmt.Sprintf("%x", hash.Sum(nil))
	return hashed
}
