package middlewares

import (
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const identityKey = "user_id"

type loginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	UserID   string
	Username string
}

type AuthMiddleware struct {
	jwt    *jwt.GinJWTMiddleware
	logger *zap.Logger
	viper  *viper.Viper
}

func NewAuthMiddleware(logger *zap.Logger, viper *viper.Viper) *AuthMiddleware {
	auth := AuthMiddleware{
		logger: logger,
		viper:  viper,
	}
	auth.initJwt()
	return &auth
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return m.jwt.MiddlewareFunc()
}

func (m *AuthMiddleware) Jwt() *jwt.GinJWTMiddleware {
	return m.jwt
}

func (m *AuthMiddleware) initJwt() {
	m.viper.SetDefault("jwt.secret_key", "secret key for dex")
	m.viper.SetDefault("jwt.realm", "trading engine auth")
	m.viper.SetDefault("jwt.timeout", time.Hour)
	m.viper.SetDefault("jwt.max_refresh", time.Hour)

	realm := m.viper.GetString("jwt.realm")
	secretKey := m.viper.GetString("jwt.secret_key")
	timeout := m.viper.GetDuration("jwt.timeout")
	maxRefresh := m.viper.GetDuration("jwt.max_refresh")

	//TODO load config
	mid, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           realm,
		Key:             []byte(secretKey),
		Timeout:         timeout,
		MaxRefresh:      maxRefresh,
		IdentityKey:     identityKey,
		PayloadFunc:     m.payloadFunc(),
		IdentityHandler: m.identityHandler(),
		Authenticator:   m.authenticator(),
		Authorizator:    m.authorizator(),
		Unauthorized:    m.unauthorized(),
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})
	if err != nil {
		m.logger.Panic("init jwt middleware error", zap.Error(err))
	}
	m.jwt = mid
}

func (m *AuthMiddleware) payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				identityKey: v.UserID,
			}
		}
		return jwt.MapClaims{}
	}
}

func (m *AuthMiddleware) identityHandler() func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		return &User{
			UserID: claims[identityKey].(string),
		}
	}
}

func (m *AuthMiddleware) authenticator() func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var loginVals loginRequest
		if err := c.ShouldBind(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}
		userID := loginVals.Username
		password := loginVals.Password

		if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
			return &User{
				Username: userID,
			}, nil
		}
		return nil, jwt.ErrFailedAuthentication
	}
}

func (m *AuthMiddleware) authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		if v, ok := data.(*User); ok && v.Username == "admin" {
			return true
		}
		return false
	}
}

func (m *AuthMiddleware) unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    code,
			"message": message,
		})
	}
}
