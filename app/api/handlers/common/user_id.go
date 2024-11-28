package common

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func GetUserId(ctx *gin.Context) string {
	claims := jwt.ExtractClaims(ctx)
	userId := claims["userId"].(string)
	return userId
}
