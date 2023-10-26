package admin

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User demo
type User struct {
	UserId    int64
	Username  string
	Password  string
	FirstName string
	LastName  string
}

// the jwt middleware
func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	var identityKey = "UserId"
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "764da09d99b0e6dc",
		Key:         []byte("4ec85fbbeb9e3b90"),
		Timeout:     time.Hour * time.Duration(24),
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					"user_id":    v.UserId,
					"username":   v.Username,
					"first_name": v.FirstName,
					"last_name":  v.LastName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserId: func() int64 {
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
			username := loginVals.Username
			password := loginVals.Password

			users := []User{
				User{UserId: 1, Username: "admin", Password: "admin2023", FirstName: "d", LastName: "demo"},
			}

			for _, item := range users {
				if username == item.Username && password == item.Password {
					return &item, nil
				}

			}

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
			c.Redirect(301, "/admin")
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
		TokenLookup: "header: Authorization, query: token, cookie: token",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		SendCookie:    true,
		CookieName:    "token",
		CookieMaxAge:  time.Hour * time.Duration(24),
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}

func Login(ctx *gin.Context) {
	ctx.HTML(200, "login", gin.H{})
}

func Logout(ctx *gin.Context) {
	//clear session
	ctx.Redirect(http.StatusMovedPermanently, "/admin/login")
}
