package admin

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/yzimhao/trading_engine/cmd/haoadm/models"
	"github.com/yzimhao/trading_engine/config"
	"github.com/yzimhao/trading_engine/utils/app"
)

type login struct {
	Name     string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var (
	auth *jwt.GinJWTMiddleware
)

// the jwt middleware
func AuthMiddleware() *jwt.GinJWTMiddleware {
	var identityKey = "user_id"
	if auth == nil {
		auth, _ = jwt.New(&jwt.GinJWTMiddleware{
			Realm:       "HaoTrader",
			Key:         []byte(config.App.Main.SecretKey),
			Timeout:     time.Hour * time.Duration(24),
			MaxRefresh:  time.Hour,
			IdentityKey: identityKey,
			PayloadFunc: func(data interface{}) jwt.MapClaims {
				if v, ok := data.(*models.Adminuser); ok {
					return jwt.MapClaims{
						"user_id":  v.Id,
						"username": v.Username,
					}
				}
				return jwt.MapClaims{}
			},
			IdentityHandler: func(c *gin.Context) interface{} {
				claims := jwt.ExtractClaims(c)
				return &models.Adminuser{
					Id: func() int64 {
						a := claims["user_id"].(float64)
						return int64(a)
					}(),
				}
			},
			Authenticator: func(c *gin.Context) (interface{}, error) {
				var loginVals login
				if err := c.ShouldBind(&loginVals); err != nil {
					return "", jwt.ErrMissingLoginValues
				}
				username := loginVals.Name
				password := loginVals.Password

				db := app.Database().NewSession()
				defer db.Close()

				var user models.Adminuser
				exist, err := db.Table(models.Adminuser{}).Where("username=?", username).Get(&user)
				if err != nil {
					return nil, err
				}

				if !exist {
					return nil, jwt.ErrMissingLoginValues
				}

				if err := user.ComparePassword(password); err == nil {
					db.Table(models.Adminuser{}).Where("id=?", user.Id).Update(&models.Adminuser{
						// LoginIp: c.ClientIP(),
					})

					return &user, nil
				}
				//todo 记录登陆错误的数据

				return nil, jwt.ErrFailedAuthentication
			},
			Authorizator: func(data interface{}, c *gin.Context) bool {
				// if v, ok := data.(*User); ok && v.UserName == "admin" {
				// 	return true
				// }

				// return false
				return true
			},
			Unauthorized: func(c *gin.Context, code int, message string) {
				c.Redirect(301, "/admin/login")

				c.JSON(code, gin.H{
					"code":    code,
					"message": message,
				})
			},
			LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
				c.Redirect(301, "/admin/index")
			},
			LogoutResponse: func(c *gin.Context, code int) {
				c.Redirect(301, "/admin/login")
			},
			// TokenLookup is a string in the form of "<source>:<name>" that is used
			// to extract token from the request.
			// Optional. Default value "header:Authorization".
			// Possible values:
			// - "header:<name>"
			// - "query:<name>"
			// - "cookie:<name>"
			// - "param:<name>"
			TokenLookup: "header: Authorization, query: admtk, cookie: admtk",
			// TokenLookup: "query:token",
			// TokenLookup: "cookie:token",

			// TokenHeadName is a string in the header. Default value is "Bearer"
			TokenHeadName: "Bearer",
			SendCookie:    true,
			CookieName:    "admtk",
			CookieMaxAge:  time.Hour * time.Duration(24),
			// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
			TimeFunc: time.Now,
		})
	}
	return auth
}

func GetLoginUserId(c *gin.Context) int64 {
	res, err := auth.GetClaimsFromJWT(c)
	if err != nil {
		return 0
	}

	if _, ok := res["user_id"]; !ok {
		return 0
	}
	return int64(res["user_id"].(float64))
}

func Login(ctx *gin.Context) {
	ctx.HTML(200, "login", gin.H{})
}

func Logout(ctx *gin.Context) {
	//clear session
	ctx.Redirect(http.StatusMovedPermanently, "/admin/login")
}
